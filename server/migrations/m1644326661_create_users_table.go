package migrations

import (
	"basicCRUD/models"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func m1644326661CreateUsersTable() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "1644326661",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&models.Users{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("users")
		},
	}
}
