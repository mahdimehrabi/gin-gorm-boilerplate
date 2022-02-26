package utils

import (
	"io/ioutil"
	"net/http"
)

//URLToBinary -> URL to Binary
func URLToBinary(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	imageBinary, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return imageBinary, nil

}





