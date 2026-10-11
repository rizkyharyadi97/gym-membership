package service

import (
	"context"
	"errors"
	"strings"

	"gym-membership/internal/model"
	"gym-membership/internal/repository"
)

var (
	ErrTrainerNameRequired           = errors.New("trainer name is required")
	ErrTrainerSpecializationRequired = errors.New("trainer specialization is required")
)

type TrainerService interface {
	Create(ctx context.Context, trainer *model.Trainer) error
	GetByID(ctx context.Context, id uint) (*model.Trainer, error)
	GetAll(ctx context.Context) ([]model.Trainer, error)
	Update(ctx context.Context, trainer *model.Trainer) error
	Delete(ctx context.Context, id uint) error
}

type trainerService struct {
	trainerRepository repository.TrainerRepository
}

func NewTrainerService(
	trainerRepository repository.TrainerRepository,
) TrainerService {
	return &trainerService{
		trainerRepository: trainerRepository,
	}
}

func (s *trainerService) Create(ctx context.Context, trainer *model.Trainer) error {
	trainer.Name = strings.TrimSpace(trainer.Name)
	trainer.Specialization = strings.TrimSpace(trainer.Specialization)

	if trainer.Name == "" {
		return ErrTrainerNameRequired
	}

	if trainer.Specialization == "" {
		return ErrTrainerSpecializationRequired
	}

	return s.trainerRepository.Create(ctx, trainer)
}

func (s *trainerService) GetByID(ctx context.Context, id uint) (*model.Trainer, error) {
	return s.trainerRepository.GetByID(ctx, id)
}

func (s *trainerService) GetAll(ctx context.Context) ([]model.Trainer, error) {
	return s.trainerRepository.GetAll(ctx)
}

func (s *trainerService) Update(ctx context.Context, trainer *model.Trainer) error {
	if trainer.ID == 0 {
		return repository.ErrTrainerNotFound
	}

	if strings.TrimSpace(trainer.Name) == "" {
		return ErrTrainerNameRequired
	}

	if strings.TrimSpace(trainer.Specialization) == "" {
		return ErrTrainerSpecializationRequired
	}

	return s.trainerRepository.Update(ctx, trainer)
}

func (s *trainerService) Delete(ctx context.Context, id uint) error {
	return s.trainerRepository.Delete(ctx, id)
}
