package blobs

import (
	"time"

	"github.com/google/uuid"
)

type RedisBlob struct {
	Title       string    `redis:"title"`
	Description string    `redis:"description"`
	Datetime    time.Time `redis:"datetime"`
	URL         string    `redis:"url"`

	Length int    `redis:"length"`
	UUID   string `redis:"uuid"`
}

func (r *RedisBlob) TransformToBlob() *Blob {
	res, err := uuid.Parse(r.UUID)
	if err != nil {
		res = uuid.Nil
	}

	blob := CreateBlob()

	blob.Title = r.Title
	blob.Description = r.Description
	blob.Datetime = r.Datetime
	blob.URL = r.URL

	blob.Length = r.Length
	blob.UUID = res

	return blob
}

func (b *Blob) SaveBlobInformation() *RedisBlob {
	return &RedisBlob{
		Title:       b.Title,
		Description: b.Description,
		Datetime:    b.Datetime,
		URL:         b.URL,
		Length:      b.Length,
		UUID:        b.GetUUID(),
	}
}
