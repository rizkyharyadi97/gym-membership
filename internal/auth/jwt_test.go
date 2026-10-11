package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestJWTGenerateAndValidate(t *testing.T) {
	m := NewJWTManager("test-secret", time.Hour)
	token, err := m.GenerateToken(7)
	require.NoError(t, err)
	claims, err := m.ValidateToken(token)
	require.NoError(t, err)
	require.Equal(t, uint(7), claims.UserID)
}

func TestJWTRejectsInvalidToken(t *testing.T) {
	m := NewJWTManager("test-secret", time.Hour)
	_, err := m.ValidateToken("invalid.token.value")
	require.ErrorIs(t, err, ErrInvalidToken)
}
