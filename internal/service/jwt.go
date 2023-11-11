package service

import (
	"adorable-star/internal/dao"
	"adorable-star/internal/pkg/config"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var Token = &TokenService{}

type TokenService struct{}

type tokenClaims struct {
	jwt.RegisteredClaims
	UID    int
	Status int
	Email  string
}

// Generate a jwt token with uid, user status and user email
func (s *TokenService) GenerateToken(uid int, status int, email string) (token string, err error) {
	// Initialize token with claims
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims{
		UID:    uid,
		Status: status,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    config.Config.Server.JwtIssuer,
			NotBefore: jwt.NewNumericDate(time.Now().Add(-time.Minute)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	})

	// Sign the token with encryption key
	token, err = t.SignedString([]byte(config.Config.Server.JwtEncrypt))
	return
}

// Verify if a token is valid
func (s *TokenService) VerifyToken(token string) error {
	// Decode token
	t, err := jwt.ParseWithClaims(token, &tokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.Config.Server.JwtEncrypt), nil
	})
	if err != nil {
		return err
	}

	// Check if valid
	if !t.Valid {
		return errors.New("invalidToken")
	}

	// Bind to claims
	if claims, ok := t.Claims.(*tokenClaims); ok {
		// Check if uid matches decoded info
		user, err := dao.User.GetUserByID(claims.UID)
		if err != nil || user.Status != claims.Status || user.Email != claims.Email {
			return errors.New("invalidToken")
		}

		// Update user active time
		if dao.User.UpdateActiveTime(claims.UID) != nil {
			return errors.New("invalidToken")
		}

		return nil
	} else {
		return errors.New("invalidToken")
	}
}
