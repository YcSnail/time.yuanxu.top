package database

import (
	"log"

	"github.com/YcSnail/time.yuanxu.top/backend/config"
	"github.com/YcSnail/time.yuanxu.top/backend/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Connect opens the MySQL connection and runs auto-migration.
func Connect(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(cfg.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&models.User{}, &models.Countdown{}); err != nil {
		return nil, err
	}

	log.Println("database connected and migrated")
	return db, nil
}
