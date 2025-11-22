package minio

import (
	"context"
	"io"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Minio struct {
	client *minio.Client
	ctx    context.Context
}

func NewMinio(ctx context.Context) *Minio {
	client, err := minio.New("localhost:9000", &minio.Options{
		Creds:  credentials.NewStaticV4("admin", "admin123", ""),
		Secure: false,
	})
	if err != nil {
		return nil
	}
	return &Minio{client: client, ctx: ctx}
}

func (m *Minio) CreateBucket(bucketName string) {
	if err := m.client.MakeBucket(m.ctx, bucketName, minio.MakeBucketOptions{}); err != nil {
		exists, errBucketExists := m.client.BucketExists(m.ctx, bucketName)
		if errBucketExists == nil && exists {
			log.Println("Bucket already exists")
		} else {
			log.Fatalln(err)
		}
	}
}

func (m *Minio) UploadFile(file io.Reader, size int64, filename string) error {
	_, err := m.client.PutObject(m.ctx, "videos", filename, file, size, minio.PutObjectOptions{ContentType: "video/mp4"})
	if err != nil {
		return err
	}

	return nil
}
