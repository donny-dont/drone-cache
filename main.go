package main

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
)

var build = "0" // build number set at compile-time

func main() {
	app := cli.NewApp()
	app.Name = "cache plugin"
	app.Usage = "cache plugin"
	app.Version = fmt.Sprintf("1.0.%s", build)
	app.Flags = []cli.Flag{
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
	app.Commands = []cli.Command{
		s3Cmd,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
