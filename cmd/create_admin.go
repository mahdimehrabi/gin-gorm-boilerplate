package main

import (
	"boilerplate/core/infrastructure"
	"boilerplate/core/models"
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

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
