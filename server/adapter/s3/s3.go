package s3

import (
	"fmt"
	"io"

	"github.com/Mushus/image-server/server/internal"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type Config interface {
	S3URL() string
	Bucket() string
}

func ProvideStorage(config Config) internal.Storage {
	var endpoint *string

	if config.S3URL() == "" {
		endpoint = aws.String(config.S3URL())
	}

	cfg := &aws.Config{
		Endpoint:         endpoint,
		S3ForcePathStyle: aws.Bool(endpoint != nil),
	}
	sess := session.Must(session.NewSession(cfg))

	instance := s3.New(sess)

	return &Storage{
		bucket: config.Bucket(),
		s3:     instance,
	}
}

type Storage struct {
	bucket string
	s3     *s3.S3
}

var _ internal.Storage = &Storage{}

func (s Storage) Get(path string) (io.ReadCloser, error) {
	obj, err := s.s3.GetObject(&s3.GetObjectInput{
		Key:    aws.String(path),
		Bucket: aws.String(s.bucket),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchKey:
				return nil, internal.ErrNotFound
			}
		}
		return nil, fmt.Errorf("unhandled s3 get error: %w", err)
	}

	return obj.Body, nil
}

func (s Storage) Put(path string, rs io.ReadSeeker) error {
	_, err := s.s3.PutObject(&s3.PutObjectInput{
		Key:    aws.String(path),
		Bucket: aws.String(s.bucket),
		Body:   rs,
	})

	if err != nil {
		return fmt.Errorf("unhandled s3 put error: %w", err)
	}

	return nil
}
