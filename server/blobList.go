package main

import (
	"cmp"
	"math"
	"slices"
)

type BlobList struct {
	Blobs []*Blob
}

func (bl *BlobList) AppendBlob(blob *Blob) {
	if bl.Blobs == nil {
		bl.Blobs = []*Blob{blob}
		return
	}

	bl.Blobs = append(bl.Blobs, blob)
}

// all documents containing that word
func (bl *BlobList) ContainingWordInBlobs(word string) float64 {
	var count float64 = 0
	for _, b := range bl.Blobs {
		if b.GetWordCount(word) > 0 {
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

func (bl *BlobList) tf_idf(query *Blob) []*Blob {
	for _, blob := range bl.Blobs {
		// if the query only has 1 word
		// if query.TotalWords == 1 {
		// 	word := query.blobFile[0]
		// 	tf := blob.TermFrecuency(word)
		// 	idf := bl.InverseDocumentFrequency(word)
		// 	blob.Score = tf * idf
		// 	continue
		// }

		blob.Score = blob.CalculateDotProduct(query)
	}

	orderedBlobs := bl.Blobs
	slices.SortFunc(orderedBlobs, func(a, b *Blob) int {
		return cmp.Compare(b.Score, a.Score)
	})

	return orderedBlobs
}
