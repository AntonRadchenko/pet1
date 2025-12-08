package userService

import (
	"errors"
	"strings"
	"time"

	"github.com/AntonRadchenko/WebPet1/internal/db"
	"github.com/AntonRadchenko/WebPet1/internal/taskService"
	"gorm.io/gorm"
)

type UserRepoInterface interface {
	Create(user *UserStruct) (*UserStruct, error)
	GetAll() ([]UserStruct, error)
	GetByID(id uint) (UserStruct, error)
	GetTasksForUser(userID uint) ([]taskService.TaskStruct, error)
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
			return UserStruct{}, err
		}
		return user, err
	}
	return user, nil
}

func (r *UserRepo) GetTasksForUser(userID uint) ([]taskService.TaskStruct, error) {
	var user UserStruct

	// Загружает пользователя вместе со всеми его задачами одним запросом
	err := db.DB.Preload("Tasks").First(&user, userID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("user not found")
		}
		return nil, err		
	}
	return user.Tasks, nil // возвращаем все таски пользователя (Tasks - поле в модели бд (слайс тасок))
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