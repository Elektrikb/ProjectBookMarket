package services

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Db *gorm.DB

// InitDB инициализирует подключение к базе данных и выполняет миграцию схемы.
// @Summary Инициализация базы данных
// @Description Устанавливает соединение с базой данных и выполняет миграцию для модели Book.
// @Tags database
// @Success 200 {string} string "Database initialized successfully"
// @Failure 500 {string} string "Failed to connect to database"
// @Router /initdb [post]
func InitDB() {
	dsn := "host=213.171.10.112 user=postgres password=67 dbname=bookdb port=5432 sslmode=disable"
	var err error
	Db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	Db.AutoMigrate(&Book{})
}

type Book struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	Title     string `json:"title"`
	Author    string `json:"author"`
	Year      int    `json:"year"`
	Publisher int    `json:"publisher"`
}
