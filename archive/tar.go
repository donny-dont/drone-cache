package archive

// special thanks to this medium article:
// https://medium.com/@skdomino/taring-untaring-files-in-go-6b07cf56bc07

import (
	"archive/tar"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	log "github.com/Sirupsen/logrus"
)

type TarArchiveOptions struct {
	DryRun bool
}

type tarArchive struct{
	opts *TarArchiveOptions
}

// NewTarArchive creates an Archive that uses the .tar file format.
func NewTarArchive(opts *TarArchiveOptions) Archive {
	return &tarArchive{
		opts: opts,
	}
}

func (a *tarArchive) Pack(srcs []string, w io.Writer) error {
	tw := tar.NewWriter(w)
	defer tw.Close()

	// Loop through each source
	var fwErr error
	for _, s := range srcs {
		// ensure the src actually exists before trying to tar it
		if _, err := os.Stat(s); err != nil {
			return err
		}

		// walk path
		fwErr = filepath.Walk(s, func(path string, fi os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			var link string
			if fi.Mode()&os.ModeSymlink == os.ModeSymlink {
				if link, err = os.Readlink(path); err != nil {
					return err
				}
				log.Infof("Symbolic link found at %s", link)
			}

			header, err := tar.FileInfoHeader(fi, fi.Name())
			if err != nil {
				return err
			}

			header.Name = strings.TrimPrefix(filepath.ToSlash(path), "/")

			if err = tw.WriteHeader(header); err != nil {
				return err
			}

			if !fi.Mode().IsRegular() {
				log.Debugf("Directory found at %s", path)
				return nil
			}

			log.Debugf("File found at %s", path)

			file, err := os.Open(path)
			if err != nil {
				return err
			}

			defer file.Close()
			_, err = io.Copy(tw, file)
			return err
		})

		if fwErr != nil {
			return fwErr
		}
	}

	return fwErr
}

func (a *tarArchive) Unpack(dst string, r io.Reader) error {
	tr := tar.NewReader(r)

	for {
		header, err := tr.Next()

		switch {

		// if no more files are found return
		case err == io.EOF:
			return nil

		// return any other error
		case err != nil:
			return err

		// if the header is nil, just skip it (not sure how this happens)
		case header == nil:
			continue
		}

		// the target location where the dir/file should be created
		target := filepath.Join(dst, header.Name)

		// the following switch could also be done using fi.Mode(), not sure if there
		// a benefit of using one vs. the other.
		// fi := header.FileInfo()

		if a.opts.DryRun {
			_, err = io.Copy(ioutil.Discard, tr)

			if err != nil {
				return err
			}

			return nil
		}

		// check the file type
		switch header.Typeflag {

		// if its a dir and it doesn't exist create it
		case tar.TypeDir:
			log.Debugf("Directory found at %s", target)
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}

		// if it's a file create it
		case tar.TypeReg:
			log.Debugf("File found at %s", target)
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			// copy over contents
			_, err = io.Copy(f, tr)

			// Explicitly close otherwise too many files remain open
			f.Close()

			if err != nil {
				return err
			}
		}
	}
}
