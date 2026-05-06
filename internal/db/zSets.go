package db

import (
	"context"
	"fmt"
	"log"
	"sync"

	"search_engine/internal/blobs"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Db *redis.Client
}

func zortedKey(idTitle string) string {
	return fmt.Sprintf("%s:%s", ZSET, idTitle)
}

// TODO: use sortedSets for documents and maybe a set for stopWords
// and a HASHES for metadata
func (r *RedisClient) SaveBlobsToRedis(ctx context.Context, blob *blobs.Blob) error {
	if len(blob.TermSpace) == 0 {
		return fmt.Errorf("not enough termSpace size")
	}

	zortedSet := []redis.Z{}

	for m, s := range blob.TermSpace {
		zortedSet = append(zortedSet, redis.Z{Score: float64(s), Member: m})
	}

	id := r.GetBlobUniqueIdentifier(blob)

	if err := r.Db.ZAdd(ctx, zortedKey(id), zortedSet...).Err(); err != nil {
		return err
	}

	if err := r.SetMetaData(ctx, blob); err != nil {
		return err
	}

	return nil
}

func (r *RedisClient) GetZSort(ctx context.Context, idTitle uuid.UUID) (*[]redis.Z, error) {
	res := r.Db.ZRangeWithScores(ctx, zortedKey(idTitle.String()), 0, -1)

	if res.Err() != nil {
		return nil, res.Err()
	}

	results := res.Val()
	return &results, nil
}

func (r *RedisClient) GetAllZBlobs(ctx context.Context) (*blobs.BlobList, error) {
	names, err := r.GetAllBlobsNames(ctx, HASH, 0)
	if err != nil {
		return nil, err
	}

	blist := blobs.CreateBlobList()
	for _, title := range names {
		// todo: make go func
		// pipe := r.Db.TxPipeline()

		res, err := r.Db.ZRangeWithScores(ctx, zortedKey(title), 0, -1).Result()
		if err != nil {
			log.Println("failed searching the termSpace")
			continue
		}

		var redisblob blobs.RedisBlob
		if err := r.Db.HGetAll(ctx, hashKey(title)).Scan(&redisblob); err != nil {
			log.Println("failed while scanning the blob")
			continue
		}
		blob := redisblob.TransformToBlob()

		// if _, err := pipe.Exec(ctx); err != nil {
		// 	log.Println("failed on the execution of the pipeline ")
		// 	continue
		// }

		for _, w := range res {
			if val, ok := w.Member.(string); ok {
				blob.AddToTermSpace(val, int(w.Score))
			}
		}

		blist.AppendBlob(blob)
	}

	return blist, nil
}

// what this does is:
// > query a certain limit of zortedSets
// > check each member of the termSpace of the blob to rank them
// > select until the top X is satisfied
// > query their metadata
// > transform them into blobs
// > return them
func (r *RedisClient) FilterByTermInSpace(ctx context.Context, blobTerm *blobs.Blob) (*blobs.BlobList, error) {
	var wg sync.WaitGroup
	var limitCount int64 = 50 // can be increased

	zNames, err := r.GetAllBlobsNames(ctx, ZSET, 0, limitCount)
	if err != nil {
		return nil, err
	}

	tx := r.Db.TxPipeline()
	for range zNames {
		wg.Go(func() {
			// todo
		})
	}

	tx.Exec(ctx)

	wg.Wait()
	return nil, nil
}
