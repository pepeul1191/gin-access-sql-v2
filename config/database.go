package config

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB inicializa la conexión a la base de datos
func InitDB() (*gorm.DB, error) {
	var err error

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)

	dbDriver := GetEnv("DB_DRIVER", "sqlite")
	dbName := GetEnv("DB_NAME", "database.db")

	switch dbDriver {
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(dbName), &gorm.Config{
			Logger: newLogger,
		})
	// case "postgres":
	// 	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
	// 		GetEnv("DB_HOST", "localhost"),
	// 		GetEnv("DB_PORT", "5432"),
	// 		GetEnv("DB_USER", "postgres"),
	// 		GetEnv("DB_PASSWORD", ""),
	// 		dbName)
	// 	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
	// 		Logger: newLogger,
	// 	})
	default:
		log.Fatal("Database driver not supported")
	}

	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Configuración del pool de conexiones
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}

// GetDB devuelve la instancia de la base de datos
func GetDB() *gorm.DB {
	return db
}
