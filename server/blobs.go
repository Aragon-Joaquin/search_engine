package main

import (
	"math"
)

type Blob struct {
	// document
	blobFile []string

	TotalWords int
}

func (b *Blob) CountWordInBlob(word string) float64 {
	var count int = 0
	for _, w := range b.blobFile {
		if w == word {
			count++
		}
	}

	return float64(count)
}

// stands for tf - total times a word appears in blob
func (b *Blob) TermFrecuency(word string) float64 {
	return b.CountWordInBlob(word) / float64(b.TotalWords)
}

// NOTE: BlobList below
type BlobList struct {
	Blobs []*Blob
}

// all documents containing that word
func (bl *BlobList) ContainingWordInBlobs(word string) float64 {
	var count float64 = 0
	for _, b := range bl.Blobs {
		if b.CountWordInBlob(word) > 0 {
			count++
		}
	}
	return count
}

// stands for idf - measures how common a word is
func (bl *BlobList) InverseDocumentFrequency(word string) float64 {
	count := float64(len(bl.Blobs)) / (1 + bl.ContainingWordInBlobs(word))

	return math.Log(count)
}

// computes the score
func (bl *BlobList) tf_idf(word string, blob *Blob) float64 {
	tf := blob.TermFrecuency(word)
	idf := bl.InverseDocumentFrequency(word)

	return tf * idf
}

// private methods lol
func (bl *BlobList) AppendBlob(blob *Blob) {
	if bl.Blobs == nil {
		bl.Blobs = []*Blob{blob}
		return
	}

	bl.Blobs = append(bl.Blobs, blob)
}
