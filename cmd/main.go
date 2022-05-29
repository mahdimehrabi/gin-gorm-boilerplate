package main

import (
	"boilerplate/core/infrastructure"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//send console argument for executing command like "go run ./cmd/ create_admin"
func main() {
	arg := os.Args[len(os.Args)-1]
	env := LoadEnv()
	switch arg {
	case "create_admin":
		CreateAdmin(env)
	default:
		log.Fatal("Unkown command!")
	}
}

func LoadEnv() infrastructure.Env {
	return infrastructure.NewEnv()
}

func GetDB(env infrastructure.Env) *gorm.DB {
	url := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s TimeZone=Europe/London",
		env.DBHost, env.DBUsername, env.DBPassword, env.DBName,
		env.DBPort)

	if env.Environment != "development" {
		url = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/London",
			env.DBHost, env.DBUsername, env.DBPassword, env.DBName,
			env.DBPort)
	}

	db, err := gorm.Open(postgres.Open(url))
	if err != nil {
		log.Fatal(err)
	}
	return db
}
