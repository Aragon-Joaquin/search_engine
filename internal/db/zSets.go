package db

import (
	"context"
	"fmt"
	"log"
	"sync"

	"search_engine/internal/blobs"
	"search_engine/internal/indexers"

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

		res, err := r.Db.ZRangeWithScores(ctx, zortedKey(title.Name), 0, -1).Result()
		if err != nil {
			log.Println("failed searching the termSpace")
			continue
		}

		var redisblob blobs.RedisBlob
		if err := r.Db.HGetAll(ctx, hashKey(title.Name)).Scan(&redisblob); err != nil {
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
				blob.AddToTermSpace(val, uint64(w.Score))
			}
		}

		blist.AppendBlob(blob)
	}

	return blist, nil
}

// what this does is:
// > query a certain limit of zortedSets
// > transform them into blobs
// > check their local blob
// > rank it
// > append it to the blobList
// TODO: simplify it
func (r *RedisClient) FilterByTermInSpace(ctx context.Context, blobTerm *blobs.Blob) (*blobs.BlobList, error) {
	var wg sync.WaitGroup
	var limitCount int64 = 50 // can be increased

	zNames, err := r.GetAllBlobsNames(ctx, ZSET, 0, limitCount)
	if err != nil {
		return nil, err
	}

	blobWords := []string{}
	for w := range blobTerm.TermSpace {
		blobWords = append(blobWords, w)
	}

	blobList := blobs.CreateBlobList()
	if len(blobWords) == 0 {
		return blobList, nil
	}

	tx := r.Db.TxPipeline()
	for _, b := range zNames {
		wg.Go(func() {
			tx.ZMScore(ctx, b.FullName(), blobWords...)
		})
	}
	wg.Wait()

	cmder, err := tx.Exec(ctx)
	if err != nil {
		return nil, err
	}

	for _, k := range cmder {
		wg.Go(func() {
			res, ok := k.(*redis.FloatSliceCmd)

			if !ok || k.Err() != nil {
				return
			}

			args := res.Args() // [type, blobname, SATURN, PLANET, etc...]
			if len(args) < 2 {
				return
			}

			// blobname
			n := args[1].(string)
			blobName, err := r.GetBlobFolderAndTitleFromIdentifier(n)
			if err != nil {
				return
			}

			// finding blobname in local folder
			b, err := blobList.ReadSpecificBlobFromLocalFolder(
				indexers.GetWikipediaIndexer(),
				blobList.GetBlobPath(indexers.GetWikipediaIndexer(), blobName.Name),
			)
			if err != nil {
				return
			}

			messages := args[2:] // [SATURN, PLANET, etc...]
			lengthMsg := len(messages)
			for i, v := range res.Val() {
				if i >= lengthMsg {
					break
				}
				b.AddToTermSpace(messages[i].(string), uint64(v))
			}

			b.CalculateDotProduct(blobTerm)
			if b.Score < indexers.MIN_SCORE_THRESHOLD {
				return
			}

			blobList.AppendBlob(b)
		})
	}

	wg.Wait()
	return blobList, nil
}
