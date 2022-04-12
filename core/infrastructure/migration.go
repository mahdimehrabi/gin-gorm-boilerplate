package infrastructure

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

//Migrations -> Migration Struct
type Migrations struct {
	logger Logger
	env    Env
}

//NewMigrations -> return new Migrations struct
func NewMigrations(
	logger Logger,
	env Env,
) Migrations {
	return Migrations{
		logger: logger,
		env:    env,
	}
}

//Migrate -> migrates all table
func (m Migrations) Migrate() {
	m.logger.Zap.Info("Migrating schemas...")

	USER := m.env.DBUsername
	PASS := m.env.DBPassword
	HOST := m.env.DBHost
	PORT := m.env.DBPort
	DBNAME := m.env.DBName
	ENVIRONMENT := m.env.Environment

	dsn := fmt.Sprintf("%s:%s@%s:%s/%s", USER, PASS, HOST, PORT, DBNAME)

	if ENVIRONMENT != "development" {
		dsn = fmt.Sprintf("%s:%s@%s:%s/%s?sslmode=disable", USER, PASS, HOST, PORT, DBNAME)
	}

	migrations, err := migrate.New("file://"+m.env.BasePath+"/core/migrations", "postgres://"+dsn)
	m.logger.Zap.Info(m.env.BasePath + "/core/migrations")
	if err != nil {
		m.logger.Zap.Error("Error loading migration file: ", err.Error())
	}

	m.logger.Zap.Info("--- Running Migration ---")
	err = migrations.Up()
	if err != nil {
		m.logger.Zap.Error("Error in migration: ", err.Error())
	}
	m.logger.Zap.Info("--- Migration Completed ---")
}
