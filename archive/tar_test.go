package archive

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"testing"
	. "github.com/franela/goblin"
)

func TestTarArchive(t *testing.T) {
	g := Goblin(t)

	g.Describe("NewTarArchive", func() {
		g.It("Should return tarArchive", func() {
			ta := NewTarArchive(&TarArchiveOptions{})
			g.Assert(ta != nil).IsTrue("failed to create tarArchive")
		})
	})

	g.Describe("Pack", func() {
		g.It("Should return no error", func() {
			ta := NewTarArchive(&TarArchiveOptions{})
			g.Assert(ta != nil).IsTrue("failed to create tarArchive")

			err, werr := packIt(ta, validMount)

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
			ta := NewTarArchive(&TarArchiveOptions{})
			g.Assert(ta != nil).IsTrue("failed to create tarArchive")

			err, werr := packIt(ta, invalidMount)

			g.Assert(err == nil).IsTrue("Failed to read the stream")
			g.Assert(werr != nil).IsTrue("Failed to properly stat 'mount'")
			g.Assert(werr.Error()).Equal("stat mount1: no such file or directory")
		})
	})

	g.Describe("Unpack", func() {
		g.It("Should return no error", func() {
			ta := NewTarArchive(&TarArchiveOptions{DryRun: true})
			g.Assert(ta != nil).IsTrue("failed to create tarArchive")

			err := unpackIt(ta, validFile)

			if err != nil {
				fmt.Printf("Received unexpected err: %s\n", err)
			}
			g.Assert(err == nil).IsTrue("Failed to unpack")
		})

		g.It("Should return error on invalid tarfile", func() {
			ta := NewTarArchive(&TarArchiveOptions{DryRun: true})
			g.Assert(ta != nil).IsTrue("failed to create tarArchive")

			err := unpackIt(ta, invalidFile)

			g.Assert(err != nil).IsTrue("Failed to return error")
			g.Assert(err.Error()).Equal("unexpected EOF")
		})

		g.It("Should return error on missing file", func() {
			ta := NewTarArchive(&TarArchiveOptions{DryRun: true})
			g.Assert(ta != nil).IsTrue("failed to create tarArchive")

			err := unpackIt(ta, missingFile)

			g.Assert(err != nil).IsTrue("Failed to return error")
			g.Assert(err.Error()).Equal("open fixtures/test2.tar: no such file or directory")
		})
	})
}

func packIt(a Archive, srcs []string) (error, error) {
	reader, writer := io.Pipe()
	defer reader.Close()

	cw := make(chan error, 1)
	defer close(cw)

	go func() {
		defer writer.Close()

		err := a.Pack(srcs, writer)

		if err != nil {
			cw <- err
			return
		} else {
			cw <- nil
		}
	}()

	_, err := ioutil.ReadAll(reader)

	werr := <-cw

	return err, werr
}

func unpackIt(a Archive, src string) error {
	reader, writer := io.Pipe()

	cw := make(chan error, 1)
	defer close(cw)

	f, err := os.Open(src)

	if err != nil {
		fmt.Printf("Failed to open file %s\n", src)
		return err
	}

	go func() {
		defer writer.Close()

		fmt.Printf("Copying %s\n", src)
		_, err = io.Copy(writer, f)
		fmt.Printf("Finished copying %s\n", src)

		if err != nil {
			fmt.Printf("Copy error'd for %s\n", src)
			cw <- err
			return
		}
	}()

	fmt.Printf("Unpacking %s\n", src)
	return a.Unpack("/dev/null", reader)
}

var (
	invalidMount = []string{
		"mount1",
		"mount2",
	}

	validMount = []string{
		"fixtures/test.txt",
		"fixtures/subdir",
	}

	validFile = "fixtures/test.tar"
	invalidFile = "fixtures/bad.tar"
	missingFile = "fixtures/test2.tar"
)
