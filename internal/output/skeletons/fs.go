package skeletons

import (
	"embed"
	"io/fs"
)

//go:embed templates/**
var RootFS embed.FS
var FS fs.FS

func init() {
	var err error

	FS, err = fs.Sub(RootFS, "templates")
	if err != nil {
		panic(err)
	}
}
