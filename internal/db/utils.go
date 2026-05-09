package db

import (
	"errors"
	"fmt"
	"strings"

	"search_engine/internal/blobs"
)

type DATA_REDIS string

const (
	ZSET DATA_REDIS = "zset"
	HASH DATA_REDIS = "hash"
)

func (r *RedisClient) GetBlobUniqueIdentifier(blob *blobs.Blob) string {
	res := string(blob.Folder)
	return fmt.Sprintf("%s|%s", res, blob.Title)
}

type BlobName struct {
	TypeBlob DATA_REDIS
	Folder   string
	Name     string
}

func (r *RedisClient) GetBlobFolderAndTitleFromIdentifier(id string) (*BlobName, error) {
	if len(id) < 3 {
		return nil, errors.New("insufficient length for the string")
	}

	bName := BlobName{}

	// we split first for the "hash:" or "zset:"
	before, after, found := strings.Cut(id, ":")
	if !found {
		return nil, errors.New("could find the datatype")
	}

	bName.TypeBlob = DATA_REDIS(before)

	nA, nB, nF := strings.Cut(after, "|")

	if !nF {
		return nil, errors.New("couldnt find both name and folder")
	}

	bName.Folder = nA
	bName.Name = nB

	return &bName, nil
}

func (b *BlobName) FullName() string {
	return fmt.Sprintf("%s:%s|%s", b.TypeBlob, b.Folder, b.Name)
}
