package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/joho/godotenv"
	"github.com/mtanzim/unsplash-wallpapers/collections"
)

const baseApi = "https://api.unsplash.com"

// cap out at 10x10 = 100 images
const maxPageLimit = 10

func main() {
	fmt.Println("Hello unsplash")
	collectionPtr := flag.String("c", "", "a collection id")

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	access := os.Getenv("ACCESS")

	flag.Parse()

	if *collectionPtr == "" {
		flag.Usage()
		return
	}

	collectionIds := []string{*collectionPtr}
	downloadUrls := make(map[string]string)

	type mapMsg struct {
		key   string
		value string
	}
	writes := make(chan mapMsg)

	// isolate map mutations in a single goroutine
	go func() {
		for msg := range writes {
			downloadUrls[msg.key] = msg.value
		}
	}()

	var wgCollect sync.WaitGroup
	for _, collectionId := range collectionIds {
		for i := 1; i <= maxPageLimit; i++ {
			wgCollect.Add(1)

			go func(page int, colId string) {
				defer wgCollect.Done()
				apiUrl := fmt.Sprintf("%s/collections/%s/photos/?client_id=%s&page=%d", baseApi, colId, access, page)
				fmt.Println(apiUrl)
				resp, err := http.Get(apiUrl)
				if err != nil {
					log.Println(err)
					return
				}
				defer resp.Body.Close()

				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					log.Fatalln(err)
				}
				if resp.StatusCode != http.StatusOK {
					log.Printf("unexpected response, stausCode: %d; body: %s", resp.StatusCode, string(body))
					return
				}
				var dat collections.Collections
				if err := json.Unmarshal(body, &dat); err != nil {
					log.Printf("err: %s, body: %s", err, string(body))
					return
				}
				for _, v := range dat {
					writes <- mapMsg{key: v.ID, value: v.Urls.Full}
				}
			}(i, collectionId)

		}
	}
	wgCollect.Wait()

	var wg sync.WaitGroup
	for k, v := range downloadUrls {
		wg.Add(1)
		go func(url, fileName string) {
			defer wg.Done()
			ext := "jpg"
			fn := fmt.Sprintf("%s.%s", fileName, ext)
			err := downloadFile(url, fn)
			if err != nil {
				log.Println(err)
			}

		}(v, k)
	}
	wg.Wait()
	fmt.Println("Done")
}

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
