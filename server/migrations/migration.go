package migrations

import (
	"basicCRUD/config"
	"log"

	"github.com/go-gormigrate/gormigrate/v2"
)

/*
	Doc reference: https://pkg.go.dev/github.com/go-gormigrate/gormigrate/v2
*/

func Migrate() {
	db := config.GetDB()

	m := gormigrate.New(
		db,
		gormigrate.DefaultOptions,
		[]*gormigrate.Migration{
			m1644326661CreateUsersTable(),
			m1644374260CreateRolesTable(),
		},
	)

	err := m.Migrate()
	if err != nil {
		log.Fatalf("Could not migrate: %v", err)
	}

	log.Printf("Migration did run successfully")

}
