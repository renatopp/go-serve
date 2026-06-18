package servers

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type ProxyServerOptions struct {
	TargetAddress string
	Address       string
	VerboseLevel  int // verbose logging level 0, 1 or 2. 0=silent, 1=endpoints, 2=headers and body
	Logger        *log.Logger
}

func (o ProxyServerOptions) orDefault() ProxyServerOptions {
	if o.Address == "" {
		o.Address = ":8080"
	}
	if o.Logger == nil {
		o.Logger = defaultLogger()
	}
	return o
}

func NewProxyServer(opts ProxyServerOptions) (*http.Server, error) {
	opts = opts.orDefault()

	targetURL, err := parseProxyTargetAddress(opts.TargetAddress)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	proxy.ErrorLog = opts.Logger

	var handler http.Handler
	handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	})

	if opts.VerboseLevel > 0 {
		handler = withVerbose(handler, opts.Logger, opts.VerboseLevel)
	}

	return &http.Server{
		Addr:    opts.Address,
		Handler: handler,
	}, nil
}

func parseProxyTargetAddress(targetAddress string) (*url.URL, error) {
	targetAddress = strings.TrimSpace(targetAddress)
	if targetAddress == "" {
		return nil, fmt.Errorf("target address is required")
	}

	if !strings.Contains(targetAddress, "://") {
		targetAddress = "http://" + targetAddress
	}

	targetURL, err := url.Parse(targetAddress)
	if err != nil {
		return nil, fmt.Errorf("invalid target address %q: %w", targetAddress, err)
	}

	if targetURL.Scheme == "" || targetURL.Host == "" {
		return nil, fmt.Errorf("invalid target address %q", targetAddress)
	}

	return targetURL, nil
}
