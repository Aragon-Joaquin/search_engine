package main

import (
	"fmt"
	"strings"
)

var (
	userQuery = "Linux operative system"
	documents = []string{
		"Linux is an operative system low on resources",
		"Linux From Scratch Drops SysVinit Support. Linux . Linux . Linux .",
		"IronClaw: a Rust-based clawd that runs tools in isolated WASM sandboxes",
		"Major European payment processor can't send email to Google Workspace users",
	}
)

// WARN: first initial test - have lower expectations. might be considerable improved later on
func main() {
	blobList := BlobList{
		Blobs: []*Blob{},
	}

	// lets suppose we've already this in db, locally or whatever
	for _, strblob := range documents {
		splitted := strings.Fields(strblob)

		blob := &Blob{
			blobFile:   splitted,
			TotalWords: len(splitted),
		}
		blobList.AppendBlob(blob)
	}

	for i, blob := range blobList.Blobs {
		score := blobList.tf_idf("Linux", blob)
		fmt.Printf("%d - Score: %f\nDocument: %s\n\n", i+1, score, blob.blobFile)
	}
}
