package templates

import (
	"embed"
	"io/fs"
)

// FS contains the embedded default templates for installation and generation.
//
//go:embed config.toml commands.allow env gitignore.block instructions/*.md slash-commands/*.md docs/agent-layer/*.md
var FS embed.FS

// Read returns the contents of an embedded template file.
func Read(path string) ([]byte, error) {
	return fs.ReadFile(FS, path)
}

// Walk walks the embedded template filesystem under the given root.
func Walk(root string, fn fs.WalkDirFunc) error {
	return fs.WalkDir(FS, root, fn)
}
