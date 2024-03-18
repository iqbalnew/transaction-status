package miniowrapper

import (
	"context"
	"io"

	"github.com/minio/minio-go/v7"
)

type MinioInterface interface {
	New(endpoint string, opts *minio.Options) (*minio.Client, error)
	GetObject(ctx context.Context, minioClient *minio.Client, bucketName, objectName string, opts minio.GetObjectOptions) (*minio.Object, error)
	PutObject(ctx context.Context, minioClient *minio.Client, bucketName, objectName string, reader io.Reader, objectSize int64, opts minio.PutObjectOptions) (info minio.UploadInfo, err error)
	RemoveObject(ctx context.Context, minioClient *minio.Client, bucketName, objectName string, opts minio.RemoveObjectOptions) error
}
