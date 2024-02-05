package handler

import (
	"net/http"
	"time"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/models"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/SawitProRecruitment/UserService/util"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type JwtCustomClaims struct {
	ID int `json:"id"`
	jwt.RegisteredClaims
}

// UserHandler struct
type UserHandler struct {
	UserRepo repository.UserRepository
}

// NewUserHandler create new user handler
func NewUserHandler(userRepo repository.UserRepository) *UserHandler {
	return &UserHandler{UserRepo: userRepo}
}

// Register handler for user registration
func (h *UserHandler) Register(c echo.Context) error {
	var input generated.RegisterRequest
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: "fail to bind input, it might be bad request",
		})
	}
	if err := c.Validate(input); err != nil {
		return err
	}
	if input.Phone[:3] != "+62" {
		return c.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: "phone number must start with +62",
		})
	}

	if !util.ValidatePassword(input.Password) {
		return c.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: "password must contains at least 1 uppercase, 1 number, and 1 special character",
		})
	}

	salt := util.GenerateSalt()
	hashedPassword := util.HashPassword(input.Password, salt)
	existingUser, err := h.UserRepo.FindByPhone(input.Phone)
	if err != nil && err.Error() != "record not found" {
		return c.JSON(http.StatusInternalServerError, generated.ErrorResponse{
			Message: err.Error(),
		})
	}

	if existingUser != nil {
		return c.JSON(http.StatusConflict, generated.ErrorResponse{
			Message: "phone number already registered",
		})
	}

	user := &models.User{
		PhoneNumber: input.Phone,
		Password:    hashedPassword,
		Fullname:    input.Fullname,
		SaltToken:   salt,
	}

	err = h.UserRepo.Create(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, generated.ErrorResponse{
			Message: err.Error(),
		})
	}

	user, err = h.UserRepo.FindByPhone(input.Phone)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, generated.ErrorResponse{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, generated.RegisterResponse{
		Id: user.ID,
	})
}

// Login hander for user login
func (h *UserHandler) Login(c echo.Context) error {
	var input generated.LoginRequest
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: "fail to bind input, it might be bad request",
		})
	}
	if err := c.Validate(input); err != nil {
		return err
	}
	if input.Phone[:3] != "+62" {
		return c.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: "phone number must start with +62",
		})
	}

	user, err := h.UserRepo.FindByPhone(input.Phone)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, generated.ErrorResponse{
			Message: err.Error(),
		})
	}

	hashedPassword := util.HashPassword(input.Password, user.SaltToken)
	if user.Password != hashedPassword {
		return c.JSON(http.StatusUnauthorized, generated.ErrorResponse{
			Message: "invalid phone or password",
		})
	}

	claims := &JwtCustomClaims{
		user.ID,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, generated.LoginResponse{
		Id:    user.ID,
		Token: t,
	})
}

// Profile handler for user profile
func (h *UserHandler) Profile(c echo.Context) error {
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(*JwtCustomClaims)

	user, err := h.UserRepo.FindByID(claims.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, generated.ErrorResponse{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, generated.ProfileResponse{
		Fullname: user.Fullname,
		Phone:    user.PhoneNumber,
	})
}

// UpdateProfile handler for updating user profile
func (h *UserHandler) UpdateProfile(c echo.Context) error {
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(*JwtCustomClaims)

	var input generated.UpdateProfileRequest
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: "fail to bind input, it might be bad request",
		})
	}
	if err := c.Validate(input); err != nil {
		return err
	}
	if input.Phone[:3] != "+62" {
		return c.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: "phone number must start with +62",
		})
	}

	user, err := h.UserRepo.FindByID(claims.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, generated.ErrorResponse{
			Message: err.Error(),
		})
	}

	existingUser, err := h.UserRepo.FindByPhone(input.Phone)
	if err != nil && err.Error() != "record not found" {
		return c.JSON(http.StatusInternalServerError, generated.ErrorResponse{
			Message: err.Error(),
		})
	}

	if existingUser != nil && existingUser.ID != user.ID {
		return c.JSON(http.StatusConflict, generated.ErrorResponse{
			Message: "phone number already registered",
		})
	}

	user.PhoneNumber = input.Phone
	if input.Fullname != "" {
		user.Fullname = input.Fullname
	}

	err = h.UserRepo.Update(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, generated.ErrorResponse{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, generated.ProfileResponse{
		Fullname: user.Fullname,
		Phone:    user.PhoneNumber,
	})
}
