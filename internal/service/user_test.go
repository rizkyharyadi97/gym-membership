package service

import (
	"context"
	"gym-membership/internal/model"
	"gym-membership/internal/repository"
	"testing"

	"github.com/stretchr/testify/require"
)

type mockUserRepository struct {
	users      map[uint]*model.User
	byEmail    map[string]*model.User
	byUsername map[string]*model.User
	created    *model.User
	deposit    float64
}

func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{users: map[uint]*model.User{}, byEmail: map[string]*model.User{}, byUsername: map[string]*model.User{}}
}
func (m *mockUserRepository) Create(_ context.Context, u *model.User) error {
	u.ID = 1
	m.created = u
	return nil
}
func (m *mockUserRepository) GetByID(_ context.Context, id uint) (*model.User, error) {
	u, ok := m.users[id]
	if !ok {
		return nil, repository.ErrUserNotFound
	}
	return u, nil
}
func (m *mockUserRepository) GetByEmail(_ context.Context, e string) (*model.User, error) {
	u, ok := m.byEmail[e]
	if !ok {
		return nil, repository.ErrUserNotFound
	}
	return u, nil
}
func (m *mockUserRepository) GetByUsername(_ context.Context, n string) (*model.User, error) {
	u, ok := m.byUsername[n]
	if !ok {
		return nil, repository.ErrUserNotFound
	}
	return u, nil
}
func (m *mockUserRepository) GetAll(_ context.Context) ([]model.User, error) {
	return nil, nil
}
func (m *mockUserRepository) UpdateDeposit(_ context.Context, _ uint, a float64) error {
	m.deposit += a
	return nil
}

func TestUserServiceRegisterHashesPassword(t *testing.T) {
	repo := newMockUserRepository()
	svc := NewUserService(repo)
	u := &model.User{Username: "rizky", Email: "rizky@mail.com", Password: "secret123"}
	err := svc.Register(context.Background(), u)
	require.NoError(t, err)
	require.NotEqual(t, "secret123", repo.created.Password)
	require.Equal(t, float64(0), repo.created.DepositAmount)
}
func TestUserServiceRegisterRejectsDuplicateUsername(t *testing.T) {
	repo := newMockUserRepository()
	repo.byUsername["rizky"] = &model.User{ID: 1, Username: "rizky"}
	svc := NewUserService(repo)
	err := svc.Register(context.Background(), &model.User{Username: "rizky", Email: "new@mail.com", Password: "secret123"})
	require.ErrorIs(t, err, ErrUserExists)
}
func TestUserServiceLoginInvalidPassword(t *testing.T) {
	repo := newMockUserRepository()
	repo.byEmail["rizky@mail.com"] = &model.User{ID: 1, Email: "rizky@mail.com", Password: "not-a-bcrypt-hash"}
	svc := NewUserService(repo)
	_, err := svc.Login(context.Background(), "rizky@mail.com", "secret123")
	require.ErrorIs(t, err, ErrInvalidCredentials)
}
