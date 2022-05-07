package download

import (
	"encoding/json"
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
	destPath     string
}

func NewDownloader(baseApi, accessKey, destPath string, maxPageLimit int) *downloader {
	return &downloader{
		baseApi:      baseApi,
		maxPageLimit: maxPageLimit,
		accessKey:    accessKey,
		destPath:     destPath,
	}
}

func (d *downloader) Download(collectionID string) {
	collectionIds := []string{collectionID}
	downloadUrls := make(map[string]string)

	access := d.accessKey
	maxPageLimit := d.maxPageLimit
	baseApi := d.baseApi

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
	wg.Wait()
	close(writes)

	var wgD sync.WaitGroup
	for k, v := range downloadUrls {
		wgD.Add(1)
		go func(url, fileName string) {
			defer wgD.Done()
			ext := "jpg"
			fn := fmt.Sprintf("%s.%s", fileName, ext)
			err := d.downloadFile(url, fn)
			if err != nil {
				log.Println(err)
			}

		}(v, k)
	}
	wgD.Wait()
	fmt.Println("Done")
}
