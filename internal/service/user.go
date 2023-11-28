package service

import (
	"adorable-star/internal/dao"
	"adorable-star/internal/model"
	"adorable-star/internal/pkg/config"
	"adorable-star/internal/pkg/crawler"
	"errors"
	"math/rand"
	"net/smtp"
	"strconv"
	"strings"
	"time"

	"github.com/jordan-wright/email"
	"golang.org/x/crypto/bcrypt"
)

var User = &UserService{dao.User, Token}

type UserService struct {
	d *dao.UserDAO
	s *TokenService
}

func (s *UserService) SendValidationCode(userMail string) error {
	// If user with this email already exist
	_, err := s.d.GetUserByUsernameOrEmail(userMail)
	if err == nil {
		return errors.New("使用此邮箱的账户已经存在")
	}

	// Check if there's unexpired validation code (3 min)
	expiration, err := dao.Redis.TTL("vc-" + userMail).Result()
	if err != nil {
		return errors.New("internalErr")
	}
	if expiration > 120 {
		return errors.New("验证码仍在5分钟有效期内 (没小于2分钟禁止重发)")
	}

	// Generate random code with 6 digits
	code := int(rand.Float32()*499999) + 100000

	// Put the code into Redis for 5 minute expiration
	dao.Redis.Set("vc-"+userMail, code, time.Minute*5)

	// Send email with verification code
	mail := &email.Email{
		From:    "萌媛星 <KiriRainCat@163.com>",
		To:      []string{userMail},
		Subject: "验证码 (Validation Code)",
		Text:    []byte("验证码将在5分钟后失效，请在时效内进行验证\n Validation code will expire in 5 minutes, please validate with in time limit\n\n" + strconv.Itoa(code)),
	}

	return mail.Send(
		config.Config.SMTP.Host+":"+strconv.Itoa(config.Config.SMTP.Port),
		smtp.PlainAuth("", config.Config.SMTP.Mail, config.Config.SMTP.Key, config.Config.SMTP.Host),
	)
}

func (s *UserService) Register(email string, validationCode string, username string, pwd string) error {
	// Verify verification code
	val, err := dao.Redis.Get("vc-" + email).Result()
	if err != nil {
		return errors.New("验证码不存在或过期")
	}

	if val != validationCode {
		return errors.New("验证码错误")
	}

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

	// Remove validation code from Redis
	dao.Redis.Del("vc-" + email)

	return nil
}

func (s *UserService) Login(name string, pwd string) (token string, user *model.User, err error) {
	// Find user from DB
	user, err = s.d.GetUserByUsernameOrEmail(name)
	if err != nil {
		err = errors.New("账户不存在")
		return
	}

	// When pwd does not match
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pwd+config.Config.Server.EncryptSalt)) != nil {
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

func (s *UserService) ChangePassword(uid int, pwd string, newPwd string) error {
	// Find user
	user, err := s.d.GetUserByID(uid)
	if err != nil {
		return errors.New("internalErr")
	}

	// Check whether old password match
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pwd+config.Config.Server.EncryptSalt)) != nil {
		return errors.New("旧密码不匹配")
	}

	// Update user password
	encodedPwd, err := bcrypt.GenerateFromPassword([]byte(newPwd+config.Config.Server.EncryptSalt), bcrypt.MinCost)
	if err != nil {
		return errors.New("internalErr")
	}

	err = s.d.UpdatePassword(uid, string(encodedPwd))
	if err != nil {
		return errors.New("internalErr")
	}

	return nil
}

func (s *UserService) ChangeCfbp(uid int, cfbp string) error {
	return s.d.UpdateCfbp(uid, cfbp)
}

func (s *UserService) CompleteInfo(uid int, account string, pwd string) error {
	// Check if user exist or already has Jupiter info
	_, err := s.d.GetUserByID(uid)
	_, exist := dao.Jupiter.GetDataByUID(uid)
	if err != nil || exist == nil {
		return errors.New("internalErr")
	}

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
