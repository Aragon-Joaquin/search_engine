package blobs

import (
	"math"
	"search_engine/stemmer"
	"strings"
	"time"
	"unicode"
)

type Blob struct {
	// Location string
	Score float64 // this is calculated at the end, represents how exact was the match made with the query from 0 to 1

	// document
	blobFile []string // string of words strimmed

	termSpace map[string]int // map of the frecuency of each word
	magnitude float64        // magnitude based on word frecuency

	Headers struct {
		Title       string
		Description string
		Datetime    time.Time
		URL         string
	}
}

func CreateBlob() *Blob {
	return &Blob{
		blobFile:  []string{},
		termSpace: map[string]int{},
		magnitude: -1,
		Score:     -1,
	}
}

func (b *Blob) StemWords(content string) {
	parsed := strings.FieldsFunc(content, func(r rune) bool {
		return unicode.IsPunct(r) || unicode.IsSpace(r) || unicode.IsSymbol(r)
	})

	stemmer := stemmer.StemMultiple(parsed)

	b.blobFile = stemmer
}

func (b *Blob) GetWordCount(word string) int {
	if val, ok := b.termSpace[word]; ok {
		return val
	}

	var count int = 0
	for _, w := range b.blobFile {
		// maybe
		val, _ := b.termSpace[w]
		b.termSpace[w] = val + 1

		if w == word {
			count++
		}
	}

	return count
}

// calculate the magnitude of the vector using the pythagorean theorem
// sqrt(a² + b² + c² ... + n²)
func (b *Blob) GetVectorMagnitute() float64 {
	if b.magnitude >= 0 {
		return b.magnitude
	}

	var magnitude int = 0

	termSpace := b.GetTermSpace()
	for _, val := range termSpace {
		magnitude += val * val
	}

	finalMag := math.Sqrt(float64(magnitude))
	b.magnitude = finalMag

	return finalMag
}

func (b *Blob) GetTermSpace() map[string]int {
	if len(b.termSpace) > 0 {
		return b.termSpace
	}

	for _, w := range b.blobFile {
		val, _ := b.termSpace[w]
		b.termSpace[w] = val + 1
	}

	return b.termSpace
}

// stands for tf - total times a word appears in blob
func (b *Blob) TermFrecuency(word string) float64 {
	return float64(b.GetWordCount(word)) / float64(len(b.blobFile))
}

//	Q 	* 	V
//	-------------
// ||Q|| x ||V1||

func (b *Blob) CalculateDotProduct(query *Blob) float64 {
	var dotProduct int = 0
	for key, value := range query.GetTermSpace() {
		val := b.GetWordCount(key)
		dotProduct += val * value

	}

	return float64(dotProduct) / (query.GetVectorMagnitute() * b.GetVectorMagnitute())
}
