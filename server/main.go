package main

import (
	"fmt"
	"log"
	"runtime"
	"search_engine/blobs"
	"time"
)

var (
	userQuery     = "Linux operative system"
	systemThreads = runtime.NumCPU()
)

func main() {
	start := time.Now()

	blobList := blobs.LoadBlobsFromFolder()
	log.Println("LOADING AND PARSE BLOBS", "time:", time.Since(start))

	query := blobs.CreateBlob()
	query.StemWords("saturn")

	ranking := blobList.Calculate_tf_idf(query)

	fmt.Println("\nRanking in order:")
	for i, b := range ranking {
		fmt.Printf("<%d>\n - Title: %s\n - Description: %s\n - URL: %s\n - DateTime: %v\n [%f out of 1.0]\n\n", i, b.Headers.Title, b.Headers.Description, b.Headers.URL, b.Headers.Datetime, b.Score)
	}

	log.Println("TOTAL TIME", "time:", time.Since(start))
}
