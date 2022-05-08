package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/mtanzim/unsplash-wallpapers/pkg/downloader"
)

const baseApi = "https://api.unsplash.com"

func main() {
	collectionPtr := flag.String("c", "", "a collection id")
	pages := flag.Int("p", 1, "number of pages")
	dest := flag.String("d", "", "destination of download, please ensure the directory exists")

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	access := os.Getenv("ACCESS")

	if access == "" {
		log.Fatal("please setup .env with access key")
		return
	}

	flag.Parse()

	if *collectionPtr == "" {
		flag.Usage()
		return
	}

	if *dest == "" {
		flag.Usage()
		return
	}

	downloader := downloader.NewDownloader(baseApi, access, *pages)
	res := downloader.Download(*collectionPtr, *dest)
	fmt.Println(res)

}
