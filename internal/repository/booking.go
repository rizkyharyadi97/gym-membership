package repository

import (
	"context"
	"errors"
	"time"

	"gym-membership/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrBookingNotFound         = errors.New("booking not found")
	ErrBookingClassUnavailable = errors.New("gym class is unavailable")
	ErrBookingClassFull        = errors.New("gym class quota is full")
	ErrBookingAlreadyExists    = errors.New("user already booked this class")
)

type BookingRepository interface {
	Create(ctx context.Context, booking *model.Booking) error
	GetByID(ctx context.Context, id uint) (*model.Booking, error)
	GetByUserID(ctx context.Context, userID uint) ([]model.Booking, error)
	GetByClassID(ctx context.Context, classID uint) ([]model.Booking, error)
	CountActiveByClassID(ctx context.Context, classID uint) (int64, error)
	CreateAtomic(ctx context.Context, userID uint, classID uint) (*model.Booking, error)
	UpdateStatus(ctx context.Context, id uint, status string) error
}

type bookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) BookingRepository {
	return &bookingRepository{db: db}
}

func (r *bookingRepository) Create(ctx context.Context, booking *model.Booking) error {
	return r.db.WithContext(ctx).Create(booking).Error
}

func (r *bookingRepository) GetByID(ctx context.Context, id uint) (*model.Booking, error) {

	var booking model.Booking

	err := r.db.WithContext(ctx).Preload("User").Preload("GymClass").First(&booking, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrBookingNotFound
	}

	if err != nil {
		return nil, err
	}

	return &booking, nil
}

func (r *bookingRepository) GetByUserID(ctx context.Context, userID uint) ([]model.Booking, error) {

	var bookings []model.Booking

	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Preload("GymClass").Order("booking_date DESC").
		Find(&bookings).Error

	if err != nil {
		return nil, err
	}

	return bookings, nil
}

func (r *bookingRepository) GetByClassID(ctx context.Context, classID uint) ([]model.Booking, error) {

	var bookings []model.Booking

	err := r.db.WithContext(ctx).Where("gym_class_id = ?", classID).Preload("User").Order("booking_date ASC").
		Find(&bookings).Error

	if err != nil {
		return nil, err
	}

	return bookings, nil
}

func (r *bookingRepository) CountActiveByClassID(ctx context.Context, classID uint) (int64, error) {

	var count int64

	err := r.db.WithContext(ctx).Model(&model.Booking{}).Where("gym_class_id = ?", classID).Where("status = ?", "booked").
		Count(&count).Error

	return count, err
}

func (r *bookingRepository) UpdateStatus(ctx context.Context, id uint, status string) error {

	result := r.db.WithContext(ctx).Model(&model.Booking{}).Where("id = ?", id).Update("status", status)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrBookingNotFound
	}

	return nil
}

func (r *bookingRepository) CreateAtomic(ctx context.Context, userID uint, classID uint) (*model.Booking, error) {

	var booking model.Booking
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var class model.GymClass
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&class, classID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrGymClassNotFound
			}
			return err
		}
		if !class.Availability {
			return ErrBookingClassUnavailable
		}
		var count int64
		if err := tx.Model(&model.Booking{}).Where("gym_class_id = ? AND status = ?", classID, "booked").Count(&count).Error; err != nil {
			return err
		}
		if count >= int64(class.Quota) {
			return ErrBookingClassFull
		}
		var existing int64
		if err := tx.Model(&model.Booking{}).Where("user_id = ? AND gym_class_id = ? AND status = ?", userID, classID, "booked").Count(&existing).Error; err != nil {
			return err
		}
		if existing > 0 {
			return ErrBookingAlreadyExists
		}
		booking = model.Booking{UserID: userID, GymClassID: classID, BookingDate: time.Now(), Status: "booked"}
		return tx.Create(&booking).Error
	})
	if err != nil {
		return nil, err
	}
	return &booking, nil
}
