package repository

import (
	"context"
	"errors"

	"gym-membership/internal/model"

	"gorm.io/gorm"
)

var ErrMembershipNotFound = errors.New("membership not found")

type MembershipRepository interface {
	Create(ctx context.Context, membership *model.Membership) error
	GetByID(ctx context.Context, id uint) (*model.Membership, error)
	GetAll(ctx context.Context) ([]model.Membership, error)
	Update(ctx context.Context, membership *model.Membership) error
	Delete(ctx context.Context, id uint) error
}

type membershipRepository struct {
	db *gorm.DB
}

func NewMembershipRepository(db *gorm.DB) MembershipRepository {
	return &membershipRepository{db: db}
}

func (r *membershipRepository) Create(ctx context.Context, membership *model.Membership) error {
	return r.db.WithContext(ctx).Create(membership).Error
}

func (r *membershipRepository) GetByID(ctx context.Context, id uint) (*model.Membership, error) {

	var membership model.Membership

	err := r.db.WithContext(ctx).First(&membership, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrMembershipNotFound
	}

	if err != nil {
		return nil, err
	}

	return &membership, nil
}

func (r *membershipRepository) GetAll(ctx context.Context) ([]model.Membership, error) {

	var memberships []model.Membership

	err := r.db.WithContext(ctx).Find(&memberships).Error
	if err != nil {
		return nil, err
	}

	return memberships, nil
}

func (r *membershipRepository) Update(ctx context.Context, membership *model.Membership) error {
	result := r.db.WithContext(ctx).Model(&model.Membership{}).Where("id = ?", membership.ID).
		Updates(map[string]interface{}{
			"name":          membership.Name,
			"costs":         membership.Costs,
			"duration_days": membership.DurationDays,
			"availability":  membership.Availability,
			"category":      membership.Category,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrMembershipNotFound
	}

	return nil
}

func (r *membershipRepository) Delete(ctx context.Context, id uint) error {

	result := r.db.WithContext(ctx).Delete(&model.Membership{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrMembershipNotFound
	}

	return nil
}
