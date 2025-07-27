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

func TestDeleteUser(t *testing.T) {
	t.Run("delete user - success", deleteUserSuccess)
	t.Run("delete user - empty email", deleteUserEmptyEmail)
	t.Run("delete user - not found", deleteUserNotFound)
	t.Run("delete user - internal error", deleteUserInternalError)
}

func TestGetUserWithEmail(t *testing.T) {
	t.Run("get user with email - success", getUserWithEmailSuccess)
	t.Run("get user with email - email not found", getUserWithEmailEmailNotFound)
	t.Run("get user with email - internal error", getUserWithEmailInternalError)
}

func TestAdminUserPresent(t *testing.T) {
	t.Run("admin user present - admin exists", adminUserPresentAdminExists)
	t.Run("admin user present - no admin", adminUserPresentNoAdmin)
	t.Run("admin user present - internal error", adminUserPresentInternalError)
}

func TestGetUsers(t *testing.T) {
	t.Run("get users - success", getUsersSuccess)
	t.Run("get users - empty list", getUsersEmptyList)
	t.Run("get users - internal error", getUsersInternalError)
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

func deleteUserSuccess(t *testing.T) {
	userService, mocks := setUp(t)

	email := "test@example.com"

	mocks.storageServiceMock.EXPECT().
		DeleteUser(gomock.Any(), email).
		Return(nil)

	err := userService.DeleteUser(context.Background(), email)

	require.NoError(t, err)
}

func deleteUserEmptyEmail(t *testing.T) {
	userService, _ := setUp(t)

	email := ""

	expectedError := errors.NewBadRequestError("email cannot be empty")
	err := userService.DeleteUser(context.Background(), email)

	require.Error(t, err)
	require.Equal(t, expectedError.Error(), err.Error())
}

func deleteUserNotFound(t *testing.T) {
	userService, mocks := setUp(t)

	email := "nonexistent@example.com"

	mocks.loggerMocks.EXPECT().
		Errorf(gomock.Any(), gomock.Any())

	mocks.storageServiceMock.EXPECT().
		DeleteUser(gomock.Any(), email).
		Return(storage.ErrNotFound)

	expectedError := errors.NewNotFoundError(fmt.Sprintf("user with email %s not found", email))
	err := userService.DeleteUser(context.Background(), email)

	require.Error(t, err)
	require.Equal(t, expectedError.Error(), err.Error())
}

func deleteUserInternalError(t *testing.T) {
	userService, mocks := setUp(t)

	email := "test@example.com"
	mockedError := fmt.Errorf("database connection failed")

	mocks.loggerMocks.EXPECT().
		Errorf(gomock.Any(), gomock.Any())

	mocks.storageServiceMock.EXPECT().
		DeleteUser(gomock.Any(), email).
		Return(mockedError)

	expectedError := errors.NewInternalServerError(fmt.Sprintf("failed to delete user with email %s: %v", email, mockedError))
	err := userService.DeleteUser(context.Background(), email)

	require.Error(t, err)
	require.Equal(t, expectedError.Error(), err.Error())
}

func getUserWithEmailSuccess(t *testing.T) {
	userService, mocks := setUp(t)

	email := "test@example.com"
	username := "testuser"
	role := "user"

	userObj := &obj.User{
		Username:          username,
		Email:             email,
		EncryptedPassword: "hashedpassword",
		Role:              role,
	}

	userModel := &model.User{
		Username: username,
		Email:    email,
		Role:     role,
	}

	mocks.storageServiceMock.EXPECT().
		GetUserWithEmail(gomock.Any(), email).
		Return(userObj, nil)

	mocks.translatorMocks.EXPECT().
		ToUserModel(userObj).
		Return(userModel)

	result, err := userService.GetUserWithEmail(context.Background(), email)

	require.NoError(t, err)
	require.Equal(t, userModel, result)
}

func getUserWithEmailEmailNotFound(t *testing.T) {
	userService, mocks := setUp(t)

	email := "test@example.com"

	mocks.loggerMocks.EXPECT().
		Errorf(gomock.Any(), gomock.Any())

	mocks.storageServiceMock.EXPECT().
		GetUserWithEmail(gomock.Any(), email).
		Return(nil, storage.ErrNotFound)

	expectedError := errors.NewNotFoundError(fmt.Sprintf("user with email %s not found", email))
	result, err := userService.GetUserWithEmail(context.Background(), email)

	require.Error(t, err)
	require.Nil(t, result)
	require.Equal(t, expectedError.Error(), err.Error())
}

func getUserWithEmailInternalError(t *testing.T) {
	userService, mocks := setUp(t)

	email := "test@example.com"
	mockedError := fmt.Errorf("database connection failed")

	mocks.loggerMocks.EXPECT().
		Errorf(gomock.Any(), gomock.Any())

	mocks.storageServiceMock.EXPECT().
		GetUserWithEmail(gomock.Any(), email).
		Return(nil, mockedError)

	expectedError := errors.NewInternalServerError(fmt.Sprintf("failed to retrieve user with email %s: %v", email, mockedError))
	result, err := userService.GetUserWithEmail(context.Background(), email)

	require.Error(t, err)
	require.Nil(t, result)
	require.Equal(t, expectedError.Error(), err.Error())
}

func adminUserPresentAdminExists(t *testing.T) {
	userService, mocks := setUp(t)

	adminUser := &obj.User{
		Username:          "admin",
		Email:             "admin@example.com",
		EncryptedPassword: "hashedpassword",
		Role:              AdminUserRole,
	}

	mocks.storageServiceMock.EXPECT().
		GetUserWithRole(gomock.Any(), AdminUserRole).
		Return(adminUser, nil)

	result, err := userService.AdminUserPresent(context.Background())

	require.NoError(t, err)
	require.True(t, result)
}

func adminUserPresentNoAdmin(t *testing.T) {
	userService, mocks := setUp(t)

	mocks.storageServiceMock.EXPECT().
		GetUserWithRole(gomock.Any(), AdminUserRole).
		Return(nil, storage.ErrNotFound)

	result, err := userService.AdminUserPresent(context.Background())

	require.NoError(t, err)
	require.False(t, result)
}

func adminUserPresentInternalError(t *testing.T) {
	userService, mocks := setUp(t)

	mockedError := fmt.Errorf("database connection failed")

	mocks.storageServiceMock.EXPECT().
		GetUserWithRole(gomock.Any(), AdminUserRole).
		Return(nil, mockedError)

	expectedError := errors.NewInternalServerError(fmt.Sprintf("failed to check for admin user: %v", mockedError))
	result, err := userService.AdminUserPresent(context.Background())

	require.Error(t, err)
	require.False(t, result)
	require.Equal(t, expectedError.Error(), err.Error())
}

func getUsersSuccess(t *testing.T) {
	userService, mocks := setUp(t)

	userObj1 := &obj.User{
		Username:          "testuser1",
		Email:             "test1@example.com",
		EncryptedPassword: "hashedpassword1",
		Role:              "user",
	}

	userObj2 := &obj.User{
		Username:          "testuser2",
		Email:             "test2@example.com",
		EncryptedPassword: "hashedpassword2",
		Role:              "admin",
	}

	userObjs := []*obj.User{userObj1, userObj2}

	userModel1 := &model.User{
		Username: "testuser1",
		Email:    "test1@example.com",
		Role:     "user",
	}

	userModel2 := &model.User{
		Username: "testuser2",
		Email:    "test2@example.com",
		Role:     "admin",
	}

	userModels := []*model.User{userModel1, userModel2}

	mocks.storageServiceMock.EXPECT().
		GetUsersWithFilter(gomock.Any(), "").
		Return(userObjs, nil)

	mocks.translatorMocks.EXPECT().
		ToUserModels(userObjs).
		Return(userModels)

	result, err := userService.GetUsers(context.Background())

	require.NoError(t, err)
	require.Equal(t, userModels, result)
	require.Len(t, result, 2)
}

func getUsersEmptyList(t *testing.T) {
	userService, mocks := setUp(t)

	var userObjs []*obj.User
	var userModels []*model.User

	mocks.storageServiceMock.EXPECT().
		GetUsersWithFilter(gomock.Any(), "").
		Return(userObjs, nil)

	mocks.translatorMocks.EXPECT().
		ToUserModels(userObjs).
		Return(userModels)

	result, err := userService.GetUsers(context.Background())

	require.NoError(t, err)
	require.Equal(t, userModels, result)
	require.Len(t, result, 0)
}

func getUsersInternalError(t *testing.T) {
	userService, mocks := setUp(t)

	mockedError := fmt.Errorf("database connection failed")

	mocks.loggerMocks.EXPECT().
		Errorf(gomock.Any(), gomock.Any())

	mocks.storageServiceMock.EXPECT().
		GetUsersWithFilter(gomock.Any(), "").
		Return(nil, mockedError)

	expectedError := errors.NewInternalServerError(fmt.Sprintf("failed to retrieve users: %v", mockedError))
	result, err := userService.GetUsers(context.Background())

	require.Error(t, err)
	require.Nil(t, result)
	require.Equal(t, expectedError.Error(), err.Error())
}
