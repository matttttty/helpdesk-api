package service

import (
	"context"
	"errors"

	//"database/sql"
	"golang.org/x/crypto/bcrypt"
	"helpdesk-api/internal/model"
	"helpdesk-api/internal/repository"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(ctx context.Context, user *model.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	if user.Role == "" {
		user.Role = model.RoleClient
	}
	if err = s.repo.CreateUser(ctx, user); err != nil {
		return err
	}

	return nil
}

func (s *UserService) Login(ctx context.Context, email, password string) (string, error) {

	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid password")
	}

	return "token", nil
}
