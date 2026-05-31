package service_test

import (
	"context"
	"os"
	"testing"

	"helpdesk-api/internal/model"
	"helpdesk-api/internal/service"
	"helpdesk-api/pkg/middleware"

	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) CreateUser(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepo) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	args := m.Called(ctx, email)
	u, _ := args.Get(0).(*model.User)
	return u, args.Error(1)
}

func (m *MockUserRepo) GetUserByID(ctx context.Context, id int64) (*model.User, error) {
	args := m.Called(ctx, id)
	u, _ := args.Get(0).(*model.User)
	return u, args.Error(1)
}

func TestLogin_UserNotFound(t *testing.T) {
	repo := new(MockUserRepo)
	repo.On("GetUserByEmail", mock.Anything, "ghost@test.com").Return(nil, nil)

	svc := service.NewUserService(repo)

	token, err := svc.Login(context.Background(), "ghost@test.com", "whatever")

	assert.Empty(t, token)
	assert.ErrorIs(t, err, service.ErrInvalidCredentials)
	repo.AssertExpectations(t)
}

func TestLogin_WrongPassword(t *testing.T) {
	repo := new(MockUserRepo)

	hash, _ := bcrypt.GenerateFromPassword([]byte("correct-pass"), bcrypt.DefaultCost)
	repo.On("GetUserByEmail", mock.Anything, "user@test.com").
		Return(&model.User{ID: 1, Email: "user@test.com", Password: string(hash), Role: model.RoleClient}, nil)

	svc := service.NewUserService(repo)

	token, err := svc.Login(context.Background(), "user@test.com", "wrong-pass")

	assert.Empty(t, token)
	assert.ErrorIs(t, err, service.ErrInvalidCredentials)
	repo.AssertExpectations(t)
}

func TestMain(m *testing.M) {
	middleware.LoadJWTSecret("test-secret-at-least-32-chars-long-xx")
	os.Exit(m.Run())
}

func TestLogin_Success(t *testing.T) {
	repo := new(MockUserRepo)
	hash, _ := bcrypt.GenerateFromPassword([]byte("plain-pass"), bcrypt.DefaultCost)
	repo.On("GetUserByEmail", mock.Anything, "user@test.com").
		Return(&model.User{ID: 1, Email: "user@test.com", Password: string(hash), Role: model.RoleClient}, nil)

	svc := service.NewUserService(repo)

	token, err := svc.Login(context.Background(), "user@test.com", "plain-pass")

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	repo.AssertExpectations(t)
}

func TestRegister_Success(t *testing.T) {
	repo := new(MockUserRepo)

	repo.On("CreateUser", mock.Anything, mock.Anything).Return(nil)

	svc := service.NewUserService(repo)
	user := &model.User{Email: "new@test.com", Password: "plain-pass"}
	err := svc.Register(context.Background(), user)

	assert.NoError(t, err)
	assert.NotEqual(t, "plain-pass", user.Password)
	repo.AssertExpectations(t)
}
func TestRegister_DefaultRole(t *testing.T) {
	repo := new(MockUserRepo)

	repo.On("CreateUser", mock.Anything, mock.Anything).Return(nil)

	svc := service.NewUserService(repo)
	user := &model.User{Email: "new@test.com", Password: "plain-pass"}
	err := svc.Register(context.Background(), user)

	assert.NoError(t, err)
	assert.Equal(t, model.RoleClient, user.Role)
	repo.AssertExpectations(t)
}

func TestRegister_DuplicateEmail(t *testing.T) {
	repo := new(MockUserRepo)

	repo.On("CreateUser", mock.Anything, mock.Anything).
		Return(&pq.Error{Code: "23505"})

	svc := service.NewUserService(repo)
	user := &model.User{Email: "new@test.com", Password: "plain-pass"}
	err := svc.Register(context.Background(), user)

	assert.ErrorIs(t, err, service.ErrEmailAlreadyExists)
	repo.AssertExpectations(t)
}
