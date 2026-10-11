package main

import (
	"fmt"
	"gym-membership/config"
	"gym-membership/internal/auth"
	"gym-membership/internal/handler"
	"gym-membership/internal/integration"
	"gym-membership/internal/repository"
	"gym-membership/internal/service"
	"gym-membership/middleware"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

func main() {
	db := config.ConnectDB()
	jwtConfig, err := config.LoadJWTConfig()
	if err != nil {
		log.Fatalf("JWT configuration error: %v", err)
	}

	jwtManager := auth.NewJWTManager(jwtConfig.SecretKey, jwtConfig.ExpirationHours)

	userRepo := repository.NewUserRepository(db)
	membershipRepo := repository.NewMembershipRepository(db)
	trainerRepo := repository.NewTrainerRepository(db)
	classRepo := repository.NewGymClassRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)
	transactionManager := repository.NewTransactionManager(db)
	bookingRepo := repository.NewBookingRepository(db)
	userSvc := service.NewUserService(userRepo)
	membershipSvc := service.NewMembershipService(membershipRepo)
	trainerSvc := service.NewTrainerService(trainerRepo)
	classSvc := service.NewGymClassService(classRepo, trainerRepo)
	transactionSvc := service.NewTransactionService(transactionRepo, transactionManager)
	bookingSvc := service.NewBookingService(bookingRepo, userRepo, classRepo)
	authHandler := handler.NewAuthHandler(userSvc, jwtManager)
	emailClient := integration.NewResendClient()
	userHandler := handler.NewUserHandler(userSvc, emailClient)
	paymentHandler := handler.NewPaymentHandler(db, emailClient)
	membershipHandler := handler.NewMembershipHandler(membershipSvc)
	trainerHandler := handler.NewTrainerHandler(trainerSvc)
	classHandler := handler.NewGymClassHandler(classSvc)
	transactionHandler := handler.NewTransactionHandler(transactionSvc)
	bookingHandler := handler.NewBookingHandler(bookingSvc)

	e := echo.New()
	e.HideBanner = true

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status": "ok", "message": "Gym Membership API is running",
		})
	})

	e.GET("/health/db", func(c echo.Context) error {
		sqlDB, err := db.DB()
		if err != nil {
			return c.JSON(500, map[string]interface{}{
				"status": "error", "message": "failed to access database",
			})
		}
		if err = sqlDB.Ping(); err != nil {
			return c.JSON(500, map[string]interface{}{
				"status": "error", "message": "database connection failed",
			})
		}
		return c.JSON(200, map[string]interface{}{
			"status": "ok", "message": "database connection is healthy",
		})
	})

	api := e.Group("/api/v1")

	api.POST("/users", userHandler.Register)
	api.POST("/auth/login", authHandler.Login)
	api.GET("/users/:id", userHandler.GetByID)
	api.GET("/memberships", membershipHandler.GetAll)
	api.GET("/memberships/:id", membershipHandler.GetByID)
	api.GET("/trainers", trainerHandler.GetAll)
	api.GET("/trainers/:id", trainerHandler.GetByID)
	api.GET("/gym-classes", classHandler.GetAll)
	api.GET("/gym-classes/:id", classHandler.GetByID)

	protected := api.Group("")

	protected.Use(middleware.JWTMiddleware(jwtManager))
	protected.POST("/users/me/deposit", userHandler.AddDeposit)
	protected.POST("/memberships", membershipHandler.Create)
	protected.PUT("/memberships/:id", membershipHandler.Update)
	protected.DELETE("/memberships/:id", membershipHandler.Delete)
	protected.POST("/trainers", trainerHandler.Create)
	protected.PUT("/trainers/:id", trainerHandler.Update)
	protected.DELETE("/trainers/:id", trainerHandler.Delete)
	protected.POST("/gym-classes", classHandler.Create)
	protected.PUT("/gym-classes/:id", classHandler.Update)
	protected.DELETE("/gym-classes/:id", classHandler.Delete)
	protected.POST("/transactions/purchase", transactionHandler.PurchaseMembership)
	protected.POST("/payments/midtrans", paymentHandler.CreateMidtransPayment)
	api.POST("/payments/midtrans/notification", paymentHandler.MidtransNotification)
	protected.GET("/transactions/:id", transactionHandler.GetByID)
	protected.GET("/me/transactions", transactionHandler.GetByUserID)
	protected.POST("/bookings", bookingHandler.Create)
	protected.GET("/bookings/:id", bookingHandler.GetByID)
	protected.GET("/me/bookings", bookingHandler.GetByUserID)
	protected.GET("/gym-classes/:classID/bookings", bookingHandler.GetByClassID)
	protected.PATCH("/bookings/:id/cancel", bookingHandler.Cancel)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}
	address := fmt.Sprintf(":%s", port)
	log.Printf("server running on %s", address)
	if err := e.Start(address); err != nil {
		log.Fatal(err)
	}
}
