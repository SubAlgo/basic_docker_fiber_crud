package controllers

import (
	"basicCRUD/models"
	"basicCRUD/service"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type UserController struct {
	DB *gorm.DB
}

// Create: create new user data
func (u *UserController) Create(c *fiber.Ctx) error {
	var (
		form CreateUserForm
		err  error
	)

	// parse data from request body to form
	if err := c.BodyParser(&form); err != nil {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"error": "bad request: " + err.Error(),
		})
	}

	// parse image
	form.Image, err = c.FormFile("image")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"error": "bad request: " + err.Error(),
		})
	}

	// validate input
	errs := service.ValidateStruct(&form)
	if errs != nil {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"error": errs,
		})
	}

	// copy ค่าจาก formInput -> userModel
	var userModel models.Users
	copier.Copy(&userModel, &form)

	// save data to db (save data ก่อนเพื่อเอา id มาสร้าง folder เก็บรูป)
	if err = u.DB.Create(&userModel).Error; err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{
			"error": err.Error(),
		})
	}

	// save image
	u.setUserImage(c, &userModel)

	// set response data
	serializedUser := userResponse{}
	copier.Copy(&serializedUser, &userModel)

	// return response
	return c.Status(http.StatusCreated).JSON(&fiber.Map{
		"user": serializedUser,
	})
}

// FindAll: search all user data
func (u *UserController) FindAll(c *fiber.Ctx) error {
	var users []models.Users

	// /users => {limit: 10, page: 1}
	// /users?limit=5 => {limit: 5, page: 1}
	// /users?limit=5?page=2 => {limit:5, page: 2}

	pagination := pagination{
		ctx:     c,
		query:   u.DB.Order("id desc"), //u.DB.Model(&users), //,
		records: &users,
	}
	paging := pagination.paginate() //pagingResource(c, u.DB.Order("id desc"), &users)

	var serializedUsers []userResponse

	copier.Copy(&serializedUsers, &users)

	return c.Status(http.StatusOK).JSON(&fiber.Map{
		"users": userPaging{
			Items:  serializedUsers,
			Paging: paging,
		},
	})
}

// FindOne: select one user by id
func (u *UserController) FindOne(c *fiber.Ctx) error {
	user, err := u.findUserByID(c)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(&fiber.Map{
			"error": "user not found",
		})
	}

	serializedUser := userResponse{}
	copier.Copy(&serializedUser, &user)

	return c.Status(http.StatusOK).JSON(&fiber.Map{
		"user": serializedUser,
	})

}

// Update: update user data by id
func (u *UserController) Update(c *fiber.Ctx) error {
	var form updateUserForm
	var err error
	// parse data
	if err := c.BodyParser(&form); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{
			"error": err.Error(),
		})
	}

	form.Image, err = c.FormFile("image")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"error": "image input: " + err.Error(),
		})
	}

	user, err := u.findUserByID(c)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(&fiber.Map{
			"error": err.Error(),
		})
	}
	copier.Copy(&user, &form)
	if err := u.DB.Model(&user).Select("id").Updates(&user).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(&fiber.Map{
			"error": "update: " + err.Error(),
		})
	}

	if err := u.setUserImage(c, user); err != nil {
		return c.Status(http.StatusNotFound).JSON(&fiber.Map{
			"error": "upload image: " + err.Error(),
		})
	}
	var serializedUser userResponse
	copier.Copy(&serializedUser, user)

	return c.Status(http.StatusOK).JSON(&fiber.Map{
		"user": serializedUser,
	})
}

func (u *UserController) Delete(c *fiber.Ctx) error {
	// 1. ค้นหา user by id
	user, err := u.findUserByID(c)
	var form deleteUserForm
	var disable bool
	if err := c.BodyParser(&form); err != nil {
		disable = false
	} else {
		disable = form.Disable
	}

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"error": err.Error(),
		})
	}

	// 2. ลบข้อมูบ
	// Unscoped().
	if disable {
		if err := u.DB.Delete(&user).Error; err != nil {
			return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
				"error": "delete user: " + err.Error(),
			})
		}
		return c.Status(http.StatusOK).JSON(&fiber.Map{
			"message": "disable user success",
		})
	}

	if err := u.DB.Unscoped().Delete(&user).Error; err != nil {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"error": "delete user: " + err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "delete user success",
	})

}

// *********************************
// Internal function               *
// *********************************
func (u *UserController) findUserByID(c *fiber.Ctx) (*models.Users, error) {
	var user models.Users

	id := c.Params("id")
	if err := u.DB.Unscoped().First(&user, id).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// setUserImage: func สำหรับทำการบันทึกรูปภาพ
func (u *UserController) setUserImage(c *fiber.Ctx, userModel *models.Users) error {

	// 1. check ว่ามีการ upload file เข้ามาหรือไม่
	file, err := c.FormFile("image")
	if err != nil || file == nil {
		return errors.New("user image not found")
	}

	// 2. ระบุ path สำหรับจัดเก็บไฟล์รูปภาพ
	// http://<host>:<port>/uploads/users/<id>/<image_file>
	host := fmt.Sprintf("%v:%v", os.Getenv("HOST"), os.Getenv("APP_PORT")) // 127.0.0.1:9000

	// 2.1 ถ้า path จัดเก็บรูปภาพไม่เท่ากับค่าว่าง (กรณีเป็นการ update ข้อมูล) ให้กำหนดค่าของ path ใหม่ โดย output จะเท่ากับ host:port
	if userModel.Image != "" {
		// เป้าหมายของขั้นตอนนี้ คือแปลง 127.0.0.1:9000/uploads/users/<id>/image.png -> /uploads/users/<id>/image.png
		// สาเหตุที่ต้องมีขั้นตอนนี้ก็เพื่อว่า กรณีที่เป็นการ update ข้อมูลเคยมีการ upload ไฟล์ไว้แล้ว เราก็ต้องทำการลบไฟล์เก่าก่อน แล้วค่อย upload ไฟล์ใหม่ขึ้นแทน
		userModel.Image = strings.Replace(userModel.Image, host, "", 1)
		pwd, _ := os.Getwd()
		os.Remove(pwd + userModel.Image)
	}

	// 3. กำหนด path ใหม่ upload/users/<userID>
	path := "uploads/users/" + strconv.Itoa(int(userModel.ID))
	os.MkdirAll(path, 0755)
	// filePath => upload/users/<userID>/fileName
	filePath := path + "/" + file.Filename
	if err := c.SaveFile(file, filePath); err != nil {
		return err
	}
	// 127.0.0.1:9000/upload/users/<userID>/image.png
	userModel.Image = host + "/" + filePath

	// 4. save data to DB
	u.DB.Save(userModel)

	return nil
}
