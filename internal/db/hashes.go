package db

import (
	"context"
	"fmt"

	"search_engine/internal/blobs"
)

const (
	SET_NAMES_KEY = "set_info:list_names"
)

func hashKey(idTitle string) string {
	return fmt.Sprintf("%s:%s", HASH, idTitle)
}

func (r *RedisClient) SetMetaData(ctx context.Context, blob *blobs.Blob) error {
	id := r.GetBlobUniqueIdentifier(blob)
	if err := r.Db.HSet(ctx, hashKey(id), blob.ParseToRedisBlob()).Err(); err != nil {
		return err
	}
	return nil
}

func (r *RedisClient) GetMetaData(ctx context.Context, blobname string) (*blobs.Blob, error) {
	var redisblob blobs.RedisBlob

	if err := r.Db.HGetAll(ctx, blobname).Scan(redisblob); err != nil {
		return nil, err
	}

	blob := redisblob.TransformToBlob()
	return blob, nil
}

func (r *RedisClient) GetAllBlobsNames(ctx context.Context, typeName DATA_REDIS, start uint64, limit ...int64) ([]*BlobName, error) {
	var defaultLimit int64 = 20
	if len(limit) > 1 {
		defaultLimit = limit[0]
	}

	// really anything, they both share the same name lol
	scan := r.Db.ScanType(ctx, start, "*", defaultLimit, string(typeName))

	keys := []*BlobName{}
	if err := scan.Err(); err != nil {
		return keys, err
	}

	iter := scan.Iterator()
	for iter.Next(ctx) {
		if err := iter.Err(); err != nil {
			return keys, err
		}

		val := iter.Val()
		if names, err := r.GetBlobFolderAndTitleFromIdentifier(val); err == nil {
			keys = append(keys, names)
		}

	}

	return keys, nil
}
