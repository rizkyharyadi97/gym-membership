package service

import (
	"context"
	"errors"

	"gym-membership/internal/model"
	"gym-membership/internal/repository"
)

var (
	ErrClassUnavailable = errors.New("gym class is unavailable")
	ErrClassFull        = errors.New("gym class quota is full")
	ErrAlreadyBooked    = errors.New("user already booked this class")
)

type BookingService interface {
	Create(ctx context.Context, userID uint, gymClassID uint) (*model.Booking, error)
	GetByID(ctx context.Context, id uint) (*model.Booking, error)
	GetByUserID(ctx context.Context, userID uint) ([]model.Booking, error)
	GetByClassID(ctx context.Context, classID uint) ([]model.Booking, error)
	Cancel(ctx context.Context, bookingID uint) error
}

type bookingService struct {
	bookingRepository  repository.BookingRepository
	userRepository     repository.UserRepository
	gymClassRepository repository.GymClassRepository
}

func NewBookingService(bookingRepository repository.BookingRepository, userRepository repository.UserRepository, gymClassRepository repository.GymClassRepository) BookingService {
	return &bookingService{
		bookingRepository:  bookingRepository,
		userRepository:     userRepository,
		gymClassRepository: gymClassRepository,
	}
}

func (s *bookingService) Create(ctx context.Context, userID uint, gymClassID uint) (*model.Booking, error) {
	if _, err := s.userRepository.GetByID(ctx, userID); err != nil {
		return nil, err
	}
	booking, err := s.bookingRepository.CreateAtomic(ctx, userID, gymClassID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrBookingClassUnavailable):
			return nil, ErrClassUnavailable
		case errors.Is(err, repository.ErrBookingClassFull):
			return nil, ErrClassFull
		case errors.Is(err, repository.ErrBookingAlreadyExists):
			return nil, ErrAlreadyBooked
		default:
			return nil, err
		}
	}
	return booking, nil
}

func (s *bookingService) GetByID(ctx context.Context, id uint) (*model.Booking, error) {
	return s.bookingRepository.GetByID(ctx, id)
}

func (s *bookingService) GetByUserID(ctx context.Context, userID uint) ([]model.Booking, error) {
	return s.bookingRepository.GetByUserID(ctx, userID)
}

func (s *bookingService) GetByClassID(ctx context.Context, classID uint) ([]model.Booking, error) {
	return s.bookingRepository.GetByClassID(ctx, classID)
}

func (s *bookingService) Cancel(ctx context.Context, bookingID uint) error {
	booking, err := s.bookingRepository.GetByID(ctx, bookingID)
	if err != nil {
		return err
	}

	if booking.Status == "cancelled" {
		return nil
	}

	return s.bookingRepository.UpdateStatus(ctx, bookingID, "cancelled")
}
