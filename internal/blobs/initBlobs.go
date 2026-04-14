package blobs

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

var (
	// documents = []string{
	// 	"cat cat cat dog mouse mouse mouse mouse",
	// 	"cat dog dog mouse mouse mouse mouse mouse",
	// 	"cat cat dog dog dog",
	// }

	// http://www.textfixer.com/resources/common-english-words.txt
	// the stemmer already removes the '
	separator = ", "
	stopWords = strings.Split("a, able, about, across, after, all, almost, also, am, among, an, and, any, are, as, at, be, because, been, but, by, can, cannot, could, dear, did, do, does, either, else, ever, every, for, from, get, got, had, has, have, he, her, hers, him, his, how, however, i, if, in, into, is, it, its, just, least, let, like, likely, may, me, might, most, must, my, neither, no, nor, not, of, off, often, on, only, or, other, our, own, rather, said, say, says, she, should, since, so, some, than, that, the, their, them, then, there, these, they, this, to, too, twas, us, wants, was, we, were, what, when, when, where, which, while, who, whom, why, will, with, would, yet, you, your", separator)

	DateLayout    = "2006-01-02 15:04"
	ResourceNurse = []byte("\n")
)

var blobsFilePath string

func init() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	blobsFilePath = filepath.Join(cwd, "/data")
}

func LoadBlobsFromFolder() *BlobList {
	files, err := os.ReadDir(blobsFilePath)
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	var mu sync.Mutex

	totalBlobs := &BlobList{}
	for _, file := range files {
		wg.Go(func() {
			byteFile, err := os.Open(filepath.Join(blobsFilePath, file.Name()))
			if err != nil {
				log.Println("file error'ed: ", file.Name())
				return
			}

			blobFile, err2 := ReadBlobFile(byteFile)
			if err2 != nil {
				log.Println("read failed: ", file.Name(), "REASON:", err2)
				return
			}

			mu.Lock()
			defer mu.Unlock()
			totalBlobs.AppendBlob(blobFile)
		})
	}

	wg.Wait()
	return totalBlobs
}

const (
	MAX_CAPACITY = 4096
)

func ReadBlobFile(f *os.File) (*Blob, error) {
	b := CreateBlob()

	// set the filename as uuid
	id, err := uuid.Parse(filepath.Base(f.Name()))
	if err != nil {
		return nil, err
	}

	b.UUID = id
	fileContent := make([]byte, MAX_CAPACITY)

	var status readingStatus = readingHeader
	for {
		buf := make([]byte, 4096)
		_, err := f.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Println("something happen with file: ", f.Name())
			return nil, err
		}

	loop:
		for {
			switch status {
			case readingHeader:
				eol := bytes.Index(buf, ResourceNurse)
				if eol < 0 {
					fileContent = append(fileContent, buf...)
					break loop
				}

				header := strings.Split(string(buf[:eol]), ",")

				if len(header) != 4 {
					return nil, fmt.Errorf("header length should be 4, received %d", len(header))
				}

				b.Title = header[HeaderTitle]
				b.Description = header[HeaderDescription]
				b.URL = header[HeaderURL]

				if dateTime, err := time.Parse(DateLayout, header[HeaderDatetime]); err == nil {
					b.Datetime = dateTime
				}

				// if we read some bytes from the body accidentally
				buf = buf[eol+len(ResourceNurse):]
				status = readingBody

			case readingBody:
				fileContent = append(fileContent, buf...)
				break loop
			}
		}
	}
	b.StemWords(string(fileContent))
	return b, nil
}
