package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestJWT_GenerateAndValidate(t *testing.T) {
	secret := "test-secret-key"
	userID := uint(42)
	clientIP := "192.168.1.100"
	expHours := 2

	// 1. Test successful token generation and validation
	tokenStr, err := GenerateToken(userID, clientIP, secret, expHours)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenStr)

	claims, err := ValidateToken(tokenStr, secret)
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, clientIP, claims.IP)

	// 2. Test validation failure due to invalid secret
	wrongSecret := "wrong-secret-key"
	wrongClaims, err := ValidateToken(tokenStr, wrongSecret)
	assert.Error(t, err)
	assert.Nil(t, wrongClaims)
}

func TestJWT_ExpiredToken(t *testing.T) {
	secret := "test-secret-key"
	userID := uint(42)
	clientIP := "192.168.1.100"
	// Use negative expiration hours to simulate an expired token
	expHours := -1

	tokenStr, err := GenerateToken(userID, clientIP, secret, expHours)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenStr)

	claims, err := ValidateToken(tokenStr, secret)
	assert.Error(t, err)
	assert.Nil(t, claims)
	assert.Contains(t, err.Error(), "token is expired")
}

func TestJWT_InvalidTokenFormat(t *testing.T) {
	secret := "test-secret-key"
	invalidTokenStr := "this.is.not.a.valid.jwt.token"

	claims, err := ValidateToken(invalidTokenStr, secret)
	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestJWT_UnexpectedSigningMethod(t *testing.T) {
	secret := "test-secret-key"
	claims := JWTClaims{
		UserID: 42,
		IP:     "127.0.0.1",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		},
	}

	// Generate a token signed with a different method (e.g. None)
	tokenNone := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
	tokenStr, err := tokenNone.SignedString(jwt.UnsafeAllowNoneSignatureType)
	assert.NoError(t, err)

	_, err = ValidateToken(tokenStr, secret)
	assert.Error(t, err)
}
