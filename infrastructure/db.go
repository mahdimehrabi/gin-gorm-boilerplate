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

	return getNormalDB(logger, zapLogger, env)
}

func getNormalDB(logger logger.Interface, zapLogger Logger, env Env) Database {
	url := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s TimeZone=Europe/London",
		env.DBHost, env.DBUsername, env.DBPassword, env.DBName,
		env.DBPort)

	if env.Environment != "local" {
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

// func createGetTestDB() Database {

// }
