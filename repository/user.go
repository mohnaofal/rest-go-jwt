package repository

import (
	"context"
	"errors"

	"github.com/mohnaofal/rest-go-jwt/config"
	"github.com/mohnaofal/rest-go-jwt/models"
	"gorm.io/gorm"
)

type userRepository struct {
	cfg *config.Config
	DB  *gorm.DB
}

type UserRepository interface {
	Create(ctx context.Context, data *models.User) (*models.User, error)
	GetByUsername(ctx context.Context, data *models.User) (*models.User, error)
	Get(ctx context.Context, data *models.User) (*models.User, error)
}

func NewUserRepository(cfg *config.Config) UserRepository {
	return &userRepository{
		cfg: cfg,
		DB:  cfg.DB().MysqlGorm(),
	}

}

func (c *userRepository) Create(ctx context.Context, data *models.User) (*models.User, error) {
	if err := c.DB.Create(data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (c *userRepository) GetByUsername(ctx context.Context, data *models.User) (*models.User, error) {
	if err := c.DB.Where("username = ?", data.Username).First(data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return data, nil
}

func (c *userRepository) Get(ctx context.Context, data *models.User) (*models.User, error) {
	if err := c.DB.First(data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return data, nil
}
