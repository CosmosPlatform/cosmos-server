package user

import (
	"context"
	"cosmos-server/pkg/errors"
	logMock "cosmos-server/pkg/log/mock"
	"cosmos-server/pkg/model"
	userMock "cosmos-server/pkg/services/user/mock"
	"cosmos-server/pkg/storage"
	storageMock "cosmos-server/pkg/storage/mock"
	"cosmos-server/pkg/storage/obj"
	"fmt"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestRegisterUser(t *testing.T) {
	t.Run("register user - success", registerUserSuccess)
	t.Run("register admin user - success", registerAdminUserSuccess)
	t.Run("register user - invalid role", registerUserInvalidRole)
	t.Run("register user - duplicate email", registerUserDuplicateEmail)
	t.Run("register user - internal error", registerUserInternalError)
}

type mocks struct {
	controller         *gomock.Controller
	storageServiceMock *storageMock.MockService
	translatorMocks    *userMock.MockTranslator
	loggerMocks        *logMock.MockLogger
}

func setUp(t *testing.T) (Service, *mocks) {
	ctrl := gomock.NewController(t)

	mocks := &mocks{
		controller:         ctrl,
		storageServiceMock: storageMock.NewMockService(ctrl),
		translatorMocks:    userMock.NewMockTranslator(ctrl),
		loggerMocks:        logMock.NewMockLogger(ctrl),
	}

	userService := NewUserService(mocks.storageServiceMock, mocks.translatorMocks, mocks.loggerMocks)

	return userService, mocks
}

func registerUserSuccess(t *testing.T) {
	userService, mocks := setUp(t)

	username := "testuser"
	email := "test@example.com"
	password := "securepassword"
	role := "user"

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	userModel := &model.User{
		Username: username,
		Email:    email,
		Role:     role,
	}

	userObj := &obj.User{
		Username:          username,
		Email:             email,
		EncryptedPassword: string(hashedPassword),
		Role:              role,
	}

	mocks.storageServiceMock.EXPECT().
		GetUserWithEmail(gomock.Any(), email).
		Return(nil, storage.ErrNotFound)

	mocks.translatorMocks.EXPECT().
		ToUserObj(userModel, gomock.Any()). // We put the Any because the hashed password has the random salt.
		Return(userObj)

	mocks.storageServiceMock.EXPECT().
		InsertUser(gomock.Any(), userObj).
		Return(nil)

	err = userService.RegisterUser(context.Background(), username, email, password, role)

	require.NoError(t, err)
}

func registerAdminUserSuccess(t *testing.T) {
	userService, mocks := setUp(t)

	username := "testuser"
	email := "test@example.com"
	password := "securepassword"
	role := "admin"

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	userModel := &model.User{
		Username: username,
		Email:    email,
		Role:     role,
	}

	userObj := &obj.User{
		Username:          username,
		Email:             email,
		EncryptedPassword: string(hashedPassword),
		Role:              role,
	}

	mocks.storageServiceMock.EXPECT().
		GetUserWithEmail(gomock.Any(), email).
		Return(nil, storage.ErrNotFound)

	mocks.translatorMocks.EXPECT().
		ToUserObj(userModel, gomock.Any()). // We put the Any because the hashed password has the random salt.
		Return(userObj)

	mocks.storageServiceMock.EXPECT().
		InsertUser(gomock.Any(), userObj).
		Return(nil)

	err = userService.RegisterUser(context.Background(), username, email, password, role)

	require.NoError(t, err)
}

func registerUserInvalidRole(t *testing.T) {
	userService, _ := setUp(t)

	username := "testuser"
	email := "test@example.com"
	password := "securepassword"
	role := "unexisting_role"

	err := userService.RegisterUser(context.Background(), username, email, password, role)

	require.Error(t, err)
}

func registerUserDuplicateEmail(t *testing.T) {
	userService, mocks := setUp(t)

	username := "testuser"
	email := "test@example.com"
	password := "securepassword"
	role := "admin"

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	userObj := &obj.User{
		Username:          username,
		Email:             email,
		EncryptedPassword: string(hashedPassword),
		Role:              role,
	}

	mocks.loggerMocks.EXPECT().
		Errorf(gomock.Any(), gomock.Any())

	mocks.storageServiceMock.EXPECT().
		GetUserWithEmail(gomock.Any(), email).
		Return(userObj, nil)

	expectedError := errors.NewBadRequestError(fmt.Sprintf("user with email %s already exists", email))

	err = userService.RegisterUser(context.Background(), username, email, password, role)

	require.Error(t, err)
	require.Equal(t, expectedError.Error(), err.Error())
}

func registerUserInternalError(t *testing.T) {
	userService, mocks := setUp(t)

	username := "testuser"
	email := "test@example.com"
	password := "securepassword"
	role := "admin"

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	userModel := &model.User{
		Username: username,
		Email:    email,
		Role:     role,
	}

	userObj := &obj.User{
		Username:          username,
		Email:             email,
		EncryptedPassword: string(hashedPassword),
		Role:              role,
	}

	mocks.storageServiceMock.EXPECT().
		GetUserWithEmail(gomock.Any(), email).
		Return(nil, storage.ErrNotFound)

	mocks.translatorMocks.EXPECT().
		ToUserObj(userModel, gomock.Any()). // We put the Any because the hashed password has the random salt.
		Return(userObj)

	mockedError := fmt.Errorf("internal error")

	mocks.storageServiceMock.EXPECT().
		InsertUser(gomock.Any(), userObj).
		Return(mockedError)

	expectedError := errors.NewInternalServerError(fmt.Sprintf("failed to insert user into storage: %v", mockedError))
	err = userService.RegisterUser(context.Background(), username, email, password, role)

	require.Error(t, err)
	require.Equal(t, expectedError.Error(), err.Error())
}
