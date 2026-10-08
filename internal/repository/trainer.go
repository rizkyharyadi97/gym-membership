package repository

import (
	"context"
	"errors"

	"gym-membership/internal/model"

	"gorm.io/gorm"
)

var ErrTrainerNotFound = errors.New("trainer not found")

type TrainerRepository interface {
	Create(ctx context.Context, trainer *model.Trainer) error
	GetByID(ctx context.Context, id uint) (*model.Trainer, error)
	GetAll(ctx context.Context) ([]model.Trainer, error)
	Update(ctx context.Context, trainer *model.Trainer) error
	Delete(ctx context.Context, id uint) error
}

type trainerRepository struct {
	db *gorm.DB
}

func NewTrainerRepository(db *gorm.DB) TrainerRepository {
	return &trainerRepository{db: db}
}

func (r *trainerRepository) Create(ctx context.Context, trainer *model.Trainer) error {
	return r.db.WithContext(ctx).Create(trainer).Error
}

func (r *trainerRepository) GetByID(ctx context.Context, id uint) (*model.Trainer, error) {

	var trainer model.Trainer

	err := r.db.WithContext(ctx).First(&trainer, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrTrainerNotFound
	}

	if err != nil {
		return nil, err
	}

	return &trainer, nil
}

func (r *trainerRepository) GetAll(ctx context.Context) ([]model.Trainer, error) {

	var trainers []model.Trainer

	err := r.db.WithContext(ctx).Find(&trainers).Error
	if err != nil {
		return nil, err
	}

	return trainers, nil
}

func (r *trainerRepository) Update(ctx context.Context, trainer *model.Trainer) error {
	result := r.db.WithContext(ctx).Model(&model.Trainer{}).Where("id = ?", trainer.ID).
		Updates(map[string]interface{}{
			"name":           trainer.Name,
			"specialization": trainer.Specialization,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrTrainerNotFound
	}

	return nil
}

func (r *trainerRepository) Delete(ctx context.Context, id uint) error {

	result := r.db.WithContext(ctx).Delete(&model.Trainer{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrTrainerNotFound
	}

	return nil
}
