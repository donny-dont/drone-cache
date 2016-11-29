package main

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
)

// Plugin contains the common values for the cache.
type Plugin struct {
	Filename string
	Path     string
}

func NewPlugin(c *cli.Context) (Plugin, error) {
	// Get the path for the cache files
	path := c.String("path")

	// Defaults to <owner>/<repo>/<branch>/
	if len(path) == 0 {
		log.Info("No path specified. Creating default")

		path = fmt.Sprintf(
			"/%s/%s/%s/",
			c.GlobalString("repo.owner"),
			c.GlobalString("repo.name"),
			c.GlobalString("commit.branch"),
		)
	}

	p := Plugin{
		Filename: c.GlobalString("filename"),
		Path:     path,
	}

	return p, nil
}
