package cache

import (
	"fmt"
	"testing"
	. "github.com/franela/goblin"

	"github.com/drone-plugins/drone-cache/storage"
)

func TestCache(t *testing.T) {
	g := Goblin(t)

	g.Describe("NewCache", func() {
		g.It("Should create new Cache", func() {
			s, err := storage.NewDummyStorage(dummyOpts)
			g.Assert(err == nil).IsTrue("failed to create storage")

			_, err = NewCache(s)
			g.Assert(err == nil).IsTrue("failed to create cache")
		})
	})

	g.Describe("Rebuild", func() {
		g.It("Should rebuild with no errors", func() {
			s, err := storage.NewDummyStorage(dummyOpts)
			g.Assert(err == nil).IsTrue("failed to create storage")

			c, err := NewCache(s)
			g.Assert(err == nil).IsTrue("failed to create cache")

			err = c.Rebuild([]string{"fixtures/test.txt", "fixtures/subdir"}, "file.tar")
			if err != nil {
				fmt.Printf("Received unexpected error: %s\n", err)
			}
			g.Assert(err == nil).IsTrue("failed to rebuild the cache")
		})

		g.It("Should return error on failure", func() {
			s, err := storage.NewDummyStorage(dummyOpts)
			g.Assert(err == nil).IsTrue("failed to create storage")

			c, err := NewCache(s)
			g.Assert(err == nil).IsTrue("failed to create cache")

			err = c.Rebuild([]string{"mount1", "mount2"}, "file.ttt")
			g.Assert(err != nil).IsTrue("failed to return error")
			g.Assert(err.Error()).Equal("Unknown file format for archive file.ttt")
		})

		g.It("Should return error from channel", func() {
			s, err := storage.NewDummyStorage(dummyOpts)
			g.Assert(err == nil).IsTrue("failed to create storage")

			c, err := NewCache(s)
			g.Assert(err == nil).IsTrue("failed to create cache")

			err = c.Rebuild([]string{"mount1", "mount2"}, "file.tar")
			g.Assert(err != nil).IsTrue("failed to return error")
			g.Assert(err.Error()).Equal("stat mount1: no such file or directory")
		})
	})
}

var (
	dummyOpts = &storage.DummyOptions{
		Server:   "myserver.com",
		Username: "johndoe",
		Password: "supersecret",
	}
)
