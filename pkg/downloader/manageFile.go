package downloader

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func downloadFile(URL, fileName, destPath string) (string, error) {
	response, err := http.Get(URL)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", errors.New("received non 200 response code")
	}

	dir := destPath
	path := fmt.Sprintf("%s/%s", dir, fileName)
	file, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	_, err = io.Copy(file, response.Body)
	fmt.Printf("Downloaded %s\n", path)
	if err != nil {
		return "", err
	}

	return path, nil
}
