package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mohnaofal/rest-go-jwt/middleware"
	"github.com/mohnaofal/rest-go-jwt/models"
	"github.com/mohnaofal/rest-go-jwt/repository"
	"github.com/mohnaofal/rest-go-jwt/utils"
	"gorm.io/gorm"
)

type UserController struct {
	userRepo repository.UserRepository
}

func NewUserController(userRepo repository.UserRepository) UserController {
	return UserController{userRepo: userRepo}
}

func (h *UserController) Apply(c *echo.Group) {
	c.POST("/registry", h.Registry)
	c.POST("/login", h.Login)
	c.GET("/profile", h.Profile, middleware.VerifyJwt())
}

func (h *UserController) Registry(c echo.Context) error {
	ctx := c.Request().Context()

	form := new(models.User)
	if err := c.Bind(form); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	if err := c.Validate(form); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	user, err := h.userRepo.GetByUsername(ctx, &models.User{Username: form.Username})
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	if user != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"success": false,
			"message": "user already exist",
		})
	}

	passwordHash, err := utils.HashPassword(form.Password)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	data := &models.User{
		Username: form.Username,
		Password: passwordHash,
		Name:     form.Name,
	}

	data, err = h.userRepo.Create(ctx, data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	data.Password = ``

	return c.JSON(http.StatusCreated, echo.Map{
		"success": true,
		"message": `Registry success`,
		"data":    data,
	})
}

func (h *UserController) Login(c echo.Context) error {
	ctx := c.Request().Context()

	form := new(models.LoginRequest)
	if err := c.Bind(form); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	if err := c.Validate(form); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	user, err := h.userRepo.GetByUsername(ctx, &models.User{Username: form.Username})
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	if user == nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"success": false,
			"message": "user not found",
		})
	}

	if ok := utils.CheckPasswordHash(user.Password, form.Password); !ok {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"success": false,
			"message": "wrong password",
		})
	}

	token := middleware.GenerateJwt(user)
	if token == `` {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"success": false,
			"message": "failed generate token",
		})
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"success": true,
		"message": `Registry success`,
		"data": echo.Map{
			"token": token,
		},
	})
}

func (h *UserController) Profile(c echo.Context) error {
	ctx := c.Request().Context()

	userID := c.Get("user_id").(int)

	user, err := h.userRepo.Get(ctx, &models.User{Model: gorm.Model{ID: uint(userID)}})
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	if user == nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"success": false,
			"message": "user not found",
		})
	}

	user.Password = ``

	return c.JSON(http.StatusCreated, echo.Map{
		"success": true,
		"message": `Registry success`,
		"data":    user,
	})
}
