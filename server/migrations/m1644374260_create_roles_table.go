package migrations

import (
	"basicCRUD/models"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func m1644374260CreateRolesTable() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "1644374260",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&models.Roles{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("roles")
		},
	}
}
