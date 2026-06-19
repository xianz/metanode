package services

import (
	"blog/models"
	"blog/utils"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	Db *gorm.DB
}

// func (us *UserService) GetUserByUsername(username string) (*models.User, error) {
// 	var user models.User
// 	if err := us.Db.Where("username = ?", username).First(&user).Error; err != nil {
// 		if err == gorm.ErrRecordNotFound {
// 			return nil, &utils.AppError{Code: 404, Message: "用户不存在"}
// 		}
// 		return nil, err
// 	}
// 	return &user, nil
// }

func (us *UserService) Authenticate(username, password string) (*models.User, error) {
	var user models.User
	if err := us.Db.Where("username = ? ", username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.NewAppError(404, "用户不存在")
		}
		return nil, err
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, utils.NewAppError(401, "密码错误")
	}

	return &user, nil
}

func (us *UserService) CreateUser(req models.UserRequest) (*models.User, error) {
	var user models.User
	user.Username = req.Username
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(hashedPassword)
	err = us.Db.Create(&user).Error
	if err != nil {
		if err.Error() == "UNIQUE constraint failed: users.username" {
			return nil, utils.NewAppError(409, "用户名已存在")
		}
		return nil, err
	}
	return &user, nil
}
