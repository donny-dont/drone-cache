package tar

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"testing"
	. "github.com/franela/goblin"

	"github.com/drone-plugins/drone-cache/archive"
)

type mountFile struct {
	Path string
	Content string
}

func TestTarArchive(t *testing.T) {
	g := Goblin(t)
	wd, _ := os.Getwd()

	// Create necessary fixtures
	createFixtures()

	g.Describe("New", func() {
		g.It("Should return tarArchive", func() {
			ta := New()
			g.Assert(ta != nil).IsTrue("failed to create tarArchive")
		})
	})

	g.Describe("Pack", func() {
		g.It("Should return no error", func() {
			ta := New()
			g.Assert(ta != nil).IsTrue("failed to create tarArchive")

			os.Chdir("/tmp/fixtures/mounts")
			err, werr := packIt(ta, validMount, "/tmp/fixtures/tarfiles/test.tar")
			os.Chdir(wd)

			if err != nil {
				fmt.Printf("Received unexpected err: %s\n", err)
			}
			g.Assert(err == nil).IsTrue("Failed to read the stream")
			if werr != nil {
				fmt.Printf("Received unexpected werr: %s\n", werr)
			}
			g.Assert(werr == nil).IsTrue("Failed to pack")
		})

		g.It("Should return error if mount does not exist", func() {
			ta := New()
			g.Assert(ta != nil).IsTrue("failed to create tarArchive")

			err, werr := packIt(ta, invalidMount, "/tmp/fixtures/tarfiles/invalidMount.tar")

			g.Assert(err == nil).IsTrue("Failed to read the stream")
			g.Assert(werr != nil).IsTrue("Failed to properly stat 'mount'")
			g.Assert(werr.Error()).Equal("stat mount1: no such file or directory")
		})
	})

	g.Describe("Unpack", func() {
		g.It("Should return no error", func() {
			ta := New()
			g.Assert(ta != nil).IsTrue("failed to create tarArchive")

			err := unpackIt(ta, validFile)

			if err != nil {
				fmt.Printf("Received unexpected err: %s\n", err)
			}
			g.Assert(err == nil).IsTrue("Failed to unpack")
		})

		g.It("Should create files in correct strucutre", func() {
			g.Assert(exists("/tmp/extracted/test.txt")).IsTrue("failed to create test.txt")
			g.Assert(exists("/tmp/extracted/subdir")).IsTrue("failed to create subdir")
			g.Assert(exists("/tmp/extracted/subdir/test2.txt")).IsTrue("failed to create subdir/test2.txt")
		})

		g.It("Should create files with correct content", func() {
			var err error
			var content []byte
			for _, element := range mountFiles {
				content, err = ioutil.ReadFile("/tmp/extracted/" + element.Path)
				g.Assert(err == nil).IsTrue("failed to read" + element.Path)
				g.Assert(string(content)).Equal(element.Content)
			}
		})

		g.It("Should return error on invalid tarfile", func() {
			ta := New()
			g.Assert(ta != nil).IsTrue("failed to create tarArchive")

			err := unpackIt(ta, invalidFile)

			g.Assert(err != nil).IsTrue("Failed to return error")
			g.Assert(err.Error()).Equal("unexpected EOF")
		})

		g.It("Should return error on missing file", func() {
			ta := New()
			g.Assert(ta != nil).IsTrue("failed to create tarArchive")

			err := unpackIt(ta, missingFile)

			g.Assert(err != nil).IsTrue("Failed to return error")
			g.Assert(err.Error()).Equal("open /tmp/fixtures/tarfiles/test2.tar: no such file or directory")
		})
	})
}

func packIt(a archive.Archive, srcs []string, dst string) (error, error) {
	reader, writer := io.Pipe()
	defer reader.Close()

	cw := make(chan error, 1)
	defer close(cw)

	go func() {
		defer writer.Close()

		cw <- a.Pack(srcs, writer)
	}()

	bytes, err := ioutil.ReadAll(reader)
	ioutil.WriteFile(dst, bytes, 0644)

	werr := <-cw

	return err, werr
}

func unpackIt(a archive.Archive, src string) error {
	reader, writer := io.Pipe()

	cw := make(chan error, 1)
	defer close(cw)

	f, err := os.Open(src)

	if err != nil {
		return err
	}

	go func() {
		defer writer.Close()

		_, err = io.Copy(writer, f)

		if err != nil {
			cw <- err
			return
		}
	}()

	return a.Unpack("/tmp/extracted", reader)
}

func createBadTarfile() {
	content := []byte("hello\ngo\n")
	err := ioutil.WriteFile("/tmp/fixtures/tarfiles/bad.tar", content, 0644)
	if err != nil {
		log.Fatalln(err)
	}
}

func createMountContent() {
	var err error
	for _, element := range mountFiles {
		err = ioutil.WriteFile("/tmp/fixtures/mounts/" + element.Path, []byte(element.Content), 0644)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func createFixtures() {
	// Create directories if does not exist
	if _, err := os.Stat("/tmp/fixtures/tarfiles"); os.IsNotExist(err) {
		os.MkdirAll("/tmp/fixtures/tarfiles", os.FileMode(int(0755)))
	}
	if _, err := os.Stat("/tmp/fixtures/mounts/subdir"); os.IsNotExist(err) {
		os.MkdirAll("/tmp/fixtures/mounts/subdir", os.FileMode(int(0755)))
	}
	if _, err := os.Stat("/tmp/extracted"); os.IsNotExist(err) {
		os.MkdirAll("/tmp/extracted", os.FileMode(int(0755)))
	}

	createBadTarfile()
	createMountContent()
}

func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil { return true }
	if os.IsNotExist(err) { return false }
	return true
}

var (
	invalidMount = []string{
		"mount1",
		"mount2",
	}

	mountFiles = []mountFile{
		{Path: "test.txt", Content: "hello\ngo\n"},
		{Path: "subdir/test2.txt", Content: "hello2\ngo\n"},
	}

	validMount = []string{
		"test.txt",
		"subdir",
	}

	validFile = "/tmp/fixtures/tarfiles/test.tar"
	invalidFile = "/tmp/fixtures/tarfiles/bad.tar"
	missingFile = "/tmp/fixtures/tarfiles/test2.tar"
)