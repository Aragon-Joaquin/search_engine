package utils

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

const (
	BLOBS_FOLDER = "/data"
)

var err_empty_filename = errors.New("empty filname")

func (i INDEXERS) CreateFile(filename string) (file *os.File, err error) {
	ffolder := i.Indexer
	if ffolder == "" {
		return nil, err_empty_filename
	}

	indexerFolderPath, err := i.AppendToWorkingDirectory()
	if err != nil {
		return nil, err
	}

	_, err = os.Stat(indexerFolderPath)
	if err != nil {
		if err2 := os.MkdirAll(indexerFolderPath, 0o755); err2 != nil {
			panic(err2)
		}
	}

	if filename == "" {
		return nil, err_empty_filename
	}

	file, err = os.Create(filepath.Join(indexerFolderPath, filename))
	return file, err
}

func (i INDEXERS) FindFileName(filename string) (pathFound string, err error) {
	if filename == "" {
		return "", err_empty_filename
	}

	ffolder := i.Indexer
	if ffolder == "" {
		return "", err_empty_filename
	}

	basepath, err := i.AppendToWorkingDirectory()
	if err != nil {
		return "", err
	}

	err = filepath.WalkDir(filepath.Join(basepath, ffolder), func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		if d.Name() == filename {
			pathFound = path
			return io.EOF
		}

		return nil
	})

	if err != nil && pathFound == "" {
		return "", err
	}

	return pathFound, nil
}
