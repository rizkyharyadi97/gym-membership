package repository

import (
	"context"
	"errors"

	"gym-membership/internal/model"

	"gorm.io/gorm"
)

var ErrGymClassNotFound = errors.New("gym class not found")

type GymClassRepository interface {
	Create(ctx context.Context, gymClass *model.GymClass) error
	GetByID(ctx context.Context, id uint) (*model.GymClass, error)
	GetAll(ctx context.Context) ([]model.GymClass, error)
	Update(ctx context.Context, gymClass *model.GymClass) error
	Delete(ctx context.Context, id uint) error
}

type gymClassRepository struct {
	db *gorm.DB
}

func NewGymClassRepository(db *gorm.DB) GymClassRepository {
	return &gymClassRepository{db: db}
}

func (r *gymClassRepository) Create(ctx context.Context, gymClass *model.GymClass) error {
	return r.db.WithContext(ctx).Create(gymClass).Error
}

func (r *gymClassRepository) GetByID(ctx context.Context, id uint) (*model.GymClass, error) {

	var gymClass model.GymClass

	err := r.db.WithContext(ctx).Preload("Trainer").First(&gymClass, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrGymClassNotFound
	}

	if err != nil {
		return nil, err
	}

	return &gymClass, nil
}

func (r *gymClassRepository) GetAll(ctx context.Context) ([]model.GymClass, error) {

	var gymClasses []model.GymClass

	err := r.db.WithContext(ctx).Preload("Trainer").Find(&gymClasses).Error
	if err != nil {
		return nil, err
	}

	return gymClasses, nil
}

func (r *gymClassRepository) Update(ctx context.Context, gymClass *model.GymClass) error {

	result := r.db.WithContext(ctx).Model(&model.GymClass{}).Where("id = ?", gymClass.ID).
		Updates(map[string]interface{}{
			"trainer_id":   gymClass.TrainerID,
			"class_name":   gymClass.ClassName,
			"schedule":     gymClass.Schedule,
			"quota":        gymClass.Quota,
			"availability": gymClass.Availability,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrGymClassNotFound
	}

	return nil
}

func (r *gymClassRepository) Delete(ctx context.Context, id uint) error {

	result := r.db.WithContext(ctx).Delete(&model.GymClass{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrGymClassNotFound
	}

	return nil
}
