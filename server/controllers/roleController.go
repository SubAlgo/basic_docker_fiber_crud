package controllers

import (
	"basicCRUD/models"
	"basicCRUD/service"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type RoleController struct {
	DB *gorm.DB
}

type roleForm struct {
	Name string `json:"name" validate:"required"`
	Desc string `json:"desc"`
}

type roleResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Desc string `json:"desc"`
}

func (r *RoleController) Create(c *fiber.Ctx) error {
	var form roleForm
	if err := c.BodyParser(&form); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{
			"error": err.Error(),
		})
	}

	// validate input
	errs := service.ValidateStruct(&form)
	if errs != nil {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"error": errs,
		})
	}

	var roleModel models.Roles
	copier.Copy(&roleModel, &form)

	if err := r.DB.Create(&roleModel).Error; err != nil {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"error": errs,
		})
	}

	serializedRole := roleResponse{}
	copier.Copy(&serializedRole, &roleModel)
	return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
		"role": serializedRole,
	})
}

// FindAll
func (r *RoleController) FindAll(c *fiber.Ctx) error {
	var roles []models.Roles

	pagination := pagination{
		ctx:     c,
		query:   r.DB.Order("name asc"),
		records: &roles,
	}
	paging := pagination.paginate()

	// set response
	var serializedRoles []roleResponse

	copier.Copy(&serializedRoles, &roles)

	return c.Status(http.StatusOK).JSON(&fiber.Map{
		"roles": rolePaging{
			Items:  serializedRoles,
			Paging: paging,
		},
	})
}

// FindOne
func (r *RoleController) FindOne(c *fiber.Ctx) error {
	role, err := r.findRoleByID(c)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(&fiber.Map{
			"error": "role not found",
		})
	}

	serializedRole := roleResponse{}
	copier.Copy(&serializedRole, &role)
	return c.Status(http.StatusOK).JSON(&fiber.Map{
		"user": serializedRole,
	})
}

func (r *RoleController) Update(c *fiber.Ctx) error {
	var form roleForm
	var err error
	roleModel, err := r.findRoleByID(c)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"errors": err,
		})
	}

	if err = c.BodyParser(&form); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{
			"error": err.Error(),
		})
	}

	if errs := service.ValidateStruct(&form); errs != nil {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"errors": errs,
		})
	}

	copier.Copy(&roleModel, &form)
	if err = r.DB.Model(&roleModel).Select("id").Updates(&roleModel).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(&fiber.Map{
			"error": "update: " + err.Error(),
		})
	}

	var serializedRole roleResponse
	copier.Copy(&serializedRole, roleModel)

	return c.Status(http.StatusOK).JSON(&fiber.Map{
		"role": serializedRole,
	})

}

func (r *RoleController) Delete(c *fiber.Ctx) error {
	return c.Status(http.StatusOK).JSON(&fiber.Map{
		"role": "delete",
	})
}

// internal function
func (r *RoleController) findRoleByID(c *fiber.Ctx) (*models.Roles, error) {
	var role models.Roles
	id := c.Params("id")

	if err := r.DB.Unscoped().First(&role, id).Error; err != nil {
		return nil, err
	}

	return &role, nil
}
