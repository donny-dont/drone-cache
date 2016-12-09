package plugin

import (
	"os"
	"testing"
	. "github.com/franela/goblin"
)

func TestCli(t *testing.T) {
	g := Goblin(t)

	g.Describe("PluginFlags", func() {

		g.It("Should return flags given the ENV", func() {
			_ = PluginFlags()

			// g.Assert(flags.String("filename")).Equal("arch")
		})
	})
}

func testEnv() {
	os.Setenv("PLUGIN_FILENAME", "archive.tar")
	os.Setenv("PLUGIN_PATH", "plugin/path")
	os.Setenv("PLUGIN_MOUNT", "mount1,mount2")
	os.Setenv("PLUGIN_REBUILD", "true")
	os.Setenv("PLUGIN_RESTORE", "true")
	os.Setenv("PLUGIN_DEBUG", "true")
	os.Setenv("DRONE_REPO_OWNER", "johndoe")
	os.Setenv("DRONE_REPO_NAME", "octocat/hello-world")
	os.Setenv("DRONE_COMMIT_BRANCH", "develop")
}
