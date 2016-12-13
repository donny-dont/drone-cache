package plugin

import (
	"errors"
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
)

func PluginFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:   "filename",
			Usage:  "Filename for the cache",
			EnvVar: "PLUGIN_FILENAME",
		},
		cli.StringFlag{
			Name:   "path",
			Usage:  "path",
			EnvVar: "PLUGIN_PATH",
		},
		cli.StringSliceFlag{
			Name:   "mount",
			Usage:  "cache directories",
			EnvVar: "PLUGIN_MOUNT",
		},
		cli.BoolFlag{
			Name:   "rebuild",
			Usage:  "rebuild the cache directories",
			EnvVar: "PLUGIN_REBUILD",
		},
		cli.BoolFlag{
			Name:   "restore",
			Usage:  "restore the cache directories",
			EnvVar: "PLUGIN_RESTORE",
		},

		cli.BoolFlag{
			Name:   "debug",
			Usage:  "debug plugin output",
			EnvVar: "PLUGIN_DEBUG",
		},

		// Build information

		cli.StringFlag{
			Name:   "repo.owner",
			Usage:  "repository owner",
			EnvVar: "DRONE_REPO_OWNER",
		},
		cli.StringFlag{
			Name:   "repo.name",
			Usage:  "repository name",
			EnvVar: "DRONE_REPO_NAME",
		},
		cli.StringFlag{
			Name:   "commit.branch",
			Value:  "master",
			Usage:  "git commit branch",
			EnvVar: "DRONE_COMMIT_BRANCH",
		},
	}
}

func newPlugin(c *cli.Context) (*plugin, error) {
	if c.GlobalBool("debug") {
		log.SetLevel(log.DebugLevel)
	}

	// Determine the mode for the plugin
	rebuild := c.GlobalBool("rebuild")
	restore := c.GlobalBool("restore")

	if rebuild && restore {
		return nil, errors.New("Cannot rebuild and restore the cache")
	} else if !rebuild && !restore {
		return nil, errors.New("No action specified")
	}

	var mode string
	var mount []string

	if rebuild {
		// Look for the mount points to rebuild
		mount = c.GlobalStringSlice("mount")

		if len(mount) == 0 {
			return nil, errors.New("No mounts specified")
		}

		mode = rebuildMode
	} else {
		mode = restoreMode
	}

	// Get the path to place the cache files
	path := c.GlobalString("path")

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

	// Get the filename
	filename := c.GlobalString("filename")

	if len(filename) == 0 {
		log.Info("No filename specified. Creating default")

		filename = "archive.tar"
	}

	return &plugin{
		Filename: filename,
		Path:     path,
		Mount:    mount,
		Mode:     mode,
	}, nil
}
