package db

import (
	"context"

	"search_engine/internal/blobs"
)

func (r *RedisClient) UserMakeQuery(word string) []*blobs.Blob {
	query := blobs.CreateBlob()
	query.StemWords(word)

	ctx := context.Background()
	bList, err := r.GetAllZBlobs(ctx)
	if err != nil {
		panic(err)
	}

	return bList.Calculate_tf_idf(query)
}
