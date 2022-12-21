package main

import (
	"log"

	"github.com/go-playground/validator"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/mohnaofal/rest-go-jwt/config"
	"github.com/mohnaofal/rest-go-jwt/controller"
	"github.com/mohnaofal/rest-go-jwt/migration"
	"github.com/mohnaofal/rest-go-jwt/repository"
)

// CustomValidator ...
type CustomValidator struct {
	validator *validator.Validate
}

// Validate ...
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// init config
	cfg := config.InitConfig()

	// migration table
	migration.Migration(cfg.DB().MysqlGorm())

	e := echo.New()

	// validator
	e.Validator = &CustomValidator{validator: validator.New()}

	userRepositori := repository.NewUserRepository(cfg)
	userController := controller.NewUserController(userRepositori)
	userGroup := e.Group("v1/user")
	userController.Apply(userGroup)

	e.Logger.Fatal(e.Start(":8080"))
}
