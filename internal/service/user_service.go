package service

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"photoset/internal/database"
	"photoset/internal/domain"
	"photoset/internal/pkg/password"
	"photoset/internal/repository"
	"time"
)

type UserService interface {
	Register(nickname, email, password string) (*domain.User, error)
	Login(email, password string) (*domain.User, error)
	GetProfile(userID uint) (*domain.User, error)
	UpdateProfile(userID uint, nickname, bio, avatar, ipLocation string) (*domain.User, error)
	ChangePassword(userID uint, oldPassword, newPassword string) error
	ResetPassword(userID uint, newPassword string) error
	RequestPasswordReset(email string) (token string, err error)
	ResetPasswordByToken(token, newPassword string) error
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) Register(nickname, email, pwd string) (*domain.User, error) {
	existing, _ := s.userRepo.FindByEmail(email)
	if existing != nil {
		return nil, errors.New("email already exists")
	}

	hashedPassword, err := password.Hash(pwd)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Nickname:     nickname,
		Email:        email,
		PasswordHash: hashedPassword,
		Role:         domain.RoleUser,
		Status:       1,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) Login(email, pwd string) (*domain.User, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if user == nil {
		return nil, errors.New("invalid email or password")
	}

	if !password.Check(pwd, user.PasswordHash) {
		return nil, errors.New("invalid email or password")
	}

	if user.Status != 1 {
		return nil, errors.New("account is disabled")
	}

	now := time.Now()
	user.LastLoginAt = &now
	s.userRepo.Update(user)

	return user, nil
}

func (s *userService) GetProfile(userID uint) (*domain.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

// UpdateProfile 更新用户资料
func (s *userService) UpdateProfile(userID uint, nickname, bio, avatar, ipLocation string) (*domain.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	if nickname != "" {
		user.Nickname = nickname
	}
	user.Bio = bio
	user.Avatar = avatar
	if ipLocation != "" {
		user.IPLocation = ipLocation
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

// ChangePassword 用户修改自己的密码（需要验证旧密码）
func (s *userService) ChangePassword(userID uint, oldPassword, newPassword string) error {
	if len(newPassword) < 6 {
		return errors.New("新密码长度不能少于6位")
	}

	user, err := s.userRepo.FindByID(userID)
	if err != nil || user == nil {
		return errors.New("用户不存在")
	}

	if !password.Check(oldPassword, user.PasswordHash) {
		return errors.New("原密码错误")
	}

	hashedPassword, err := password.Hash(newPassword)
	if err != nil {
		return err
	}

	user.PasswordHash = hashedPassword
	return s.userRepo.Update(user)
}

// ResetPassword 管理员重置用户密码（不需要旧密码）
func (s *userService) ResetPassword(userID uint, newPassword string) error {
	if len(newPassword) < 6 {
		return errors.New("新密码长度不能少于6位")
	}

	user, err := s.userRepo.FindByID(userID)
	if err != nil || user == nil {
		return errors.New("用户不存在")
	}

	hashedPassword, err := password.Hash(newPassword)
	if err != nil {
		return err
	}

	user.PasswordHash = hashedPassword
	return s.userRepo.Update(user)
}

// generateToken 生成安全的随机 token
func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// RequestPasswordReset 请求密码重置（生成 token，发送邮件由 handler 负责调用）
func (s *userService) RequestPasswordReset(email string) (string, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errors.New("该邮箱未注册")
		}
		return "", err
	}
	if user == nil {
		return "", errors.New("该邮箱未注册")
	}

	if user.Status != 1 {
		return "", errors.New("账号已被禁用，请联系管理员")
	}

	// 生成 token
	token, err := generateToken()
	if err != nil {
		return "", errors.New("生成重置令牌失败")
	}

	db := database.GetMySQL()

	// 将之前未使用的 token 全部标记为已使用（防止重复发送）
	db.Model(&domain.PasswordResetToken{}).
		Where("email = ? AND used = ? AND expire > ?", email, false, time.Now()).
		Update("used", true)

	// 创建新 token（30 分钟有效）
	resetToken := &domain.PasswordResetToken{
		UserID: user.ID,
		Email:  email,
		Token:  token,
		Used:   false,
		Expire: time.Now().Add(30 * time.Minute),
	}
	if err := db.Create(resetToken).Error; err != nil {
		return "", errors.New("保存重置令牌失败")
	}

	return token, nil
}

// ResetPasswordByToken 通过 token 重置密码
func (s *userService) ResetPasswordByToken(token, newPassword string) error {
	if len(newPassword) < 6 {
		return errors.New("新密码长度不能少于6位")
	}

	db := database.GetMySQL()
	var resetToken domain.PasswordResetToken
	if err := db.Where("token = ? AND used = ? AND expire > ?", token, false, time.Now()).First(&resetToken).Error; err != nil {
		return errors.New("重置链接无效或已过期")
	}

	// 标记 token 为已使用
	db.Model(&resetToken).Update("used", true)

	// 更新用户密码
	hashedPassword, err := password.Hash(newPassword)
	if err != nil {
		return err
	}

	if err := db.Model(&domain.User{}).Where("id = ?", resetToken.UserID).Update("password_hash", hashedPassword).Error; err != nil {
		return errors.New("密码重置失败")
	}

	return nil
}

