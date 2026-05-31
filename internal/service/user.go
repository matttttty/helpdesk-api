package service

import (
	"context"
	"errors"

	"github.com/lib/pq"

	"helpdesk-api/internal/model"
	"helpdesk-api/pkg/middleware"

	"golang.org/x/crypto/bcrypt"
)

const pgUniqueViolationCode = "23505"

type UserRepo interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	GetUserByID(ctx context.Context, id int64) (*model.User, error)
}

type UserService struct {
	repo UserRepo
}

func NewUserService(repo UserRepo) *UserService {
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
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == pgUniqueViolationCode {
			return ErrEmailAlreadyExists
		}
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
		return "", ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {

		return "", ErrInvalidCredentials
	}

	token, err := middleware.GenerateToken(user.ID, string(user.Role))
	if err != nil {
		return "", err
	}
	return token, nil
}
