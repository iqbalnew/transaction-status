package miniowrapper

import (
	"context"
	"io"

	"github.com/minio/minio-go/v7"
)

type MinioWrapper struct {
	MinioInterface
}

func (rw *MinioWrapper) New(endpoint string, opts *minio.Options) (*minio.Client, error) {
	return minio.New(endpoint, opts)
}

func (rw *MinioWrapper) GetObject(ctx context.Context, minioClient *minio.Client, bucketName, objectName string, opts minio.GetObjectOptions) (*minio.Object, error) {
	return minioClient.GetObject(ctx, bucketName, objectName, opts)
}

func (rw *MinioWrapper) PutObject(ctx context.Context, minioClient *minio.Client, bucketName, objectName string, reader io.Reader, objectSize int64, opts minio.PutObjectOptions) (info minio.UploadInfo, err error) {
	return minioClient.PutObject(ctx, bucketName, objectName, reader, objectSize, opts)
}

func (rw *MinioWrapper) RemoveObject(ctx context.Context, minioClient *minio.Client, bucketName, objectName string, opts minio.RemoveObjectOptions) error {
	return minioClient.RemoveObject(ctx, bucketName, objectName, opts)
}
