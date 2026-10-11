package service

import (
	"context"
	"errors"
	"time"

	"gym-membership/internal/model"
	"gym-membership/internal/repository"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrInsufficientDeposit   = errors.New("insufficient deposit")
	ErrInvalidTransaction    = errors.New("invalid transaction")
	ErrMembershipUnavailable = errors.New("membership is unavailable")
)

type TransactionService interface {
	PurchaseMembership(ctx context.Context, userID uint, membershipID uint) (*model.Transaction, error)

	GetByID(ctx context.Context, id uint) (*model.Transaction, error)

	GetByUserID(ctx context.Context, userID uint) ([]model.Transaction, error)
}

type transactionService struct {
	transactionRepository repository.TransactionRepository
	transactionManager    repository.TransactionManager
}

func NewTransactionService(transactionRepository repository.TransactionRepository, transactionManager repository.TransactionManager) TransactionService {
	return &transactionService{
		transactionRepository: transactionRepository,
		transactionManager:    transactionManager,
	}
}

func (s *transactionService) PurchaseMembership(ctx context.Context, userID uint, membershipID uint) (*model.Transaction, error) {

	var transaction model.Transaction

	err := s.transactionManager.Transaction(ctx,
		func(tx *gorm.DB) error {
			var user model.User
			var membership model.Membership

			// Lock the user row so two simultaneous purchases
			// cannot spend the same deposit balance.
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&user, userID).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return repository.ErrUserNotFound
				}

				return err
			}

			// Get membership.
			if err := tx.First(&membership, membershipID).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return repository.ErrMembershipNotFound
				}

				return err
			}

			// Check membership availability.
			if !membership.Availability {
				return ErrMembershipUnavailable
			}

			// Check deposit balance.
			if user.DepositAmount < membership.Costs {
				return ErrInsufficientDeposit
			}

			// Deduct deposit.
			newDeposit := user.DepositAmount - membership.Costs
			if err := tx.Model(&model.User{}).Where("id = ?", user.ID).Update("deposit_amount", newDeposit).Error; err != nil {
				return err
			}

			// Create transaction.
			transaction = model.Transaction{
				UserID:          user.ID,
				MembershipID:    membership.ID,
				TotalAmount:     membership.Costs,
				DepositPayment:  membership.Costs,
				Status:          "paid",
				TransactionDate: time.Now(),
			}

			if err := tx.Create(&transaction).Error; err != nil {
				return err
			}

			return nil
		},
	)

	if err != nil {
		return nil, err
	}

	return &transaction, nil
}

func (s *transactionService) GetByID(ctx context.Context, id uint) (*model.Transaction, error) {
	return s.transactionRepository.GetByID(ctx, id)
}

func (s *transactionService) GetByUserID(ctx context.Context, userID uint) ([]model.Transaction, error) {
	return s.transactionRepository.GetByUserID(ctx, userID)
}
