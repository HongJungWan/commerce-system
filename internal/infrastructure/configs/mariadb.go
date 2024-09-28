package configs

import (
	"fmt"
	"github.com/HongJungWan/commerce-system/internal/helper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

const (
	charset        = "utf8mb4"
	parseTime      = "True"
	loc            = "Local"
	defaultLogMode = logger.Info
)

func ConnectionDB(config *Config) *gorm.DB {
	dsn := buildDSN(config)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: logger.Default.LogMode(defaultLogMode),
	})
	helper.ErrorPanic(err)

	return db
}

func buildDSN(config *Config) string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%s&loc=%s",
		config.DBUserName,
		config.DBPassword,
		config.DBHost,
		config.DBPort,
		config.DBName,
		charset,
		parseTime,
		loc,
	)
}
