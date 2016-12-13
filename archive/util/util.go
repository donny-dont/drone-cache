package util

import (
	"fmt"
	"strings"

	. "github.com/drone-plugins/drone-cache/archive"
	"github.com/drone-plugins/drone-cache/archive/tar"
)

// FromFilename determines the archive format to use based on the name.
func FromFilename(name string) (Archive, error) {
	if strings.HasSuffix(name, ".tar") {
		return tar.New(&tar.Options{}), nil
	}

	return nil, fmt.Errorf("Unknown file format for archive %s", name)
}
