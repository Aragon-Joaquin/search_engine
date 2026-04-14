package db

import (
	"context"
	"fmt"

	"search_engine/internal/blobs"

	"github.com/google/uuid"
)

const (
	SET_NAMES_KEY = "set_info:list_names"
)

func hashKey(idTitle string) string {
	return fmt.Sprintf("hash:%s", idTitle)
}

func (r *RedisClient) SetMetaData(ctx context.Context, blob *blobs.Blob) error {
	if err := r.Db.HSet(ctx, hashKey(blob.GetUUID()), blob.SaveBlobInformation()).Err(); err != nil {
		return err
	}

	if err := r.Db.SAdd(ctx, SET_NAMES_KEY, blob.GetUUID()).Err(); err != nil {
		return err
	}

	return nil
}

func (r *RedisClient) GetMetaData(ctx context.Context, idTitle uuid.UUID) (*blobs.Blob, error) {
	var redisblob blobs.RedisBlob

	if err := r.Db.HGetAll(ctx, idTitle.String()).Scan(redisblob); err != nil {
		return nil, err
	}

	blob := redisblob.TransformToBlob()
	return blob, nil
}

func (r *RedisClient) GetAllBlobsNames(ctx context.Context) (*[]string, error) {
	s, err := r.Db.SMembers(ctx, SET_NAMES_KEY).Result()
	if err != nil {
		return nil, err
	}
	return &s, nil
}
