package handler

import (
	"bytes"
	"context"
	"crypto/sha512"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"gym-membership/internal/integration"
	"gym-membership/internal/model"
	"gym-membership/middleware"
)

type PaymentHandler struct {
	db         *gorm.DB
	email      *integration.ResendClient
	httpClient *http.Client
}

func NewPaymentHandler(db *gorm.DB, email *integration.ResendClient) *PaymentHandler {
	return &PaymentHandler{db: db, email: email, httpClient: &http.Client{Timeout: 20 * time.Second}}
}

type createPaymentRequest struct {
	MembershipID uint `json:"membership_id"`
}
type midtransCreateResponse struct {
	Token         string   `json:"token"`
	RedirectURL   string   `json:"redirect_url"`
	ErrorMessages []string `json:"error_messages"`
}

// CreateMidtransPayment creates a pending local transaction and a Midtrans Snap checkout URL.
func (h *PaymentHandler) CreateMidtransPayment(c echo.Context) error {
	uid, ok := middleware.GetUserID(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"message": "authentication required",
		})
	}

	serverKey := strings.TrimSpace(os.Getenv("MIDTRANS_SERVER_KEY"))
	if serverKey == "" {
		return c.JSON(http.StatusServiceUnavailable, map[string]string{
			"message": "Midtrans is not configured; set MIDTRANS_SERVER_KEY",
		})
	}

	var req createPaymentRequest
	if err := c.Bind(&req); err != nil || req.MembershipID == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "membership_id is required",
		})
	}

	var user model.User
	if err := h.db.WithContext(c.Request().Context()).First(&user, uid).Error; err != nil {
		return handleError(c, err)
	}

	var membership model.Membership
	if err := h.db.WithContext(c.Request().Context()).First(&membership, req.MembershipID).Error; err != nil {
		return handleError(c, err)
	}

	if !membership.Availability {
		return c.JSON(http.StatusConflict, map[string]string{
			"message": "membership is unavailable",
		})
	}

	if membership.Costs <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "membership price must be greater than zero",
		})
	}

	orderID := fmt.Sprintf("GYM-%d-%d-%d", uid, membership.ID, time.Now().UnixNano())
	tx := model.Transaction{UserID: uid, MembershipID: membership.ID,
		TotalAmount:    membership.Costs,
		DepositPayment: 0, Status: "pending",
		TransactionDate: time.Now(),
		PaymentOrderID:  orderID,
	}

	if err := h.db.WithContext(c.Request().Context()).Create(&tx).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "could not create local payment transaction",
		})
	}

	baseURL := "https://app.sandbox.midtrans.com/snap/v1/transactions"
	if strings.EqualFold(os.Getenv("MIDTRANS_ENVIRONMENT"), "production") {
		baseURL = "https://app.midtrans.com/snap/v1/transactions"
	}

	payload := map[string]interface{}{
		"transaction_details": map[string]interface{}{
			"order_id":     orderID,
			"gross_amount": int64(membership.Costs),
		},
		"customer_details": map[string]interface{}{
			"first_name": user.Username,
			"email":      user.Email,
		},
		"item_details": []interface{}{
			map[string]interface{}{
				"id":       strconv.FormatUint(uint64(membership.ID), 10),
				"price":    int64(membership.Costs),
				"quantity": 1,
				"name":     membership.Name,
			},
		},
	}

	body, _ := json.Marshal(payload)
	ctx, cancel := context.WithTimeout(c.Request().Context(), 20*time.Second)

	defer cancel()

	apiReq, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL, bytes.NewReader(body))
	if err != nil {
		return c.JSON(500, map[string]string{
			"message": "could not prepare Midtrans request",
		})
	}

	apiReq.Header.Set("Content-Type", "application/json")
	apiReq.Header.Set("Accept", "application/json")
	apiReq.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(serverKey+":")))

	resp, err := h.httpClient.Do(apiReq)
	if err != nil {
		return c.JSON(http.StatusBadGateway, map[string]string{
			"message": "could not reach Midtrans; local transaction remains pending",
		})
	}

	defer resp.Body.Close()
	respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))

	var mid midtransCreateResponse
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		_ = h.db.WithContext(c.Request().Context()).Model(&model.Transaction{}).Where("id = ?", tx.ID).Update("status", "failed").Error
		return c.JSON(http.StatusBadGateway, map[string]interface{}{
			"message":           "Midtrans rejected payment request",
			"provider_status":   resp.StatusCode,
			"provider_response": json.RawMessage(respBody),
		})
	}

	if err := json.Unmarshal(respBody, &mid); err != nil || mid.Token == "" || mid.RedirectURL == "" {
		return c.JSON(http.StatusBadGateway, map[string]string{
			"message": "Midtrans returned an invalid response",
		})
	}

	_ = h.db.WithContext(c.Request().Context()).Model(&tx).Updates(map[string]interface{}{
		"payment_token": mid.Token,
		"payment_url":   mid.RedirectURL,
	}).Error

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "payment checkout created",
		"data": map[string]interface{}{
			"transaction_id": tx.ID,
			"order_id":       orderID,
			"amount":         membership.Costs,
			"status":         "pending",
			"token":          mid.Token,
			"redirect_url":   mid.RedirectURL},
	})
}

