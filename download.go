package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

var ErrNotFound = errors.New("Not found")

func download(url string, body *[]byte) error {
	resp, err := http.Get(url)
	if resp.StatusCode == http.StatusNotFound {
		return ErrNotFound
	}
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	*body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return nil
}

// Example:
// zeroStrings(3, 10) -> 0007
// zeroStrings(5, 10) -> 000010
func zeroStrings(length int, number int) string {
	numberString := strconv.Itoa(number)
	for len(numberString) < length {
		numberString = "0" + numberString
	}

	return numberString
}

func generateUrl(templateUrl string, start int) chan string {
	urlChan := make(chan string)

	go func() {
		for i := start; ; i++ {
			url := fmt.Sprintf(templateUrl, zeroStrings(5, i))
			urlChan <- url
		}
	}()

	return urlChan
}
