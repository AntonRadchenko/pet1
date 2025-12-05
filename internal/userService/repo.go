package userService

import (
	"errors"
	"strings"
	"time"

	"github.com/AntonRadchenko/WebPet1/internal/db"
	"gorm.io/gorm"
)

type UserRepoInterface interface {
	Create(user *UserStruct) (*UserStruct, error)
	GetAll() ([]UserStruct, error)
	GetByID(id uint) (UserStruct, error)
	// GetByEmail(email string) (UserStruct, error) // вместо этого метода, сделать в бд (в миграции) constraints по полю Email
	Update(user *UserStruct) (*UserStruct, error)
	Delete(user *UserStruct) error
}

type UserRepo struct{}

func (r *UserRepo) Create(user *UserStruct) (*UserStruct, error) {
	err := db.DB.Create(user).Error
	if err != nil {
		// првоеряем ошибку бд на дупликат бд
		if strings.Contains(err.Error(), "duplicate key") {
			return nil, errors.New("email already exists")
		}
		return nil, err
	}
	return user, nil
}

func (r *UserRepo) GetAll() ([]UserStruct, error) {
	var users []UserStruct

	err := db.DB.Find(&users).Error
	if err != nil {
		// если таблицы нет, то вместо ошибки возвращаем пустой массив []
		if strings.Contains(err.Error(), "relation") {
			return []UserStruct{}, nil
		}
		return nil, err
	}
	return users, nil
}

func (r *UserRepo) GetByID(id uint) (UserStruct, error) {
	var user UserStruct

	err := db.DB.First(&user, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return UserStruct{}, nil
		}
		return user, err
	}
	return user, nil
}

func (r *UserRepo) Update(user *UserStruct) (*UserStruct, error) {
	user.UpdatedAt = time.Now()
	err := db.DB.Save(&user).Error 
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepo) Delete(user *UserStruct) error {
	now := time.Now()
	user.DeletedAt = &now
	err := db.DB.Delete(user).Error
	if err != nil {
		return err
	}
	return nil
}