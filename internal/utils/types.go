package utils

import (
	"fmt"
	"net/url"
)

type INDEXERS string

const (
	INDEXER_WIKIPEDIA INDEXERS = "en.wikipedia.org/wiki"
)

func GetURL(indexer INDEXERS, term string) string {
	termEncoded := url.QueryEscape(term)
	return fmt.Sprintf("https://%s/%s", indexer, termEncoded)
}

func GetAbsoluteIndexerURL(ind INDEXERS) string {
	return fmt.Sprintf("https://%s", ind)
}
