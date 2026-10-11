package model

import "time"

type Transaction struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	UserID          uint      `gorm:"not null;index" json:"user_id"`
	MembershipID    uint      `gorm:"not null;index" json:"membership_id"`
	TotalAmount     float64   `gorm:"type:numeric(12,2);not null" json:"total_amount"`
	DepositPayment  float64   `gorm:"type:numeric(12,2);not null;default:0" json:"deposit_payment"`
	Status          string    `gorm:"size:20;not null;default:pending;index" json:"status"`
	TransactionDate time.Time `gorm:"not null;index" json:"transaction_date"`
	PaymentOrderID  string    `gorm:"size:64;uniqueIndex;default:null" json:"payment_order_id,omitempty"`
	PaymentToken    string    `gorm:"size:255" json:"payment_token,omitempty"`
	PaymentURL      string    `gorm:"type:text" json:"payment_url,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`

	// Relations
	User       User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Membership Membership `gorm:"foreignKey:MembershipID" json:"membership,omitempty"`
}
