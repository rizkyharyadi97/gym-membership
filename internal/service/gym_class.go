package service

import (
	"context"
	"errors"
	"strings"

	"gym-membership/internal/model"
	"gym-membership/internal/repository"
)

var (
	ErrClassNameRequired = errors.New("class name is required")
	ErrInvalidQuota      = errors.New("quota must be greater than zero")
	ErrInvalidSchedule   = errors.New("schedule is required")
)

type GymClassService interface {
	Create(ctx context.Context, gymClass *model.GymClass) error
	GetByID(ctx context.Context, id uint) (*model.GymClass, error)
	GetAll(ctx context.Context) ([]model.GymClass, error)
	Update(ctx context.Context, gymClass *model.GymClass) error
	Delete(ctx context.Context, id uint) error
}

type gymClassService struct {
	gymClassRepository repository.GymClassRepository
	trainerRepository  repository.TrainerRepository
}

func NewGymClassService(gymClassRepository repository.GymClassRepository, trainerRepository repository.TrainerRepository) GymClassService {
	return &gymClassService{
		gymClassRepository: gymClassRepository,
		trainerRepository:  trainerRepository,
	}
}

func (s *gymClassService) Create(ctx context.Context, gymClass *model.GymClass) error {
	gymClass.ClassName = strings.TrimSpace(gymClass.ClassName)

	if gymClass.ClassName == "" {
		return ErrClassNameRequired
	}

	if gymClass.TrainerID == 0 {
		return repository.ErrTrainerNotFound
	}

	if gymClass.Quota <= 0 {
		return ErrInvalidQuota
	}

	if gymClass.Schedule.IsZero() {
		return ErrInvalidSchedule
	}

	// Make sure trainer exists.
	_, err := s.trainerRepository.GetByID(ctx, gymClass.TrainerID)
	if err != nil {
		return err
	}

	return s.gymClassRepository.Create(ctx, gymClass)
}

func (s *gymClassService) GetByID(ctx context.Context, id uint) (*model.GymClass, error) {
	return s.gymClassRepository.GetByID(ctx, id)
}

func (s *gymClassService) GetAll(ctx context.Context) ([]model.GymClass, error) {
	return s.gymClassRepository.GetAll(ctx)
}

func (s *gymClassService) Update(ctx context.Context, gymClass *model.GymClass) error {
	if gymClass.ID == 0 {
		return repository.ErrGymClassNotFound
	}

	gymClass.ClassName = strings.TrimSpace(gymClass.ClassName)

	if gymClass.ClassName == "" {
		return ErrClassNameRequired
	}

	if gymClass.Quota <= 0 {
		return ErrInvalidQuota
	}

	if gymClass.Schedule.IsZero() {
		return ErrInvalidSchedule
	}

	_, err := s.trainerRepository.GetByID(ctx, gymClass.TrainerID)
	if err != nil {
		return err
	}

	return s.gymClassRepository.Update(ctx, gymClass)
}

func (s *gymClassService) Delete(ctx context.Context, id uint) error {
	return s.gymClassRepository.Delete(ctx, id)
}
