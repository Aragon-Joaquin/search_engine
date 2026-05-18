package blobs

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"search_engine/internal/indexers"
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
	blobList := BlobList{}

	blobFile, err2 := blobList.ReadBlobsFromLocalFolder(indexers.GetWikipediaIndexer())
	if err2 != nil {
		return &blobList
	}

	for _, b := range blobFile {
		blobList.AppendBlob(b)
	}

	return &blobList
}

func (bl *BlobList) ReadBlobsFromLocalFolder(i indexers.INDEXERS) ([]*Blob, error) {
	var blobs []*Blob
	err := filepath.WalkDir(filepath.Join(blobsFilePath, i.GetIndexer()), func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		log.Fatalln("PANIC PATH: ", path)
		b, err := bl.ReadSpecificBlobFromLocalFolder(i, path)
		if err != nil {
			return err
		}
		blobs = append(blobs, b)
		return nil
	})

	return blobs, err
}

func (bl *BlobList) ReadSpecificBlobFromLocalFolder(i indexers.INDEXERS, path string) (*Blob, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// create blob
	b := CreateBlob()
	b.Folder = i.GetIndexer()

	// separate the header from the body
	before, after, ok := bytes.Cut(file, []byte("\n"))
	if !ok {
		return nil, errors.New("there's no newlines in the file")
	}

	header := before
	body := after

	// we retrieve the header "metadata"
	headerData := strings.Split(string(header), ",")
	if len(headerData) != 4 {
		return nil, fmt.Errorf("header length should be 4, received %d", len(header))
	}

	b.Title = headerData[HeaderTitle]
	b.Description = headerData[HeaderDescription]
	b.URL = headerData[HeaderURL]

	b.Body = body
	b.SetDateTime(headerData[HeaderDatetime])

	// lastly we stem its content
	b.StemWords(string(body))

	return b, nil
}

func (bl *BlobList) GetBlobPath(i indexers.INDEXERS, name string) string {
	return filepath.Join(filepath.Join(blobsFilePath, i.GetIndexer(), name))
}
