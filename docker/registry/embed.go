package dockerregistry

import (
	"embed"
	"path"
)

//go:embed busybox/*
var images embed.FS

func read(s string) []byte {
	b, err := images.ReadFile(path.Join("busybox", s))
	check(err)
	return b
}
