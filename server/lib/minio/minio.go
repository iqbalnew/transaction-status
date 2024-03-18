package minio

import (
	"bytes"
	"context"
	"errors"
	"io"
	"time"

	miniowrapper "bitbucket.bri.co.id/scm/bricams-addons/qcash-template-service/server/lib/minio/wrapper"
	"github.com/gabriel-vasile/mimetype"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Config struct {
	Bucket          string
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	UseSsl          bool
	Timeout         int64
}

type MinioObj struct {
	client *minio.Client
	config *Config
	mcw    miniowrapper.MinioInterface
}

func NewMinio(config *Config, mcw miniowrapper.MinioInterface) (*MinioObj, error) {
	if mcw == nil {
		return nil, errors.New("failed create minio object")
	}

	return &MinioObj{
		config: config,
		mcw:    mcw,
	}, nil
}

func (mo *MinioObj) Connect() error {
	var err error
	creds := credentials.NewStaticV4(mo.config.AccessKeyID, mo.config.SecretAccessKey, "")
	mo.client, err = mo.mcw.New(mo.config.Endpoint, &minio.Options{
		Creds:  creds,
		Secure: mo.config.UseSsl,
	})

	return err
}

func (mo *MinioObj) ReadFile(ctx context.Context, fileName string) ([]byte, error) {
	ctxTimeout, ctxCanc := context.WithTimeout(ctx, time.Duration(mo.config.Timeout)*time.Second)
	defer ctxCanc()

	object, getObjErr := mo.mcw.GetObject(ctxTimeout, mo.client, mo.config.Bucket, fileName, minio.GetObjectOptions{})
	if getObjErr != nil {
		return nil, getObjErr
	}
	defer object.Close()

	return io.ReadAll(object)
}

func (mo *MinioObj) UploadFile(ctx context.Context, objectName string, data []byte) (*minio.UploadInfo, error) {
	ctxTimeout, ctxCancel := context.WithTimeout(ctx, time.Duration(mo.config.Timeout)*time.Second)
	defer ctxCancel()

	reader := bytes.NewReader(data)
	mtype, err := mimetype.DetectReader(reader)
	if err != nil {
		return nil, err
	}

	fileInfo, putErr := mo.mcw.PutObject(ctxTimeout, mo.client, mo.config.Bucket, objectName, reader, reader.Size(),
		minio.PutObjectOptions{
			ContentType: mtype.String(),
		})

	if putErr != nil {
		return nil, putErr
	}

	return &fileInfo, nil
}

func (mo *MinioObj) DeleteFile(ctx context.Context, fileName string) error {
	ctxTimeout, ctxCancel := context.WithTimeout(ctx, time.Duration(mo.config.Timeout)*time.Second)
	defer ctxCancel()

	return mo.mcw.RemoveObject(ctxTimeout, mo.client, mo.config.Bucket, fileName, minio.RemoveObjectOptions{})
}
