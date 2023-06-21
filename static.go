package main

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed htdocs
var webDir embed.FS

func StaticServerChroot() http.Handler {

	chrootedWebDir, err := fs.Sub(webDir, "htdocs")
	if err != nil {
		panic(err)
	}
	return http.FileServer(http.FS(chrootedWebDir))
}
func StaticServer() http.Handler {
	return http.FileServer(http.FS(webDir))
}

func Templates() fs.FS {
	chrootedWebDir, err := fs.Sub(webDir, "htdocs")
	if err != nil {
		panic(err)
	}
	return chrootedWebDir
}
