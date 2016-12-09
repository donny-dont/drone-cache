package storage

import (
	"io"
	"io/ioutil"

	log "github.com/Sirupsen/logrus"
)

// S3Options contains configuration for the S3 connection.
type DummyOptions struct {
  Server   string
	Username string
  Password string
}

type dummyStorage struct {
	opts   *DummyOptions
}

// NewS3Storage creates an implementation of Storage with S3 as the backend.
func NewDummyStorage(opts *DummyOptions) (Storage, error) {
	return &dummyStorage{
		opts:   opts,
	}, nil
}

func (s *dummyStorage) Get(p string, dst io.Writer) error {
	return nil
}

func (s *dummyStorage) Put(p string, src io.Reader) error {
	log.Infof("Reading for %s", p)

	_, err := ioutil.ReadAll(src)

	if err != nil {
		log.Errorf("Failed to read for %s", p)
		return err
	}

	log.Infof("Finished reading for %s", p)

	return nil
}
