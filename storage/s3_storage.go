package storage

import (
	"fmt"
	"io"
	"strings"

	log "github.com/Sirupsen/logrus"
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

	UseSSL bool
}

type s3Storage struct {
	client *minio.Client
	opts   *S3Options
}

// NewS3Storage creates an implementation of Storage with S3 as the backend.
func NewS3Storage(opts *S3Options) (Storage, error) {
	client, err := minio.New(opts.Endpoint, opts.Access, opts.Secret, opts.UseSSL)

	if err != nil {
		return nil, err
	}

	return &s3Storage{
		client: client,
		opts:   opts,
	}, nil
}

func (s *s3Storage) Get(p string, dst io.Writer) error {
	bucket, key := splitBucket(p)

	if len(bucket) == 0 || len(key) == 0 {
		return fmt.Errorf("Invalid path %s", p)
	}

	log.Infof("Retrieving file in %s at %s", bucket, key)

	exists, err := s.client.BucketExists(bucket)

	if !exists {
		return err
	}

	object, err := s.client.GetObject(bucket, key)
	if err != nil {
		return err
	}

	numBytes, err := io.Copy(dst, object)

	if err != nil {
		return err
	}

	log.Infof("Downloaded %s from server", byteSize(numBytes))

	return nil
}

func (s *s3Storage) Put(p string, src io.Reader) error {
	bucket, key := splitBucket(p)

	if len(bucket) == 0 || len(key) == 0 {
		return fmt.Errorf("Invalid path %s", p)
	}

	exists, err := s.client.BucketExists(bucket)

	if !exists || err != nil {
		if err = s.client.MakeBucket(bucket, s.opts.Region); err != nil {
			return err
		}
		log.Debugf("Bucket %s created", bucket)
	} else {
		log.Debugf("Bucket %s already exists", bucket)
	}

	log.Debugf("Putting file in %s at %s", bucket, key)

	numBytes, err := s.client.PutObject(bucket, key, src, "application/tar")

	if err != nil {
		return err
	}

	log.Infof("Uploaded %s to server", byteSize(numBytes))

	return nil
}

func splitBucket(p string) (string, string) {
	// Remove initial forward slash
	full := strings.TrimPrefix(p, "/")

	// Get first index
	i := strings.Index(full, "/")

	if i != -1 && len(full) != i+1 {
		return full[0:i], full[i+1:]
	}

	return "", ""
}
