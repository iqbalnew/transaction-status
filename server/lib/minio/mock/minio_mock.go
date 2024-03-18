package miniomock

import (
	"context"
	"errors"
	"io"

	miniowrapper "bitbucket.bri.co.id/scm/bricams-addons/qcash-template-service/server/lib/minio/wrapper"
	"github.com/minio/minio-go/v7"
)

type MinioMock struct {
	miniowrapper.MinioInterface
}

func (mm *MinioMock) New(endpoint string, opts *minio.Options) (*minio.Client, error) {
	if endpoint == "failed-connect" {
		return nil, errors.New("failed to connect minio")
	}

	return &minio.Client{}, nil
}

func (mm *MinioMock) GetObject(ctx context.Context, minioClient *minio.Client, bucketName, objectName string, opts minio.GetObjectOptions) (*minio.Object, error) {
	if objectName == "failed-object.xlsx" {
		return nil, errors.New("failed get object")
	}

	return &minio.Object{}, nil
}

func (mm *MinioMock) PutObject(ctx context.Context, minioClient *minio.Client, bucketName, objectName string, reader io.Reader, objectSize int64, opts minio.PutObjectOptions) (info minio.UploadInfo, err error) {
	if objectName == "failed-object.xlsx" {
		return minio.UploadInfo{
			Bucket: bucketName,
		}, errors.New("failed put object")
	}

	return minio.UploadInfo{
		Bucket: bucketName,
	}, nil
}

func (mm *MinioMock) RemoveObject(ctx context.Context, minioClient *minio.Client, bucketName, objectName string, opts minio.RemoveObjectOptions) error {
	if objectName == "failed-object.xlsx" {
		return errors.New("failed remove object")
	}

	return nil
}
