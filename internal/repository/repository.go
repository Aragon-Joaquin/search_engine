package repository

import (
	"context"

	"search_engine/internal/blobs"
	"search_engine/internal/db"
	"search_engine/internal/utils"
)

type Repository struct {
	db *db.RedisClient
}

func CreateRepository(db *db.RedisClient) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) UserMakeQuery(word string) []*blobs.Blob {
	query := blobs.CreateBlob()
	query.StemWords(word)

	// TODO: finish
	// ctx := context.Background()
	// bList, err := r.db.GetAllZBlobs(ctx)
	// if err != nil {
	// 	panic(err)
	// }
	// return bList.Calculate_tf_idf(query)
	return []*blobs.Blob{}
}

// both saves into the redisdb + local
func (r *Repository) SaveBlob(b *blobs.Blob, i utils.INDEXERS, content *[]byte) error {
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
