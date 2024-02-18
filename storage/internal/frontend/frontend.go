package frontend

import (
	"go.uber.org/zap"
	"net/http"
	"os"
	"path"
	"strings"
)

const FSPATH = "./build/"

func ServeFrontend() {
	if _, err := os.Stat(FSPATH); os.IsNotExist(err) {
		zap.S().Warn("frontend files not found")
		return
	}
	fs := http.FileServer(http.Dir(FSPATH))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// If the requested file exists then return if; otherwise return index.html (fileserver default page)
		if r.URL.Path != "/" {
			fullPath := FSPATH + strings.TrimPrefix(path.Clean(r.URL.Path), "/")
			_, err := os.Stat(fullPath)
			if err != nil {
				if !os.IsNotExist(err) {
					panic(err)
				}
				// Requested file does not exist so we return the default (resolves to index.html)
				r.URL.Path = "/"
			}
		}
		fs.ServeHTTP(w, r)
	})
	http.ListenAndServe(":3000", nil)
}
