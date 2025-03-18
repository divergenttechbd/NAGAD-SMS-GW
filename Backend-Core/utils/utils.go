package utils

import (
	"myproject/config"
	"gorm.io/gorm"
)

var (
	// db  *gorm.DB
	cfg *config.Config
)

// Init initializes the utils package and loads the configuration
func Init() {
	cfg = config.GetConfig()
}

// GetConfig returns the application configuration
func GetConfig() *config.Config {
	return cfg
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return db
}