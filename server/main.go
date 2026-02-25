package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"runtime"
	"search_engine/blobs"
	"sync"
	"time"
)

var (
	userQuery     = "Linux operative system"
	systemThreads = runtime.NumCPU()
	maxTimeout    = time.Second * 3

	start time.Time
)

func init() {
	start = time.Now()
}

func main() {
	// loads all blobs again into the redisDB
	flagLoadBlobs := flag.Bool("l", false, "requires a bool")
	flag.Parse()

	if flagLoadBlobs != nil && *flagLoadBlobs {
		log.Println("USED FLAG -l - LOADING ALL ./data/* BLOBS TO REDIS")
		loadBlobs()
		log.Println("UPLOAD FINISHED IN: ", time.Since(start))
		return
	}

	// user query
	query := blobs.CreateBlob()
	query.StemWords("saturn")

	ctx := context.Background()
	bList, err := DBRedis.GetAllZBlobs(ctx)
	if err != nil {
		panic(err)
	}

	ranking := bList.Calculate_tf_idf(query)

	fmt.Println("\nRanking in order:")
	for i, b := range ranking {
		fmt.Printf("<%d>\n - Title: %s\n - Description: %s\n - URL: %s\n - DateTime: %v\n [%f out of 1.0]\n\n", i, b.Title, b.Description, b.URL, b.Datetime, b.Score)
	}

	log.Println("TOTAL TIME", "time:", time.Since(start))
}

func loadBlobs() {
	ctx, cancel := context.WithTimeout(context.Background(), maxTimeout)
	defer cancel()
	blobList := blobs.LoadBlobsFromFolder()

	var wg sync.WaitGroup
	waitChan := make(chan any)

	go func() {
		for _, blob := range blobList.Blobs {
			wg.Go(func() {
				if err := DBRedis.AddZSort(ctx, blob); err != nil {
					log.Println("Error in one of the blobs while trying to load it to redis: ", err)
				}
			})
		}
		wg.Wait()
		close(waitChan)
	}()

	select {
	case <-waitChan:
		return
	case <-ctx.Done():
		panic("timeout'ed while loading all blobs from local to redis")
	}
}
