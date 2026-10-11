package service

import (
	"context"
	"errors"
	"strings"

	"gym-membership/internal/model"
	"gym-membership/internal/repository"
)

var (
	ErrMembershipNameRequired = errors.New("membership name is required")
	ErrInvalidMembershipCost  = errors.New("membership cost cannot be negative")
	ErrInvalidDuration        = errors.New("membership duration must be greater than zero")
)

type MembershipService interface {
	Create(ctx context.Context, membership *model.Membership) error
	GetByID(ctx context.Context, id uint) (*model.Membership, error)
	GetAll(ctx context.Context) ([]model.Membership, error)
	Update(ctx context.Context, membership *model.Membership) error
	Delete(ctx context.Context, id uint) error
}

type membershipService struct {
	membershipRepository repository.MembershipRepository
}

func NewMembershipService(membershipRepository repository.MembershipRepository) MembershipService {
	return &membershipService{membershipRepository: membershipRepository}
}

func (s *membershipService) Create(ctx context.Context, membership *model.Membership) error {
	membership.Name = strings.TrimSpace(membership.Name)
	membership.Category = strings.TrimSpace(membership.Category)

	if membership.Name == "" {
		return ErrMembershipNameRequired
	}

	if membership.Costs < 0 {
		return ErrInvalidMembershipCost
	}

	if membership.DurationDays <= 0 {
		return ErrInvalidDuration
	}

	return s.membershipRepository.Create(ctx, membership)
}

func (s *membershipService) GetByID(ctx context.Context, id uint) (*model.Membership, error) {
	return s.membershipRepository.GetByID(ctx, id)
}

func (s *membershipService) GetAll(ctx context.Context) ([]model.Membership, error) {
	return s.membershipRepository.GetAll(ctx)
}

func (s *membershipService) Update(ctx context.Context, membership *model.Membership) error {
	if membership.ID == 0 {
		return repository.ErrMembershipNotFound
	}

	if strings.TrimSpace(membership.Name) == "" {
		return ErrMembershipNameRequired
	}

	if membership.Costs < 0 {
		return ErrInvalidMembershipCost
	}

	if membership.DurationDays <= 0 {
		return ErrInvalidDuration
	}

	return s.membershipRepository.Update(ctx, membership)
}

func (s *membershipService) Delete(ctx context.Context, id uint) error {
	return s.membershipRepository.Delete(ctx, id)
}
