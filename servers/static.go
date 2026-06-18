package servers

import (
	"log"
	"net/http"
	"strings"
)

type StaticServerOptions struct {
	Address      string // host:port
	Directory    string // directory to serve
	Prefix       string // URL prefix to serve from
	SpaFallback  bool   // fallback to index.html for SPA routing
	EnableCors   bool   // enable CORS
	VerboseLevel int    // verbose logging level 0, 1 or 2. 0=silent, 1=endpoints, 2=headers and body
	Logger       *log.Logger
}

func (o StaticServerOptions) orDefault() StaticServerOptions {
	if o.Address == "" {
		o.Address = ":8080"
	}
	if o.Directory == "" {
		o.Directory = "."
	}
	if o.Prefix == "" {
		o.Prefix = "/"
	}
	if o.Logger == nil {
		o.Logger = defaultLogger()
	}

	return o
}

func NewStaticServer(opts StaticServerOptions) *http.Server {
	opts = opts.orDefault()

	handler := http.Handler(http.FileServer(http.Dir(opts.Directory)))
	if opts.SpaFallback {
		handler = withSpaFallback(handler, opts.Directory)
	}
	if opts.Prefix != "/" {
		handler = http.StripPrefix(normalizePrefix(opts.Prefix), handler)
	}
	if opts.EnableCors {
		handler = withCors(handler)
	}
	if opts.VerboseLevel > 0 {
		handler = withVerbose(handler, opts.Logger, opts.VerboseLevel)
	}

	mux := http.NewServeMux()
	mux.Handle(opts.Prefix, handler)

	return &http.Server{
		Addr:    opts.Address,
		Handler: mux,
	}
}

func normalizePrefix(prefix string) string {
	prefix = strings.Trim(strings.TrimSpace(prefix), "/")
	if prefix == "" {
		return "/"
	}
	return "/" + prefix + "/"
}
