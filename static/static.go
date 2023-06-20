package static

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed static
var webDir embed.FS

func StaticServerChroot() http.Handler {

	chrootedWebDir, err := fs.Sub(webDir, "web_dist")
	if err != nil {
		panic(err)
	}
	return http.FileServer(http.FS(chrootedWebDir))
}
func StaticServer() http.Handler {
	return http.FileServer(http.FS(webDir))
}

func Templates() fs.FS {
	chrootedWebDir, err := fs.Sub(webDir, "app")
	if err != nil {
		panic(err)
	}
	return chrootedWebDir
}
