package taskService

import "time"

// модель базы данных
type TaskStruct struct {
	ID        uint `gorm:"primaryKey;autoIncrement"`
	UserId    uint `gorm:"not null;index"`
	Task      string
	IsDone    bool
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func (TaskStruct) TableName() string {
    return "task_structs"  // как в миграции
}

// единый формат ошибок
// type ErrorStruct struct {
// 	Error string `json:"error"`
// }
