package serve

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/renatopp/x/fmtx"
	"github.com/renatopp/x/strx"
)

func NewDefaultLogger() *log.Logger {
	logger := log.New(os.Stdout, "", 0)
	return logger
}

func withVerbose(next http.Handler, logger *log.Logger, level int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch level {
		case 1:
			logger.Printf("%s %s", fmtx.Bold(r.Method), r.URL.String())
		case 2:
			logger.Printf("%s %s", fmtx.Bold(r.Method), r.URL.String())
			for name, values := range r.Header {
				logger.Printf("%s: %s", fmtx.Dim(name), strx.Join(values, ", "))
			}
			logger.Printf("\n")
			raw, err := io.ReadAll(r.Body)
			if err == nil && len(raw) > 0 {
				logger.Printf("%s", strx.Ellipsis(string(raw), 2048))
				logger.Printf("\n")
			}
			r.Body = io.NopCloser(bytes.NewBuffer(raw))
		}

		next.ServeHTTP(w, r)
	})
}

func withCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET,HEAD,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func withSpaFallback(next http.Handler, directory string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			next.ServeHTTP(w, r)
			return
		}

		cleanPath := path.Clean("/" + strings.TrimPrefix(r.URL.Path, "/"))
		if cleanPath == "/" || strings.Contains(path.Base(cleanPath), ".") {
			next.ServeHTTP(w, r)
			return
		}

		if _, err := os.Stat(filepath.Join(directory, cleanPath)); err == nil {
			next.ServeHTTP(w, r)
			return
		}

		clone := r.Clone(r.Context())
		clone.URL.Path = "/index.html"
		next.ServeHTTP(w, clone)
	})
}
