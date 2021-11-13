package main

import (
	"bytes"
	"path"
	"strings"
)

func Worker(url, folder string) error {
	var body []byte
	err := download(url, &body)
	if err != nil {
		return err
	}

	_, filenameWithQueryParams := path.Split(url)
	filename := strings.Split(filenameWithQueryParams, "?")[0]

	path := path.Join(folder, filename)
	reader := bytes.NewReader(body)
	err = client.UploadFile(path, true, reader)
	if err != nil {
		return err
	}

	return nil
}

func RetryWorker(retryNumber int, url, folder string) (err error) {
	for i := retryNumber; i > 0; i-- {
		err = Worker(url, folder)
		if err != nil {
			continue
		}
		return nil
	}

	return err
}
