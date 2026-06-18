package servers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type EchoServerOptions struct {
	Address      string // host:port
	VerboseLevel int    // verbose logging level 0, 1 or 2. 0=silent, 1=endpoints, 2=headers and body
	Logger       *log.Logger
}

func (e EchoServerOptions) orDefault() EchoServerOptions {
	if e.Address == "" {
		e.Address = ":8080"
	}
	if e.Logger == nil {
		e.Logger = defaultLogger()
	}
	return e
}

type EchoResponse struct {
	Method  string              `json:"method"`
	Path    string              `json:"path"`
	Query   map[string][]string `json:"query"`
	Headers map[string][]string `json:"headers"`
	Body    string              `json:"body"`
}

func NewEchoServer(opts EchoServerOptions) *http.Server {
	opts = opts.orDefault()
	var handler http.Handler
	handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawBody, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "failed to read request body", http.StatusBadRequest)
			return
		}

		response := EchoResponse{
			Method:  r.Method,
			Path:    r.URL.Path,
			Query:   r.URL.Query(),
			Headers: map[string][]string(r.Header),
			Body:    string(rawBody),
		}

		data, err := json.Marshal(response)
		if err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(data)
	})

	if opts.VerboseLevel > 0 {
		handler = withVerbose(handler, opts.Logger, opts.VerboseLevel)
	}

	return &http.Server{
		Addr:    opts.Address,
		Handler: handler,
	}
}
