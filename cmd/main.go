package main

import (
	"boilerplate/core/infrastructure"
	"boilerplate/core/models"
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

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

func CreateAdmin(env infrastructure.Env) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Create admin")
	fmt.Println("---------------------")

	fmt.Print("Please enter FirstName: ")
	firstName, _ := reader.ReadString('\n')
	firstName = strings.Replace(firstName, "\n", "", -1)

	fmt.Print("Please enter LastName: ")
	lastName, _ := reader.ReadString('\n')
	lastName = strings.Replace(lastName, "\n", "", -1)

	fmt.Print("Please enter email: ")
	email, _ := reader.ReadString('\n')
	email = strings.Replace(email, "\n", "", -1)

	db := GetDB(env)
	var count int64
	err := db.Model(&models.User{}).Where("email = ?", email).Count(&count).Error
	if err != nil {
		log.Fatal(err)
	}
	if count > 0 {
		log.Fatal("Can't create admin , this email already exist")
	}

	fmt.Print("Please enter you password: ")
	password, _ := reader.ReadString('\n')
	password = strings.Replace(password, "\n", "", -1)
	encryption := infrastructure.NewEncryption(infrastructure.Logger{}, env)
	password = encryption.SaltAndSha256Encrypt(password)
	user := models.User{
		FirstName:     firstName,
		LastName:      lastName,
		Email:         email,
		Password:      password,
		IsAdmin:       true,
		VerifiedEmail: true,
	}
	err = db.Create(&user).Error
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Admin Created Successfuly!")
}
