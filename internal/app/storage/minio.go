package storage

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"mime/multipart"
	"os"
)

type MinioStorage struct {
	client *minio.Client
}

func NewMinioStorage(endpoint, accessKey, secretKey string) (*MinioStorage, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(accessKey, secretKey, ""),
	})
	if err != nil {
		return nil, err
	}
	return &MinioStorage{client: client}, nil
}

func (s *MinioStorage) LoadImg(bucketName, fileName string, file multipart.File, fileSize int64) error {
	if bucketName == "" {
		bucketName = os.Getenv("MINIO_BUCKET_NAME")
	}
	_, err := s.client.PutObject(context.Background(), bucketName, fileName, file, fileSize, minio.PutObjectOptions{})
	return err
}

func (s *MinioStorage) DeleteImg(bucketName, fileName string) error {
	if bucketName == "" {
		bucketName = os.Getenv("MINIO_BUCKET_NAME")
	}
	return s.client.RemoveObject(context.Background(), bucketName, fileName, minio.RemoveObjectOptions{})
}
