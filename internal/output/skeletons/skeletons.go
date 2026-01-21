package skeletons

import (
	"embed"
	"io/fs"
)

//go:embed skeleton-templates/**
var RootFS embed.FS
var FS fs.FS

func init() {
	var err error

	FS, err = fs.Sub(RootFS, "skeleton-templates")
	if err != nil {
		panic(err)
	}
}
