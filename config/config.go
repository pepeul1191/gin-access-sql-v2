package config

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Config struct {
	DB *gorm.DB
}

func NewConfig() (*Config, error) {
	db, err := gorm.Open(sqlite.Open("db/app.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &Config{DB: db}, nil
}
