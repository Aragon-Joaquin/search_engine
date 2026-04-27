package blobs

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"slices"
	"strings"
	"time"
	"unicode"

	"search_engine/internal/utils"
	"search_engine/third_party/stemmer"
)

type blobHeaders struct {
	Folder      utils.INDEXERS `redis:"folder"`
	Title       string         `redis:"title"`
	Description string         `redis:"description"`
	Datetime    time.Time      `redis:"datetime"`
	URL         string         `redis:"url"`
}

type Blob struct {
	// count all the ZScore of redis and stores it in this field
	Length int `redis:"length"`

	// Location string
	Score float64
	// this is calculated at the end, represents how exact was the match made with the query from 0 to 1

	TermSpace map[string]int // map of the frecuency of each word
	magnitude float64        // magnitude based on word frecuency

	blobHeaders `redis:"headers"`
}

func CreateBlob() *Blob {
	return &Blob{
		TermSpace: map[string]int{},
		Score:     -1,

		magnitude:   -1,
		blobHeaders: blobHeaders{},
	}
}

// NOTE: no idea if im going to use this
func (b *Blob) MarshalToBinary() ([]byte, error) {
	return json.Marshal(b)
}

func (b *Blob) StemWords(content string) {
	parsed := strings.FieldsFunc(content, func(r rune) bool {
		return unicode.IsPunct(r) || unicode.IsSpace(r) || unicode.IsSymbol(r)
	})

	skipStopWords := []string{}
	for _, w := range parsed {
		if ok := slices.Contains(stopWords, w); !ok {
			continue
		}

		skipStopWords = append(skipStopWords, w)
	}

	stemmer := stemmer.StemMultiple(skipStopWords)

	b.Length = len(stemmer)
	b.TermSpace = b.SetTermSpace(stemmer)
}

func (b *Blob) GetWordCount(word string) int {
	b.initTermSpace() // uhhh
	if val, ok := b.TermSpace[word]; ok {
		return val
	}
	return 0
}

func (b *Blob) initTermSpace() {
	if b.TermSpace == nil {
		b.TermSpace = map[string]int{}
	}
}

func (b *Blob) SetTermSpace(stemmedContent []string) map[string]int {
	b.initTermSpace()

	for _, w := range stemmedContent {
		val, _ := b.TermSpace[w]
		b.TermSpace[w] = val + 1
	}

	return b.TermSpace
}

func (b *Blob) AddToTermSpace(word string, score int) {
	b.initTermSpace()
	b.TermSpace[word] = score
}

// NOTE: math below
// calculate the magnitude of the vector using the pythagorean theorem
// sqrt(a² + b² + c² ... + n²)
func (b *Blob) GetVectorMagnitute() float64 {
	if b.magnitude >= 0 {
		return b.magnitude
	}

	var magnitude int = 0
	for _, val := range b.TermSpace {
		magnitude += val * val
	}

	finalMag := math.Sqrt(float64(magnitude))
	b.magnitude = finalMag

	return finalMag
}

// stands for tf - total times a word appears in blob
func (b *Blob) TermFrecuency(word string) float64 {
	return float64(b.GetWordCount(word)) / float64(b.Length)
}

//	Q 	* 	V
//	-------------
// ||Q|| x ||V1||

func (b *Blob) CalculateDotProduct(query *Blob) float64 {
	var dotProduct int = 0
	for key, value := range b.TermSpace {
		val := query.GetWordCount(key)
		dotProduct += val * value
	}

	return float64(dotProduct) / (query.GetVectorMagnitute() * b.GetVectorMagnitute())
}

// BUG: handle files with "," (comma) in their names or brief descriptions
func (b *Blob) ParseBlobContentsIntoFile(file *os.File, content *[]byte) error {
	defer file.Close()

	// header
	_, err := file.WriteString(
		fmt.Sprint(
			b.Title,
			",",
			b.Description,
			",",
			b.Datetime,
			",",
			b.URL,
			"\n",
		),
	)
	if err != nil {
		return err
	}

	_, err = file.WriteString(string(*content))
	return err
}

// both saves into the redisdb + local
func (b *Blob) SaveBlob(folder utils.INDEXERS, term string, content *[]byte) error {
	f, err := utils.CreateFile(folder, term)
	if err != nil {
		return err
	}

	return b.ParseBlobContentsIntoFile(f, content)
}
