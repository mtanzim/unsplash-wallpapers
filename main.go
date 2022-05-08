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

// cap number of images out at nx10
const destPath = "./images"

func main() {
	collectionPtr := flag.String("c", "", "a collection id")
	pages := flag.Int("p", 1, "number of pages")

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

	downloader := downloader.NewDownloader(baseApi, access, destPath, *pages)
	res := downloader.Download(*collectionPtr)
	fmt.Println(res)

}
