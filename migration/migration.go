package migration

import (
	"github.com/mohnaofal/rest-go-jwt/models"
	"gorm.io/gorm"
)

func Migration(db *gorm.DB) {
	db.AutoMigrate(
		models.User{},
	)
}
