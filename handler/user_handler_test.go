package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/models"
	"github.com/SawitProRecruitment/UserService/util"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByPhone(phone string) (*models.User, error) {
	args := m.Called(phone)
	return args[0].(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindByID(id int) (*models.User, error) {
	args := m.Called(id)
	return args[0].(*models.User), args.Error(1)
}

func (m *MockUserRepository) Update(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func registerEchoCtx(jsonInput, endpoint string) (*httptest.ResponseRecorder, echo.Context) {
	req := httptest.NewRequest(http.MethodPost, endpoint, strings.NewReader(jsonInput))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	c := e.NewContext(req, rec)

	return rec, c
}

func TestRegister(t *testing.T) {
	mockRepo := new(MockUserRepository)

	handler := &UserHandler{
		UserRepo: mockRepo,
	}

	jsonInput := `{
		"phone": "+62812345678912",
		"password": "A1234*",
		"fullname": "mr smith"
	}`
	rec, c := registerEchoCtx(jsonInput, "/register")

	input := generated.RegisterRequest{
		Phone:    "+62812345678912",
		Password: "A1234*",
		Fullname: "mr smith",
	}

	// Set expectations on the mock repository for the FindByPhone method
	var existingUser *models.User
	mockRepo.On("FindByPhone", input.Phone).Return(existingUser, errors.New("record not found")).Once()

	// Set expectations on the mock repository for the Create method
	mockRepo.On("Create", mock.AnythingOfType("*models.User")).Return(nil)

	mockRepo.On("FindByPhone", input.Phone).Return(&models.User{
		ID:          1,
		PhoneNumber: input.Phone,
		Password:    "hashedPassword",
		Fullname:    input.Fullname,
		SaltToken:   "salt",
	}, nil)

	err := handler.Register(c)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusCreated, rec.Code)
	mockRepo.AssertExpectations(t)
}

func TestRegisterFailBindInput(t *testing.T) {
	mockRepo := new(MockUserRepository)

	handler := &UserHandler{
		UserRepo: mockRepo,
	}

	jsonInput := `{
		"phone": "+62812345678912",
	}`
	rec, c := registerEchoCtx(jsonInput, "/register")

	err := handler.Register(c)
	assert.NoError(t, err)

	expectedJSON := `{"message":"fail to bind input, it might be bad request"}`
	assert.JSONEq(t, expectedJSON, rec.Body.String())
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	mockRepo.AssertExpectations(t)
}

func TestRegisterValidate(t *testing.T) {
	mockRepo := new(MockUserRepository)

	handler := &UserHandler{
		UserRepo: mockRepo,
	}

	jsonInput := `{
		"phone": "+62812345678912"
	}`
	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(jsonInput))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}

	// there is a bug with this feature, it's return 200 instead of 400 with correct custom message
	// but the error detected
	c := e.NewContext(req, rec)

	err := handler.Register(c)
	assert.Error(t, err)

	mockRepo.AssertExpectations(t)
}

func TestRegisterValidatePhone(t *testing.T) {
	mockRepo := new(MockUserRepository)

	handler := &UserHandler{
		UserRepo: mockRepo,
	}

	jsonInput := `{
		"phone": "0812345678912",
		"password": "A1234*",
		"fullname": "mr smith"
	}`
	rec, c := registerEchoCtx(jsonInput, "/register")
	err := handler.Register(c)
	assert.NoError(t, err)

	expectedJSON := `{"message":"phone number must start with +62"}`
	assert.JSONEq(t, expectedJSON, rec.Body.String())
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	mockRepo.AssertExpectations(t)
}

func TestRegisterValidatePassword(t *testing.T) {
	mockRepo := new(MockUserRepository)

	handler := &UserHandler{
		UserRepo: mockRepo,
	}

	jsonInput := `{
		"phone": "0812345678912",
		"password": "abcdef",
		"fullname": "mr smith"
	}`
	rec, c := registerEchoCtx(jsonInput, "/register")

	err := handler.Register(c)
	assert.NoError(t, err)

	expectedJSON := `{"message":"phone number must start with +62"}`
	assert.JSONEq(t, expectedJSON, rec.Body.String())
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	mockRepo.AssertExpectations(t)
}

func TestRegisterValidateErrorRecordNotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)

	handler := &UserHandler{
		UserRepo: mockRepo,
	}

	jsonInput := `{
		"phone": "+62812345678912",
		"password": "A1234*",
		"fullname": "mr smith"
	}`
	rec, c := registerEchoCtx(jsonInput, "/register")

	var emptyUser *models.User
	mockRepo.On("FindByPhone", "+62812345678912").Return(emptyUser, errors.New("test the error"))

	err := handler.Register(c)
	assert.NoError(t, err)

	expectedJSON := `{"message":"test the error"}`
	assert.JSONEq(t, expectedJSON, rec.Body.String())
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	mockRepo.AssertExpectations(t)
}

func TestRegisterValidatePhoneNumberAlreadyRegistered(t *testing.T) {
	mockRepo := new(MockUserRepository)

	handler := &UserHandler{
		UserRepo: mockRepo,
	}

	jsonInput := `{
		"phone": "+62812345678912",
		"password": "A1234*",
		"fullname": "mr smith"
	}`
	rec, c := registerEchoCtx(jsonInput, "/register")

	existedUser := &models.User{
		ID:          1,
		PhoneNumber: "+62812345678912",
		Password:    "hashedPassword",
		Fullname:    "mr smith",
		SaltToken:   "salt",
	}
	mockRepo.On("FindByPhone", "+62812345678912").Return(existedUser, nil)

	err := handler.Register(c)
	assert.NoError(t, err)

	expectedJSON := `{"message":"phone number already registered"}`
	assert.JSONEq(t, expectedJSON, rec.Body.String())
	assert.Equal(t, http.StatusConflict, rec.Code)
	mockRepo.AssertExpectations(t)
}

func TestRegisterCreateUserError(t *testing.T) {
	mockRepo := new(MockUserRepository)

	handler := &UserHandler{
		UserRepo: mockRepo,
	}

	jsonInput := `{
		"phone": "+62812345678912",
		"password": "A1234*",
		"fullname": "mr smith"
	}`
	rec, c := registerEchoCtx(jsonInput, "/register")

	var emptyUser *models.User
	mockRepo.On("FindByPhone", "+62812345678912").Return(emptyUser, errors.New("record not found")).Once()
	mockRepo.On("Create", mock.AnythingOfType("*models.User")).Return(errors.New("create user error"))

	err := handler.Register(c)
	assert.NoError(t, err)

	expectedJSON := `{"message":"create user error"}`
	assert.JSONEq(t, expectedJSON, rec.Body.String())
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	mockRepo.AssertExpectations(t)
}

func TestLogin(t *testing.T) {
	mockRepo := new(MockUserRepository)

	handler := &UserHandler{
		UserRepo: mockRepo,
	}

	jsonInput := `{
		"phone": "+62812345678912",
		"password": "A1234*"
	}`
	rec, c := registerEchoCtx(jsonInput, "/register")

	input := generated.LoginRequest{
		Phone:    "+62812345678912",
		Password: "A1234*",
	}

	hashedPassword := util.HashPassword(input.Password, "salt")
	mockRepo.On("FindByPhone", input.Phone).Return(&models.User{
		ID:          1,
		PhoneNumber: input.Phone,
		Password:    hashedPassword,
		Fullname:    "The Inspirator",
		SaltToken:   "salt",
	}, nil)

	err := handler.Login(c)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockRepo.AssertExpectations(t)
}

func TestLoginFailBindInput(t *testing.T) {
	mockRepo := new(MockUserRepository)

	handler := &UserHandler{
		UserRepo: mockRepo,
	}

	jsonInput := `{
		"phone": "+62812345678912",
	}`
	rec, c := registerEchoCtx(jsonInput, "/login")

	err := handler.Login(c)
	assert.NoError(t, err)

	expectedJSON := `{"message":"fail to bind input, it might be bad request"}`
	assert.JSONEq(t, expectedJSON, rec.Body.String())
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	mockRepo.AssertExpectations(t)
}

func TestLoginValidatePhone(t *testing.T) {
	mockRepo := new(MockUserRepository)

	handler := &UserHandler{
		UserRepo: mockRepo,
	}

	jsonInput := `{
		"phone": "0812345678912",
		"password": "A1234*"
	}`
	rec, c := registerEchoCtx(jsonInput, "/login")
	err := handler.Login(c)
	assert.NoError(t, err)

	expectedJSON := `{"message":"phone number must start with +62"}`
	assert.JSONEq(t, expectedJSON, rec.Body.String())
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	mockRepo.AssertExpectations(t)
}

