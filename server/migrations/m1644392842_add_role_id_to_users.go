package migrations

import (
	"basicCRUD/models"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func m1644392842AddRoleIdToUsers() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "1644392842",
		Migrate: func(tx *gorm.DB) error {
			err := tx.AutoMigrate(&models.Users{})

			//var usersModel []models.Users
			//	tx.Unscoped().Find(&usersModel)
			/*
				for _, user := range usersModel {
					user.RoleID = 2
					tx.Save(&user)
				}
			*/
			return err
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropColumn(&models.Users{}, "role_id")
		},
	}
}
