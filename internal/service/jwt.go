package service

import (
	"adorable-star/internal/dao"
	"adorable-star/internal/pkg/config"
	"errors"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var Token = &TokenService{}

type TokenService struct{}

type tokenClaims struct {
	jwt.RegisteredClaims
	UID    int    `json:"uid"`
	Status int    `json:"status"`
	Email  string `json:"email"`
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
func (s *TokenService) VerifyToken(ctx *gin.Context) (claims *tokenClaims, err error) {
	// Decode token
	claims = &tokenClaims{}
	t, err := jwt.ParseWithClaims(ctx.Request.Header.Get("Token"), claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.Config.Server.JwtEncrypt), nil
	})
	if err != nil && !strings.Contains(err.Error(), "expired") {
		return
	}

	// Check if uid matches decoded info
	user, err := dao.User.GetUserByID(claims.UID)
	if err != nil || user.Status != claims.Status || user.Email != claims.Email {
		return nil, errors.New("invalidToken")
	}

	// Check if valid (If the token does not expire for more than 24 hours, renew token)
	if !t.Valid {
		if withinLimit(claims.ExpiresAt.Unix(), 3600*24) {
			// Renew the token
			token, err := s.GenerateToken(claims.UID, claims.Status, claims.Email)
			if err != nil {
				return nil, errors.New("internalErr")
			}
			ctx.Header("New-Token", token)
		} else {
			err = errors.New("invalidToken")
			return
		}
	}

	// Update user active time
	if dao.User.UpdateActiveTime(claims.UID) != nil {
		return nil, errors.New("internalErr")
	}

	return
}

// Check whether a past time until now exceed the limited time
func withinLimit(expiredAt int64, timeLimit int64) bool {
	now := time.Now().Unix()
	return now-expiredAt < timeLimit
}
