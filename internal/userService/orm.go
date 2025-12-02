package userService

import "time"

// модель базы данных
type UserStruct struct {
	ID uint `gorm:"primaryKey;autoIncrement"`
	Email string 
	Password string 
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}