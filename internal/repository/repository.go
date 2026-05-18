package repository

import (
	"context"

	"search_engine/internal/blobs"
	"search_engine/internal/db"
	"search_engine/internal/indexers"
)

type Repository struct {
	db *db.RedisClient
}

func CreateRepository(db *db.RedisClient) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) UserMakeQuery(userQuery *blobs.Blob) (*blobs.BlobList, error) {
	ctx := context.Background()
	bList, err := r.db.FilterByTermInSpace(ctx, userQuery)
	if err != nil {
		return nil, err
	}

	return bList, nil
}

// both saves into the redisdb + local
func (r *Repository) SaveBlob(b *blobs.Blob, i indexers.INDEXERS, content *[]byte) error {
	f, err := i.CreateFile(b.Title)
	if err != nil {
		return err
	}

	if err := b.ParseBlobContentsIntoFile(f, content); err != nil {
		return err
	}

	ctx := context.Background()
	return r.db.SaveBlobsToRedis(ctx, b)
}
