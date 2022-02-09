package seed

import (
	"basicCRUD/config"
	"basicCRUD/migrations"
	"basicCRUD/models"
	"strconv"

	"github.com/bxcodec/faker/v3"
	"github.com/labstack/gommon/log"
)

func Load() {
	db := config.GetDB()

	// Clean Database
	db.Migrator().DropTable("users", "roles", "migrations")
	migrations.Migrate()

	// Add categories
	log.Info("Creating categories...")

	type roleData struct {
		id   uint
		name string
		desc string
	}

	roleList := []roleData{
		{id: 1, name: "admin", desc: "admin role"},
		{id: 2, name: "user", desc: "user role"},
	}

	for _, v := range roleList {
		role := models.Roles{
			Name: v.name,
			Desc: v.desc,
		}

		db.Create(&role)
	}

	// Add articles
	log.Info("Creating articles...")

	numOfUsers := 10
	usersList := make([]models.Users, 0, numOfUsers)

	for i := 1; i <= numOfUsers; i++ {
		user := models.Users{
			Username: faker.Username(),
			Password: "1234",
			Name:     faker.Name(),
			Surname:  faker.LastName(),
			Email:    faker.Email(),
			Image:    "https://source.unsplash.com/random/300x200?" + strconv.Itoa(i),
			RoleID:   2,
		}

		db.Create(&user)
		usersList = append(usersList, user)
		_ = usersList
	}
}