type midtransNotification struct {
	OrderID           string `json:"order_id"`
	StatusCode        string `json:"status_code"`
	GrossAmount       string `json:"gross_amount"`
	SignatureKey      string `json:"signature_key"`
	TransactionStatus string `json:"transaction_status"`
	FraudStatus       string `json:"fraud_status"`
	PaymentType       string `json:"payment_type"`
}

func (h *PaymentHandler) MidtransNotification(c echo.Context) error {
	serverKey := strings.TrimSpace(os.Getenv("MIDTRANS_SERVER_KEY"))
	if serverKey == "" {
		return c.JSON(503, map[string]string{
			"message": "Midtrans is not configured",
		})
	}

	var n midtransNotification
	if err := c.Bind(&n); err != nil || n.OrderID == "" || n.SignatureKey == "" {
		return c.JSON(400, map[string]string{
			"message": "invalid Midtrans notification",
		})
	}

	sum := sha512.Sum512([]byte(n.OrderID + n.StatusCode + n.GrossAmount + serverKey))
	expected := hex.EncodeToString(sum[:])
	if subtle.ConstantTimeCompare([]byte(strings.ToLower(expected)), []byte(strings.ToLower(n.SignatureKey))) != 1 {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"message": "invalid Midtrans signature",
		})
	}

	var tx model.Transaction
	if err := h.db.WithContext(c.Request().Context()).Where("payment_order_id = ?", n.OrderID).First(&tx).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"message": "transaction not found",
		})
	}

	status := tx.Status
	switch n.TransactionStatus {
	case "settlement":
		status = "paid"

	case "capture":
		if n.FraudStatus == "accept" || n.FraudStatus == "" {
			status = "paid"
		}

	case "pending":
		status = "pending"

	case "deny", "failure":
		status = "failed"

	case "cancel", "expire":
		status = "cancelled"
	}

	if err := h.db.WithContext(c.Request().Context()).Model(&tx).Update("status", status).Error; err != nil {
		return c.JSON(500, map[string]string{
			"message": "could not update transaction status",
		})
	}

	if status == "paid" && h.email != nil {
		var user model.User
		if h.db.WithContext(c.Request().Context()).First(&user, tx.UserID).Error == nil {
			if err := h.email.Send(c.Request().Context(), user.Email, "Pembayaran membership berhasil",
				fmt.Sprintf(
					"<p>Pembayaran membership Anda berhasil.</p><p>Nomor transaksi: %s</p><p>Total: %.2f</p>",
					n.OrderID, tx.TotalAmount)); err != nil {
				c.Logger().Errorf("Resend email failed: %v", err)
			}
		}
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "notification processed"})
}
