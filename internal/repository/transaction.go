package repository

import (
	"context"
	"errors"

	"gym-membership/internal/model"

	"gorm.io/gorm"
)

var ErrTransactionNotFound = errors.New("transaction not found")

// TransactionRepository handles transaction data.
type TransactionRepository interface {
	Create(ctx context.Context, transaction *model.Transaction) error
	GetByID(ctx context.Context, id uint) (*model.Transaction, error)
	GetByUserID(ctx context.Context, userID uint) ([]model.Transaction, error)
	UpdateStatus(ctx context.Context, id uint, status string) error
}

// TransactionManager handles database transactions.
type TransactionManager interface {
	Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error
}

type transactionRepository struct {
	db *gorm.DB
}

type transactionManager struct {
	db *gorm.DB
}

// NewTransactionRepository creates a new transaction repository.
func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

// NewTransactionManager creates a new database transaction manager.
func NewTransactionManager(db *gorm.DB) TransactionManager {
	return &transactionManager{db: db}
}

// Create creates a new transaction.
func (r *transactionRepository) Create(ctx context.Context, transaction *model.Transaction) error {
	return r.db.WithContext(ctx).Create(transaction).Error
}

// GetByID retrieves a transaction by ID.
func (r *transactionRepository) GetByID(ctx context.Context, id uint) (*model.Transaction, error) {

	var transaction model.Transaction

	err := r.db.WithContext(ctx).Preload("User").Preload("Membership").First(&transaction, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrTransactionNotFound
	}

	if err != nil {
		return nil, err
	}

	return &transaction, nil
}

// GetByUserID retrieves all transactions belonging to a user.
func (r *transactionRepository) GetByUserID(ctx context.Context, userID uint) ([]model.Transaction, error) {

	var transactions []model.Transaction

	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Preload("Membership").
		Order("transaction_date DESC").
		Find(&transactions).
		Error

	if err != nil {
		return nil, err
	}

	return transactions, nil
}

// UpdateStatus updates the status of a transaction.
func (r *transactionRepository) UpdateStatus(ctx context.Context, id uint, status string) error {

	result := r.db.WithContext(ctx).Model(&model.Transaction{}).Where("id = ?", id).Update("status", status)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrTransactionNotFound
	}

	return nil
}

// Transaction executes a function inside a database transaction.
//
// If fn returns an error, GORM automatically rolls back.
// If fn returns nil, GORM automatically commits.
func (m *transactionManager) Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return m.db.WithContext(ctx).Transaction(fn)
}
