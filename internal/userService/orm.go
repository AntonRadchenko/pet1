package userService

import (
	"time"

	"github.com/AntonRadchenko/WebPet1/internal/taskService"
)

// модель базы данных
type UserStruct struct {
	ID uint `gorm:"primaryKey;autoIncrement"`
	Tasks []taskService.TaskStruct `gorm:"foreignkey:UserID"`
	Email string 
	Password string 
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func (UserStruct) TableName() string {
    return "user_structs"  // как в миграции
}