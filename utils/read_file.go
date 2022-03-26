package utils

import (
	"io/ioutil"
	"log"
	"os"
)

func ReadFile(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	return ioutil.ReadAll(file)
}