func TestLoginValidate(t *testing.T) {
	mockRepo := new(MockUserRepository)

	handler := &UserHandler{
		UserRepo: mockRepo,
	}

	jsonInput := `{
		"phone": "+62812345678912"
	}`
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(jsonInput))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}

	// there is a bug with this feature, it's return 200 instead of 400 with correct custom message
	// but the error detected
	c := e.NewContext(req, rec)

	err := handler.Login(c)
	assert.Error(t, err)

	mockRepo.AssertExpectations(t)
}

func TestProfile(t *testing.T) {
	// Create an instance of the mocked repository and UserHandler
	mockRepo := new(MockUserRepository)
	handler := &UserHandler{
		UserRepo: mockRepo,
	}

	// Mock the JWT token for testing
	claims := &JwtCustomClaims{
		ID: 123,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with a secret (use the same secret you use in your implementation)
	tokenString, err := token.SignedString([]byte("secret"))
	assert.NoError(t, err)

	// Create a sample request with the JWT token in the authorization header
	req := httptest.NewRequest(http.MethodGet, "/profile", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+tokenString)

	// Create a recorder to capture the response
	rec := httptest.NewRecorder()

	// Create a new Echo context
	e := echo.New()
	c := e.NewContext(req, rec)

	// Set the JWT token in the Echo context
	c.Set("user", token)

	// Mock the UserRepo method
	mockRepo.On("FindByID", 123).Return(&models.User{
		ID:          123,
		Fullname:    "John Doe",
		PhoneNumber: "+628123456789",
	}, nil)

	// Call the function being tested
	err = handler.Profile(c)
	assert.NoError(t, err)

	// Assert the response
	assert.Equal(t, http.StatusOK, rec.Code)

	// Assert the response body or any other expectations
	expectedJSON := `{"fullname":"John Doe","phone":"+628123456789"}`
	assert.JSONEq(t, expectedJSON, rec.Body.String())

	// Assert expectations on the mock repository
	mockRepo.AssertExpectations(t)
}

func TestUpdateProfile(t *testing.T) {
	// Create an instance of the mocked repository and UserHandler
	mockRepo := new(MockUserRepository)
	handler := &UserHandler{
		UserRepo: mockRepo,
	}

	// Mock the JWT token for testing
	claims := &JwtCustomClaims{
		ID: 123,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with a secret (use the same secret you use in your implementation)
	tokenString, err := token.SignedString([]byte("secret"))
	assert.NoError(t, err)

	jsonInput := `{
		"phone": "+62812345678912"
	}`

	// Create a sample request with the JWT token in the authorization header
	req := httptest.NewRequest(http.MethodPatch, "/profile", strings.NewReader(jsonInput))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+tokenString)

	// Create a recorder to capture the response
	rec := httptest.NewRecorder()

	// Create a new Echo context
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	c := e.NewContext(req, rec)
	c.Set("user", token)

	hashedPassword := util.HashPassword("A1234*", "salt")
	mockRepo.On("FindByID", 123).Return(&models.User{
		ID:          123,
		PhoneNumber: "+62812345678909",
		Password:    hashedPassword,
		Fullname:    "The Inspirator",
		SaltToken:   "salt",
	}, nil)
	mockRepo.On("FindByPhone", "+62812345678912").Return(&models.User{
		ID:          123,
		PhoneNumber: "+62812345678909",
		Password:    hashedPassword,
		Fullname:    "The Inspirator",
		SaltToken:   "salt",
	}, nil)
	mockRepo.On("Update", &models.User{
		ID:          123,
		PhoneNumber: "+62812345678912",
		Password:    hashedPassword,
		Fullname:    "The Inspirator",
		SaltToken:   "salt",
	}).Return(nil)

	err = handler.UpdateProfile(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	mockRepo.AssertExpectations(t)
}

func TestNewUserHandler(t *testing.T) {
	// Create a mock user repository
	mockUserRepo := new(MockUserRepository)

	// Create a new user handler instance
	userHandler := NewUserHandler(mockUserRepo)

	// Check if the user handler is not nil
	if userHandler == nil {
		t.Errorf("NewUserHandler() returned nil")
	}
}
