package plugin

import (
	"github.com/drone-plugins/drone-cache/cache"
	"github.com/drone-plugins/drone-cache/storage"
	"github.com/urfave/cli"
)

const (
	restoreMode = "restore"
	rebuildMode = "rebuild"
)

func Exec(c *cli.Context, s storage.Storage) error {
	p, err := newPlugin(c)

	if err != nil {
		return err
	}

	ca, err := cache.NewCache(s)

	if err != nil {
		return err
	}

	path := p.Path + p.Filename

	if p.Mode == rebuildMode {
		return ca.Rebuild(p.Mount, path)
	}

	return ca.Restore(path)
}

type plugin struct {
	Filename string
	Path     string
	Mode     string
	Mount    string
}
