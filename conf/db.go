package conf

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func GetDBConnectionURL() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", configuration.Database.DbUser, configuration.Database.DbPassword,
		configuration.Database.DbHost, configuration.Database.DbPort, configuration.Database.DbName)
}

func DBConn() *gorm.DB {
	WriteDbUrl := GetDBConnectionURL()
	fmt.Println(WriteDbUrl)
	var db *gorm.DB
	var err error
	db, err = gorm.Open(mysql.Open(WriteDbUrl), GetGormConfig())
	if err != nil {
		panic(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)
	return db
}

func DBConnWithLoglevel(logMode logger.LogLevel) *gorm.DB {
	dbUrl := GetDBConnectionURL()
	var db *gorm.DB
	var err error
	db, err = gorm.Open(mysql.Open(dbUrl), &gorm.Config{
		Logger: logger.Default.LogMode(logMode),
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Database CONNECTED")

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db
}

func GetGormConfig() *gorm.Config {
	logMode := logger.Info
	return &gorm.Config{
		Logger: logger.Default.LogMode(logMode),
	}
}
