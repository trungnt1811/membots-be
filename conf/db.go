package conf

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/prometheus"
)

func GetWriteDBConnectionURL() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", configuration.Database.DbUser, configuration.Database.DbPassword,
		configuration.Database.DbHost, configuration.Database.DbPort, configuration.Database.DbName)
}

func DBConn() *gorm.DB {
	WriteDbUrl := GetWriteDBConnectionURL()
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
	WriteDbUrl := GetWriteDBConnectionURL()
	var db *gorm.DB
	var err error
	db, err = gorm.Open(mysql.Open(WriteDbUrl), &gorm.Config{
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

	db.Use(prometheus.New(prometheus.Config{
		DBName:          configuration.Database.DbName, // `DBName` as metrics label
		RefreshInterval: 15,                            // refresh metrics interval (default 15 seconds)
		StartServer:     false,                         // start http server to expose metrics
		HTTPServerPort:  configuration.AppPort,         // configure http server port, default port 8080 (if you have configured multiple instances, only the first `HTTPServerPort` will be used to start server)
		MetricsCollector: []prometheus.MetricsCollector{
			&prometheus.MySQL{VariableNames: []string{"Threads_running"}},
		},
	}))

	return db
}

func GetGormConfig() *gorm.Config {
	logMode := logger.Info
	return &gorm.Config{
		Logger: logger.Default.LogMode(logMode),
	}
}
