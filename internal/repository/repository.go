package repository

import (
	"context"

	"search_engine/internal/blobs"
	"search_engine/internal/db"
)

type Repository struct {
	db *db.RedisClient
}

func CreateRepostory(db *db.RedisClient) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) UserMakeQuery(word string) []*blobs.Blob {
	query := blobs.CreateBlob()
	query.StemWords(word)

	ctx := context.Background()
	bList, err := r.db.GetAllZBlobs(ctx)
	if err != nil {
		panic(err)
	}

	return bList.Calculate_tf_idf(query)
}
