package indexers

import "os"

type types_of_indexers struct {
	SearchTerm string
	Indexer    string
}

// specific
type INDEXERS interface {
	Search(term string) []string
	GetAbsoluteIndexerURL() string
	AppendToWorkingDirectory() (string, error)
	GetIndexer() string

	FindFileName(filename string) (pathFound string, err error)
	CreateFile(filename string) (file *os.File, err error)
}

func GetWikipediaIndexer() INDEXERS {
	return INDEXER_WIKIPEDIA
}
