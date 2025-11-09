package db

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// 1. db-слой (инструмент)

// этот слой отвечает только за подключение к базе данных и миграцию
// здесь нет бизнес-логики, нет работы с HTTP
// он просто открывает соединение и отдаёт объект GORM наружу

var DB *gorm.DB

// функция для инициализации подключения и работы с бд
func InitDB() {
	// источник данных (инфа о нашей бд)
	dsn := "host=localhost user=postgres password=yourpassword dbname=postgres port=5432 sslmode=disable"
	var err error

	// открываем соединение с бд (по нашим данным)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Could not connect to Database: %v", err)
	}
}
