package download

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func downloadFile(URL, fileName string) error {
	response, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return errors.New("received non 200 response code")
	}

	dir := "images"
	path := fmt.Sprintf("%s/%s", dir, fileName)
	file, err := create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, response.Body)
	fmt.Printf("Downloaded %s\n", path)
	if err != nil {
		return err
	}

	return nil
}

// TODO: will this work correctly if run from elsewhere?
func create(p string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(p), 0770); err != nil {
		return nil, err
	}
	return os.Create(p)
}
