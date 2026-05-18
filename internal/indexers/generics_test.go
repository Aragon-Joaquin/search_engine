package indexers

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWorkingDirectory(T *testing.T) {
	dir, err := INDEXER_WIKIPEDIA.AppendToWorkingDirectory()
	if err != nil {
		T.Errorf("Couldn't get the WKDirectory: %e\n", err)
		return
	}

	currentWd, err := os.Getwd()
	if err != nil {
		T.Errorf("Couldn't get the OS WKDirectory: %e\n", err)
		return
	}

	fb := filepath.Join(currentWd, "/data", INDEXER_WIKIPEDIA.Indexer)

	if fb != dir {
		T.Errorf("Distinct values between \nfb: %s\ndir: %s\n", fb, dir)
		return

	}
}
