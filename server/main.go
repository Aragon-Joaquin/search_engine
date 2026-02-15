package main

import (
	"fmt"
)

var (
	userQuery = "Linux operative system"
	documents = []string{
		"cat cat cat dog mouse mouse mouse mouse",
		"cat dog dog mouse mouse mouse mouse mouse",
		"cat cat dog dog dog",
		// "eating eats eaten",
		// "Linux is an operative system low on resources",
		// "Linux From Scratch Drops SysVinit Support. Linux . Linux . Linux .",
		// "IronClaw: a Rust-based clawd that runs tools in isolated WASM sandboxes",
		// "Major European payment processor can't send email to Google Workspace users",
	}
)

// WARN: first initial test - have lower expectations. might be considerable improved later on
func main() {
	blobList := BlobList{
		Blobs: []*Blob{},
	}

	// lets suppose we've already this in db, redis, locally or whatever
	for _, strblob := range documents {
		blob := CreateBlob(strblob)
		blobList.AppendBlob(blob)
	}

	query := CreateBlob("mouse")
	ranking := blobList.tf_idf(query)

	fmt.Println("\nRanking in order:")
	for i, b := range ranking {
		fmt.Printf("<%d>\n - %s \n - %f out of 1.0\n", i, b.blobFile, b.Score)
	}
}
