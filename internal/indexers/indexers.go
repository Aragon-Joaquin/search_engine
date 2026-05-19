package indexers

import "os"

type types_of_indexers struct {
	SearchTerm string
	Indexer    string
}

type IndexerExtraInfo struct {
	URL string // absolute url. takes to the resource
	Key string // unique key

	// wikipedia already gives a easy way to gather both of this instead of depending in a css class!
	// other indexers can do the same. Im just using the most of the resources!
	// Obviouly i need to add a todo if u dont mind
	// TODO: handle case when one or all are null fields
	Title       string // title. can be null.
	Description string // description. can be null.
}

// specific
type INDEXERS interface {
	Search(term string) map[string]*IndexerExtraInfo
	GetAbsoluteIndexerURL() string
	AppendToWorkingDirectory() (string, error)
	GetIndexer() string

	FindFileName(filename string) (pathFound string, err error)
	CreateFile(filename string) (file *os.File, err error)
}

func GetWikipediaIndexer() INDEXERS {
	return INDEXER_WIKIPEDIA
}
