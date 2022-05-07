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
const maxPageLimit = 1
const destPath = "./images"

func main() {
	fmt.Println("Hello unsplash")
	collectionPtr := flag.String("c", "", "a collection id")

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

	downloader := downloader.NewDownloader(baseApi, access, destPath, maxPageLimit)
	downloader.Download(*collectionPtr)

}
