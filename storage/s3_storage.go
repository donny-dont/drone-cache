package storage

import (
	"fmt"
	"io"

	"github.com/minio/minio-go"
)

// S3Options contains configuration for the S3 connection.
type S3Options struct {
	Endpoint   string
	Key        string
	Secret     string
	Encryption string
	Access     string

	// us-east-1
	// us-west-1
	// us-west-2
	// eu-west-1
	// ap-southeast-1
	// ap-southeast-2
	// ap-northeast-1
	// sa-east-1
	Region string

	// Use path style instead of domain style.
	//
	// Should be true for minio and false for AWS.
	PathStyle bool
}

type s3Storage struct {
	client *minio.Client
	opts   *S3Options
}

// NewS3Storage creates an implementation of Storage with S3 as the backend.
func NewS3Storage(opts *S3Options) (Storage, error) {
	client, err := minio.New(opts.Endpoint, opts.Access, opts.Secret, false)

	if err != nil {
		return nil, err
	}

	return &s3Storage{
		client: client,
		opts:   opts,
	}, nil
}

func (s *s3Storage) Get(p string, dst io.Writer) error {
	bucket := "again"
	exists, err := s.client.BucketExists(bucket)

	if !exists {
		return err
	}

	object, err := s.client.GetObject(bucket, p)
	if err != nil {
		return err
	}

	numBytes, err := io.Copy(dst, object)

	fmt.Printf("Num bytes written %d", numBytes)

	return err
}

func (s *s3Storage) Put(p string, src io.Reader) error {
	bucket := "again"
	exists, err := s.client.BucketExists(bucket)

	if !exists || err != nil {
		if err = s.client.MakeBucket(bucket, s.opts.Region); err != nil {
			return err
		}
	}

	fmt.Printf("BUCKET EXISTS")

	_, err = s.client.PutObject(bucket, p, src, "application/tar")

	if err != nil {
		fmt.Printf("Could not put object %s\n", err.Error())
	}

	return err
}
