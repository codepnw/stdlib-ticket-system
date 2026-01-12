package jwttoken

import (
	"errors"
	"fmt"
	"time"

	"github.com/codepnw/stdlib-ticket-system/internal/config"
	"github.com/codepnw/stdlib-ticket-system/internal/features/user"
	"github.com/golang-jwt/jwt/v5"
)

type JWTToken interface {
	GenerateAccessToken(u user.User) (string, error)
	GenerateRefreshToken(u user.User) (string, error)
	VerifyAccessToken(tokenStr string) (*Payload, error)
	VerifyRefreshToken(tokenStr string) (*Payload, error)
}

type jwtToken struct {
	secretKey  string
	refreshKey string
}

func NewJWT(secretKey, refreshKey string) (JWTToken, error) {
	if secretKey == "" || refreshKey == "" {
		return nil, errors.New("secret key & refresh key is required")
	}
	return &jwtToken{
		secretKey:  secretKey,
		refreshKey: refreshKey,
	}, nil
}

type UserClaims struct {
	ID       int64
	Username string
	*jwt.RegisteredClaims
}

func (j *jwtToken) GenerateAccessToken(u user.User) (string, error) {
	return j.generateToken([]byte(j.secretKey), u, config.AccessTokenDuration)
}

func (j *jwtToken) GenerateRefreshToken(u user.User) (string, error) {
	return j.generateToken([]byte(j.refreshKey), u, config.RefreshTokenDuration)
}

func (j *jwtToken) generateToken(key []byte, u user.User, duration time.Duration) (string, error) {
	claims := &UserClaims{
		ID:       u.ID,
		Username: u.Username,
		RegisteredClaims: &jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	ss, err := token.SignedString(key)
	if err != nil {
		return "", fmt.Errorf("sign token failed: %w", err)
	}
	return ss, nil
}

// ============= Verify Token =================

func (j *jwtToken) VerifyAccessToken(tokenStr string) (*Payload, error) {
	return j.verifyToken([]byte(j.secretKey), tokenStr)
}

func (j *jwtToken) VerifyRefreshToken(tokenStr string) (*Payload, error) {
	return j.verifyToken([]byte(j.refreshKey), tokenStr)
}

type Payload struct {
	UserID    int64     `json:"user_id"`
	Username  string    `json:"username"`
}

func (j *jwtToken) verifyToken(key []byte, tokenStr string) (*Payload, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &UserClaims{}, func(t *jwt.Token) (any, error) {
		return key, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, err
	}

	if time.Now().After(claims.ExpiresAt.Time) {
		return nil, errors.New("token expires")
	}

	payload := &Payload{
		UserID:    claims.ID,
		Username:  claims.Username,
	}
	return payload, nil
}
