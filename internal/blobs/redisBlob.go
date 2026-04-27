package blobs

import (
	"time"

	"search_engine/internal/utils"
)

type RedisBlob struct {
	Title       string    `redis:"title"`
	Description string    `redis:"description"`
	Datetime    time.Time `redis:"datetime"`
	URL         string    `redis:"url"`

	Length int    `redis:"length"`
	Folder string `redis:"folder"`
}

func (r *RedisBlob) TransformToBlob() *Blob {
	blob := CreateBlob()
	blob.Title = r.Title

	if r.Description != "" {
		blob.Description = r.Description
	} else {
		blob.Description = "No description provided"
	}

	blob.Datetime = r.Datetime
	blob.URL = r.URL

	blob.Length = r.Length
	blob.Folder = utils.INDEXERS(r.Folder)

	return blob
}

func (b *Blob) ParseToRedisBlob() *RedisBlob {
	return &RedisBlob{
		Title:       b.Title,
		Description: b.Description,
		Datetime:    b.Datetime,
		URL:         b.URL,
		Length:      b.Length,
		Folder:      string(b.Folder),
	}
}
