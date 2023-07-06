package iothub

import (
	"embed"
	"io/fs"
	"net/http"
)

var (
	//go:embed "config.default.yaml"
	DefaultConfigYaml []byte

	//go:embed "api/swagger_ui/*"
	swagFS embed.FS

	//go:embed  "ui/dist/*"
	webFS embed.FS
)

func RouteSwagger() {
	d, _ := fs.Sub(swagFS, "api/swagger_ui")
	http.Handle("/docs/", http.StripPrefix("/docs/", http.FileServer(http.FS(d))))
}

func RouteWeb() {
	d, _ := fs.Sub(webFS, "ui/dist")
	http.Handle("/ui/", http.StripPrefix("/ui/", http.FileServer(http.FS(d))))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/ui", http.StatusFound)
	})
}
