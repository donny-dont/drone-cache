package cache

import (
	"fmt"
	"io"

	"github.com/drone-plugins/drone-cache/archive"
	"github.com/drone-plugins/drone-cache/storage"
)

type Build struct {
	Owner  string
	Repo   string
	Branch string
}

func RebuildCache(b Build, s storage.Storage, a archive.Archive) error {
	path := buildPath(b.Owner, b.Repo, b.Branch)

	return rebuildCache("testing/cache.go", path, s, a)
}

func RestoreCache(b Build, s storage.Storage, a archive.Archive) error {
	path := buildPath(b.Owner, b.Repo, b.Branch)

	err := restoreCache(path, s, a)

	if err != nil {
		// Attempt fallback
		path = buildPath(b.Owner, b.Repo, "master")

		if err = restoreCache(path, s, a); err != nil {
			return err
		}
		fmt.Printf("Using cache from fallback branch %s", "master")
	} else {
		fmt.Printf("Using cache on branch %s\n", b.Branch)
	}

	return nil
}

func restoreCache(src string, s storage.Storage, a archive.Archive) error {
	reader, writer := io.Pipe()

	go func() {
		defer writer.Close()

		s.Get(src, writer)
	}()

	return a.Unpack(src, reader)
}

func rebuildCache(src string, dest string, s storage.Storage, a archive.Archive) error {
	reader, writer := io.Pipe()

	go func() {
		defer writer.Close()

		a.Pack(src, writer)
	}()

	return s.Put(dest, reader)
}

func buildPath(owner string, repo string, branch string) string {
	return fmt.Sprintf(
		"%s-%s/%s/archive.tar",
		owner,
		repo,
		branch,
	)
}
