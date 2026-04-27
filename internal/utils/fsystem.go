package utils

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
)

const (
	BLOBS_FOLDER = "/data"
)

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Join(filepath.Dir(b), BLOBS_FOLDER)

	err_empty_filename = errors.New("empty filname")
)

func CreateFile(folder INDEXERS, filename string) (file *os.File, err error) {
	ffolder := string(folder)
	if ffolder == "" {
		return nil, err_empty_filename
	}

	indexerFolderPath := filepath.Join(basepath, ffolder)

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

func FindFileName(folder INDEXERS, filename string) (pathFound string, err error) {
	if filename == "" {
		return "", err_empty_filename
	}

	ffolder := string(folder)
	if ffolder == "" {
		return "", err_empty_filename
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
