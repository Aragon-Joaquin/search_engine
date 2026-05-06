package db

import (
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

func (r *RedisClient) GetBlobFolderAndTitleFromIdentifier(id string) (typeblob string, folder string, name string) {
	// we split first for the "hash:" or "zset:"
	c := strings.SplitN(id, ":", 1)
	typeBlob := c[0]

	res := strings.SplitN(c[1], "|", 1)
	folder = res[0]
	name = res[1]

	return typeBlob, folder, name
}
