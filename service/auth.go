package service

import (
	"adorable-star/config"
	"adorable-star/dao"
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

var Auth = &AuthService{dao.User, Token}

type AuthService struct {
	d *dao.UserDAO
	s *TokenService
}

func (s *AuthService) Register(email string, username string, pwd string) error {
	// TODO: 实现邮箱验证码

	// Encrypt pwd
	encryptedPwd, err := bcrypt.GenerateFromPassword([]byte(pwd+config.ENCRYPT_SALT), bcrypt.MinCost)
	if err != nil {
		return errors.New("internalErr")
	}

	// Insert user to DB
	if err = s.d.InsertUser(email, username, string(encryptedPwd[:])); err != nil {
		if strings.Contains(err.Error(), "UNIQUE") {
			return errors.New("使用此邮箱或用户名的账户已经存在")
		}
		return errors.New("internalErr")
	}

	return nil
}

func (s *AuthService) Login(name string, pwd string) (token string, err error) {
	// Find user from DB
	user, err := s.d.GetUserByUsernameOrEmail(name)
	if err != nil {
		err = errors.New("账户不存在")
		return
	}

	// When pwd does not match
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pwd+config.ENCRYPT_SALT)); err != nil {
		err = errors.New("账号或密码错误")
		return
	}

	// Generate token
	token, err = s.s.GenerateToken(user.ID, user.Status, user.Email)
	if err != nil {
		err = errors.New("internalErr")
	}

	return
}
