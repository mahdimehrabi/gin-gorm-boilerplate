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

	if ENVIRONMENT != "local" {
		dsn = fmt.Sprintf("%s:%s@%s:%s/%s?sslmode=disable", USER, PASS, HOST, PORT, DBNAME)
	}

	migrations, err := migrate.New("file://migration/", "postgres://"+dsn)

	m.logger.Zap.Info("--- Running Migration ---")
	err = migrations.Steps(1000)
	if err != nil {
		m.logger.Zap.Error("Error in migration: ", err.Error())
	}
}
