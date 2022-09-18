package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Todo struct {
	ID     string `gorm:"primarykey"`
	Text   string
	Done   bool
	UserID string
}

type User struct {
	ID   string `gorm:"primarykey"`
	Name string
}

type DB struct {
	*gorm.DB
}

func Init() (*DB, error) {
	db, err := gorm.Open(sqlite.Open("main.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&Todo{}, &User{})
	if err != nil {
		return nil, err
	}

	return &DB{db}, err
}
