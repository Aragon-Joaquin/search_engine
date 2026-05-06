package utils

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
)

type INDEXERS struct {
	SearchTerm string
	Indexer    string
}

var INDEXER_WIKIPEDIA INDEXERS = INDEXERS{
	Indexer:    "en.wikipedia.org",
	SearchTerm: "/wiki",
}

func (i INDEXERS) Search(term string) string {
	termEncoded := url.QueryEscape(term)
	return fmt.Sprintf("https://%s/%s/%s", i.Indexer, i.SearchTerm, termEncoded)
}

func (i INDEXERS) GetAbsoluteIndexerURL() string {
	return fmt.Sprintf("https://%s", i.Indexer)
}

func (i INDEXERS) AppendToWorkingDirectory() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, BLOBS_FOLDER, i.Indexer), nil
}
