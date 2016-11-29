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
	path := "/drone/drone-cache/master/archive.tar"

	return rebuildCache("testing", path, s, a)
}

func RestoreCache(b Build, s storage.Storage, a archive.Archive) error {
	path := "/drone/drone-cache/master/archive.tar"

	err := restoreCache(path, s, a)

	/*
		if err != nil {
			// Attempt fallback
			path = "/drone/drone-cache/master/archive.tar"

			if err = restoreCache(path, s, a); err != nil {
				return err
			}
			fmt.Printf("Using cache from fallback branch %s", "master")
		} else {
			fmt.Printf("Using cache on branch %s\n", b.Branch)
		}
	*/
	return err
}

func restoreCache(src string, s storage.Storage, a archive.Archive) error {
	reader, writer := io.Pipe()

	cw := make(chan error)

	go func() {
		defer writer.Close()

		cw <- s.Get(src, writer)
	}()

	cr := make(chan error)

	go func() {
		defer reader.Close()

		cr <- a.Unpack(src, reader)
	}()

	werr := <-cw
	rerr := <-cr

	if werr != nil {
		return werr
	}

	return rerr
}

func rebuildCache(src string, dest string, s storage.Storage, a archive.Archive) error {
	reader, writer := io.Pipe()

	cw := make(chan error)

	go func() {
		defer writer.Close()

		cw <- a.Pack(src, writer)
	}()

	cr := make(chan error)

	go func() {
		defer reader.Close()

		cr <- s.Put(dest, reader)
	}()

	werr := <-cw
	rerr := <-cr

	if werr != nil {
		return werr
	}

	return rerr
}

func buildPath(owner string, repo string, branch string) string {
	return fmt.Sprintf(
		"%s-%s/%s/archive.tar",
		owner,
		repo,
		branch,
	)
}
