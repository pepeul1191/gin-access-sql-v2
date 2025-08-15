package config

import (
	"os"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

// LoadEnv carga las variables de entorno
func LoadEnv() error {
	if err := godotenv.Load(".env"); err != nil {
		return err
	}
	return nil
}

// GetEnv obtiene variables de entorno con valor por defecto
func GetEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
