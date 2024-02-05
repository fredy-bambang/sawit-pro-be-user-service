package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jinzhu/gorm"
	"github.com/rakyll/statik/fs"

	"github.com/SawitProRecruitment/UserService/handler"
	"github.com/SawitProRecruitment/UserService/models"
	"github.com/SawitProRecruitment/UserService/repository"
	_ "github.com/SawitProRecruitment/UserService/statik"
	"github.com/go-playground/validator/v10"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func main() {
	startServer()
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func startServer() {
	// dsn := "host=localhost user=gorm password=gorm dbname=gorm port=5432 sslmode=disable"
	// "postgres://postgres:postgres@postgres:5432/postgres?sslmode=disable"
	// dsn := "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" //os.Getenv("DATABASE_URL")
	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open("postgres", dsn)
	if err != nil {
		panic("Failed to connect to database")
	}

	// Auto Migrate PostgreSQL
	db.AutoMigrate(&models.User{})

	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	e.HTTPErrorHandler = handleHTTPError

	// Initialize repositories
	userRepo := repository.NewPgUserRepository(db)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userRepo)

	// create docs for swagger handler in echo
	statikFS, err := fs.New()
	if err != nil {
		panic(err)
	}
	// Serve the Swagger UI at the route /swaggerui
	e.GET("/swaggerui/*", echo.WrapHandler(http.StripPrefix("/swaggerui/", http.FileServer(statikFS))))
	e.POST("/register", userHandler.Register)
	e.POST("/login", userHandler.Login)

	// Restricted group
	r := e.Group("/profile")

	// Configure middleware with the custom claims type
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(handler.JwtCustomClaims)
		},
		SigningKey: []byte("secret"),
	}
	r.Use(echojwt.WithConfig(config))
	r.GET("", userHandler.Profile)
	r.PATCH("", userHandler.UpdateProfile)

	e.Logger.Fatal(e.Start(":1323"))
}

func handleHTTPError(err error, c echo.Context) {
	report, ok := err.(*echo.HTTPError)
	if !ok {
		report = echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if castedObject, ok := err.(validator.ValidationErrors); ok {
		errorMessage := []string{}
		for _, err := range castedObject {
			switch err.Tag() {
			case "required":
				errorMessage = append(errorMessage, fmt.Sprintf("%s is required", err.Field()))
			case "email":
				errorMessage = append(errorMessage, fmt.Sprintf("%s is not valid email", err.Field()))
			case "gte":
				errorMessage = append(errorMessage, fmt.Sprintf("%s value must be greater than %s", err.Field(), err.Param()))
			case "lte":
				errorMessage = append(errorMessage, fmt.Sprintf("%s value must be lower than %s", err.Field(), err.Param()))
			case "min":
				errorMessage = append(errorMessage, fmt.Sprintf("%s must have length or number at least %s", err.Field(), err.Param()))
			case "max":
				errorMessage = append(errorMessage, fmt.Sprintf("%s must have length or number with maximum %s", err.Field(), err.Param()))
			}
		}
		report.Message = strings.Join(errorMessage, ", ")
	}
	c.Logger().Error(report)
	c.JSON(report.Code, report)
}
