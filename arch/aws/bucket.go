package aws

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"time"

	awsService "github.com/aws/aws-sdk-go/aws"
	awsServiceS3 "github.com/aws/aws-sdk-go/service/s3"
)

type Bucket interface {
	GetInstance() *bucket
	PutObject(bucket string, key string, file *multipart.FileHeader) error
	PresignedURL(bucket string, key string, expires time.Duration) (string, error)
}

type bucket struct {
	context context.Context
	session Session
}

func NewBucket(session Session) Bucket {
	return &bucket{
		context: context.Background(),
		session: session,
	}
}

func (b *bucket) GetInstance() *bucket {
	return b
}

func (b *bucket) PutObject(bucket string, key string, file *multipart.FileHeader) error {
	session, err := b.session.GetSession()
	if err != nil {
		return err
	}

	svc := awsServiceS3.New(session)
	_, err = svc.PutObject(&awsServiceS3.PutObjectInput{
		Bucket: awsService.String(bucket),
		Key:    awsService.String(key),
		Body: func() *bytes.Reader {
			f, _ := file.Open()
			data, _ := io.ReadAll(f)
			defer f.Close()
			return bytes.NewReader(data)
		}(),
	})
	if err != nil {
		return err
	}

	return nil
}

func (b *bucket) PresignedURL(bucket string, key string, expires time.Duration) (string, error) {
	session, err := b.session.GetSession()
	if err != nil {
		return "", err
	}

	svc := awsServiceS3.New(session)

	result, _ := svc.GetObjectRequest(&awsServiceS3.GetObjectInput{
		Bucket: awsService.String(bucket),
		Key:    awsService.String(key),
	})

	url, err := result.Presign(expires)
	if err != nil {
		return "", err
	}

	return url, nil
}
