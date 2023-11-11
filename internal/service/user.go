package service

import (
	"adorable-star/internal/dao"
	"adorable-star/internal/model"
	"adorable-star/internal/pkg/config"
	"adorable-star/internal/pkg/crawler"
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

var User = &UserService{dao.User, Token}

type UserService struct {
	d *dao.UserDAO
	s *TokenService
}

func (s *UserService) Register(email string, username string, pwd string) error {
	// TODO: 实现邮箱验证码

	// Encrypt pwd
	encryptedPwd, err := bcrypt.GenerateFromPassword([]byte(pwd+config.Config.Server.EncryptSalt), bcrypt.MinCost)
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

func (s *UserService) Login(name string, pwd string) (token string, uid int, err error) {
	// Find user from DB
	user, err := s.d.GetUserByUsernameOrEmail(name)
	if err != nil {
		err = errors.New("账户不存在")
		return
	}
	uid = user.ID

	// When pwd does not match
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pwd+config.Config.Server.EncryptSalt)); err != nil {
		err = errors.New("账号或密码错误")
		return
	}

	// When Jupiter Ed Data does not exist
	if _, err = dao.Jupiter.GetDataByUID(user.ID); err != nil {
		err = errors.New("userJupiterDataNotFound")
		return
	}

	// Generate token
	token, err = s.s.GenerateToken(user.ID, user.Status, user.Email)
	if err != nil {
		err = errors.New("internalErr")
		return
	}

	// Update user active time
	if dao.User.UpdateActiveTime(user.ID) != nil {
		err = errors.New("internalErr")
	}

	return
}

func (s *UserService) CompleteInfo(uid int, account string, pwd string) error {
	// Verify the Jupiter account given by the user
	if err := crawler.VerifyAccount(uid, account, pwd); err != nil {
		if err.Error() == "invalidJupiterAccount" {
			return errors.New("jupiter 账号或密码错误")
		}
		return errors.New("internalErr")
	}

	// Insert data to database
	if err := dao.Jupiter.InsertData(&model.JupiterData{UID: uid, Account: account, Password: pwd}); err != nil {
		if strings.Contains(err.Error(), "UNIQUE") {
			return errors.New("不可重复插入用户信息")
		}
		return err
	}

	return nil
}