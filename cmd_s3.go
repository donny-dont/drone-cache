package main

import (
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/drone-plugins/drone-cache/archive"
	"github.com/drone-plugins/drone-cache/cache"
	"github.com/drone-plugins/drone-cache/storage"
	"github.com/urfave/cli"
)

var s3Cmd = cli.Command{
	Name:   "s3",
	Usage:  "cache to S3",
	Action: s3Plugin,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:   "server",
			Usage:  "s3 server",
			EnvVar: "PLUGIN_SERVER,CACHE_S3_SERVER",
		},
		cli.StringFlag{
			Name:   "access-key",
			Usage:  "s3 access key",
			EnvVar: "PLUGIN_ACCESS_KEY,CACHE_S3_ACCESS_KEY",
		},
		cli.StringFlag{
			Name:   "secret-key",
			Usage:  "s3 secret key",
			EnvVar: "PLUGIN_SECRET_KEY,CACHE_S3_SECRET_KEY",
		},
	},
}

func s3Options(c *cli.Context) (*storage.S3Options, error) {
	// Get the endpoint
	server := c.String("server")

	var endpoint string
	var useSSL bool

	if len(server) > 0 {
		useSSL = strings.HasPrefix(server, "https://")

		if !useSSL {
			if !strings.HasPrefix(server, "http://") {
				return nil, fmt.Errorf("Invalid server %s. Needs to be a HTTP URI", server)
			}

			endpoint = server[7:]
		} else {
			endpoint = server[8:]
		}
	} else {
		endpoint = "s3.amazonaws.com"
		useSSL = true
	}

	// Get the access credentials
	access := c.String("access-key")
	secret := c.String("secret-key")

	if len(access) == 0 || len(secret) == 0 {
		return nil, fmt.Errorf("No access credentials provided")
	}

	return &storage.S3Options{
		Endpoint: endpoint,
		Access:   access,
		Secret:   secret,
		UseSSL:   useSSL,
	}, nil
}

func s3Plugin(c *cli.Context) error {
	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)

	plugin, err := NewPlugin(c)

	if err != nil {
		return err
	}

	opts, err := s3Options(c)

	if err != nil {
		return err
	}

	log.Infof("Using %s as the cache", opts.Endpoint)
	log.Infof("Using %s", plugin.Path)

	s, err := storage.NewS3Storage(opts)

	if err != nil {
		return err
	}

	a, err := archive.FromFilename("archive.tar")

	if err != nil {
		return err
	}

	b := cache.Build{
		Owner:  "drone",
		Repo:   "drone-cache",
		Branch: "foo/bar",
	}

	err = cache.RestoreCache(b, s, a)

	return err
}
