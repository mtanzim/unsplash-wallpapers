package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/mtanzim/unsplash-wallpapers/pkg/download"
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

	if access == "" {
		log.Fatal("please setup .env with access key")
		return
	}

	flag.Parse()

	if *collectionPtr == "" {
		flag.Usage()
		return
	}

	downloader := download.NewDownloader(baseApi, access, maxPageLimit)
	downloader.Download(*collectionPtr)

}
