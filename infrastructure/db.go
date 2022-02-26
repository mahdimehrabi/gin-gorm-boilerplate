package infrastructure

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Database modal
type Database struct {
	DB *gorm.DB
}

// NewDatabase creates a new database instance
func NewDatabase(Zaplogger Logger, env Env) Database {

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // Disable color
		},
	)

	Zaplogger.Zap.Info(env)

	url := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s TimeZone=Europe/London",
		env.DBHost, env.DBUsername, env.DBPassword, env.DBName,
		env.DBPort)

	if env.Environment != "local" {
		url = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/London",
			env.DBHost, env.DBUsername, env.DBPassword, env.DBName,
			env.DBPort)
	}

	db, err := gorm.Open(postgres.Open(url), &gorm.Config{Logger: newLogger})
	if err != nil {
		Zaplogger.Zap.Info("Url: ", url)
		Zaplogger.Zap.Panic(err)
	}

	Zaplogger.Zap.Info("Database connection established ✔️")

	return Database{
		DB: db,
	}
}
