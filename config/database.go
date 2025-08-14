package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() error {
	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error cargando el archivo .env")
	}

	var err error
	dbDriver := os.Getenv("DB_DRIVER")
	dbName := os.Getenv("DB_NAME")

	switch dbDriver {
	case "sqlite":
		DB, err = gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	// case "postgres":
	// 	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
	// 		os.Getenv("DB_HOST"),
	// 		os.Getenv("DB_PORT"),
	// 		os.Getenv("DB_USER"),
	// 		os.Getenv("DB_PASSWORD"),
	// 		dbName,
	// 		os.Getenv("DB_SSL_MODE"))
	// 	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	default:
		log.Fatal("Driver de BD no soportado")
	}

	if err != nil {
		return err
	}
	return nil
}

func GetDB() *gorm.DB {
	return DB
}
