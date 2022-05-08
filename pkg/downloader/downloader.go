package downloader

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/mtanzim/unsplash-wallpapers/collections"
)

type downloader struct {
	baseApi      string
	maxPageLimit int
	accessKey    string
}

func NewDownloader(baseApi, accessKey string, maxPageLimit int) *downloader {

	return &downloader{
		baseApi:      baseApi,
		maxPageLimit: maxPageLimit,
		accessKey:    accessKey,
	}
}

func (d *downloader) Download(collectionID string, destPath string) []string {

	dirExists, err := exists(destPath)
	if !dirExists || err != nil {
		return []string{"Destination directroy does not exist"}
	}

	urls, urlErrorMessages := d.collectUrls(collectionID)
	downloadMessages := d.triggerDownloads(urls, destPath)
	return append(urlErrorMessages, downloadMessages...)

}

func (d *downloader) collectUrls(collectionID string) (map[string]string, []string) {

	collectionIds := []string{collectionID}
	downloadUrls := make(map[string]string)
	urlErrors := []string{}

	access := d.accessKey
	maxPageLimit := d.maxPageLimit
	baseApi := d.baseApi

	type mapMsg struct {
		key   string
		value string
	}
	writes := make(chan mapMsg)
	errorWrites := make(chan string)

	// isolate map mutations in a single goroutine
	go func() {
		for msg := range writes {
			downloadUrls[msg.key] = msg.value
		}
	}()

	// isolate slice mutations in a single goroutine
	go func() {
		for msg := range errorWrites {
			urlErrors = append(urlErrors, msg)
		}
	}()

	var wg sync.WaitGroup
	for _, collectionId := range collectionIds {
		for i := 1; i <= maxPageLimit; i++ {
			wg.Add(1)

			go func(page int, colId string) {
				defer wg.Done()
				apiUrl := fmt.Sprintf("%s/collections/%s/photos/?client_id=%s&page=%d", baseApi, colId, access, page)
				fmt.Println(apiUrl)
				resp, err := http.Get(apiUrl)
				if err != nil {
					log.Println(err)
					errorWrites <- err.Error()
					return
				}
				defer resp.Body.Close()

				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					errorWrites <- err.Error()
					log.Println(err)
					return

				}
				if resp.StatusCode != http.StatusOK {
					errMsg := fmt.Sprintf("unexpected response, stausCode: %d; body: %s", resp.StatusCode, string(body))
					errorWrites <- errors.New(errMsg).Error()
					log.Print(errMsg)
					return
				}
				var dat collections.Collections
				if err := json.Unmarshal(body, &dat); err != nil {
					errorWrites <- err.Error()
					log.Printf("err: %s, body: %s", err, string(body))
					return
				}
				for _, v := range dat {
					writes <- mapMsg{key: v.ID, value: v.Urls.Full}
				}
			}(i, collectionId)

		}
	}
	wg.Wait()
	close(writes)
	close(errorWrites)
	return downloadUrls, urlErrors
}

func (d *downloader) triggerDownloads(downloadUrls map[string]string, destPath string) []string {

	results := []string{}
	writes := make(chan string)

	// isolate slice mutations in a single goroutine
	go func() {
		for msg := range writes {
			results = append(results, msg)
		}
	}()

	var wgD sync.WaitGroup
	for k, v := range downloadUrls {
		wgD.Add(1)
		go func(url, fileName string) {
			defer wgD.Done()
			ext := "jpg"
			fn := fmt.Sprintf("%s.%s", fileName, ext)
			path, err := downloadFile(url, fn, destPath)
			if err != nil {
				log.Println(err)
				writes <- err.Error()
				return
			}
			writes <- fmt.Sprintf("Downloaded %s to %s", url, path)

		}(v, k)
	}
	wgD.Wait()
	close(writes)
	return results
}
