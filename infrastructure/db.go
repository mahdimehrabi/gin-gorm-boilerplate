package infrastructure

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Database modal
type Database struct {
	DB *gorm.DB
}

// NewDatabase creates a new database instance
func NewDatabase(zapLogger Logger, env Env) Database {
	logger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // Disable color
		},
	)
	zapLogger.Zap.Info(env)
	if env.Environment == "test" {
		RemoveDB(zapLogger, env, "")
		CreateDB(zapLogger, env, "")
	}
	return GetDB(logger, zapLogger, env)
}

// connect to database and return Database object
func GetDB(logger logger.Interface, zapLogger Logger, env Env) Database {
	url := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s TimeZone=Europe/London",
		env.DBHost, env.DBUsername, env.DBPassword, env.DBName,
		env.DBPort)

	if env.Environment != "development" {
		url = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/London",
			env.DBHost, env.DBUsername, env.DBPassword, env.DBName,
			env.DBPort)
	}

	db, err := gorm.Open(postgres.Open(url), &gorm.Config{Logger: logger})
	if err != nil {
		zapLogger.Zap.Info("Url: ", url)
		zapLogger.Zap.Panic(err)
	}

	zapLogger.Zap.Info("Database connection established ✔️")

	return Database{
		DB: db,
	}
}

//create database if no DBName passed , it automaticaly use dbname of environment varible
func CreateDB(zapLogger Logger, env Env, DBName string) {
	db, err := ConnectDBSQL(env)
	if err != nil {
		zapLogger.Zap.Panic(err)
	}
	defer db.Close()

	if DBName == "" {
		DBName = env.DBName
	}

	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s;", DBName))
	if err != nil {
		zapLogger.Zap.Panic(err)
	}
}

//remove passed database if no DBName passed , it automaticaly use dbname of environment varible
func RemoveDB(zapLogger Logger, env Env, DBName string) {
	db, err := ConnectDBSQL(env)
	if err != nil {
		zapLogger.Zap.Panic(err)
	}
	defer db.Close()

	if DBName == "" {
		DBName = env.DBName
	}

	_, err = db.Exec(fmt.Sprintf("DROP DATABASE %s;", DBName))
	if err != nil {
		zapLogger.Zap.Panic(err)
	}
}

//connect to database without passing database name and with database/sql package
//useful for doing general database sql statements that not related to a specefic database
//be sure to defer db.Close() in using function
func ConnectDBSQL(env Env) (*sql.DB, error) {
	url := fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=disable",
		env.DBHost, env.DBPort, env.DBUsername, env.DBPassword)
	return sql.Open("postgres", url)
}
