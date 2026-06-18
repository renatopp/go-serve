package servers

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type MockServerOptions struct {
	Address      string
	Status       int
	Body         string
	Headers      http.Header
	Delay        time.Duration
	VerboseLevel int // verbose logging level 0, 1 or 2. 0=silent, 1=endpoints, 2=headers and body
	Logger       *log.Logger
}

func (o MockServerOptions) orDefault() MockServerOptions {
	if o.Address == "" {
		o.Address = ":8080"
	}
	if o.Status == 0 {
		o.Status = http.StatusOK
	}
	if o.Logger == nil {
		o.Logger = defaultLogger()
	}
	return o
}

func NewMockServer(opts MockServerOptions) *http.Server {
	opts = opts.orDefault()

	var handler http.Handler
	handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if opts.Delay > 0 {
			time.Sleep(opts.Delay)
		}

		for name, values := range opts.Headers {
			for _, value := range values {
				w.Header().Add(name, value)
			}
		}

		body := opts.Body
		if body != "" {
			if w.Header().Get("Content-Type") == "" {
				w.Header().Set("Content-Type", contentTypeFromBody(body))
			}
			w.Header().Set("Content-Length", strconv.Itoa(len(body)))
		}

		w.WriteHeader(opts.Status)
		if body != "" {
			_, _ = w.Write([]byte(body))
		}
	})

	if opts.VerboseLevel > 0 {
		handler = withVerbose(handler, opts.Logger, opts.VerboseLevel)
	}

	return &http.Server{
		Addr:    opts.Address,
		Handler: handler,
	}
}

func contentTypeFromBody(body string) string {
	body = strings.TrimSpace(body)
	if strings.HasPrefix(body, "{") || strings.HasPrefix(body, "[") {
		return "application/json"
	}
	return "text/plain; charset=utf-8"
}
