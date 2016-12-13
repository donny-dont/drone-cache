package cache

import (
	"fmt"
	"testing"
	. "github.com/franela/goblin"

	"github.com/drone-plugins/drone-cache/storage/dummy"
)

func TestCache(t *testing.T) {
	g := Goblin(t)

	g.Describe("New", func() {
		g.It("Should create new Cache", func() {
			s, err := dummy.New(dummyOpts)
			g.Assert(err == nil).IsTrue("failed to create storage")

			_, err = New(s)
			g.Assert(err == nil).IsTrue("failed to create cache")
		})
	})

	g.Describe("Rebuild", func() {
		g.It("Should rebuild with no errors", func() {
			s, err := dummy.New(dummyOpts)
			g.Assert(err == nil).IsTrue("failed to create storage")

			c, err := New(s)
			g.Assert(err == nil).IsTrue("failed to create cache")

			err = c.Rebuild([]string{"fixtures/test.txt", "fixtures/subdir"}, "file.tar")
			if err != nil {
				fmt.Printf("Received unexpected error: %s\n", err)
			}
			g.Assert(err == nil).IsTrue("failed to rebuild the cache")
		})

		g.It("Should return error on failure", func() {
			s, err := dummy.New(dummyOpts)
			g.Assert(err == nil).IsTrue("failed to create storage")

			c, err := New(s)
			g.Assert(err == nil).IsTrue("failed to create cache")

			err = c.Rebuild([]string{"mount1", "mount2"}, "file.ttt")
			g.Assert(err != nil).IsTrue("failed to return error")
			g.Assert(err.Error()).Equal("Unknown file format for archive file.ttt")
		})

		g.It("Should return error from channel", func() {
			s, err := dummy.New(dummyOpts)
			g.Assert(err == nil).IsTrue("failed to create storage")

			c, err := New(s)
			g.Assert(err == nil).IsTrue("failed to create cache")

			err = c.Rebuild([]string{"mount1", "mount2"}, "file.tar")
			g.Assert(err != nil).IsTrue("failed to return error")
			g.Assert(err.Error()).Equal("stat mount1: no such file or directory")
		})
	})

	g.Describe("Restore", func() {
		g.It("Should restore with no errors", func() {
			s, err := dummy.New(dummyOpts)
			g.Assert(err == nil).IsTrue("failed to create storage")

			c, err := New(s)
			g.Assert(err == nil).IsTrue("failed to create cache")

			err = c.Restore("fixtures/test.tar")
			if err != nil {
				fmt.Printf("Received unexpected error: %s\n", err)
			}
			g.Assert(err == nil).IsTrue("failed to rebuild the cache")
		})

		g.It("Should not return error on missing file", func() {
			s, err := dummy.New(dummyOpts)
			g.Assert(err == nil).IsTrue("failed to create storage")

			c, err := New(s)
			g.Assert(err == nil).IsTrue("failed to create cache")

			err = c.Restore("fixtures/test2.tar")
			g.Assert(err == nil).IsTrue("should not have returned error on missing file")
		})

		g.It("Should return error on unknown file format", func() {
			s, err := dummy.New(dummyOpts)
			g.Assert(err == nil).IsTrue("failed to create storage")

			c, err := New(s)
			g.Assert(err == nil).IsTrue("failed to create cache")

			err = c.Restore("fixtures/test2.ttt")
			g.Assert(err != nil).IsTrue("failed to return filetype error")
		})
	})
}

var (
	dummyOpts = &dummy.Options{
		Server:   "myserver.com",
		Username: "johndoe",
		Password: "supersecret",
	}
)
