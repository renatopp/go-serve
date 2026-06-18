package main

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/renatopp/go-serve/serve"
	"github.com/renatopp/x/testx"
)

type proxiedRequest struct {
	Method string
	Path   string
	Query  string
	Host   string
	Header string
	Body   string
}

func TestProxyModeCLI(t *testing.T) {
	t.Run("proxies request and response", func(t *testing.T) {
		requests := make(chan proxiedRequest, 1)
		target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rawBody, _ := io.ReadAll(r.Body)
			requests <- proxiedRequest{
				Method: r.Method,
				Path:   r.URL.Path,
				Query:  r.URL.RawQuery,
				Host:   r.Host,
				Header: r.Header.Get("X-Client"),
				Body:   string(rawBody),
			}

			w.Header().Set("X-Target", "ok")
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte("proxied"))
		}))
		defer target.Close()

		server, err := serve.NewProxyServer(serve.ProxyServerOptions{
			TargetAddress: target.URL,
			Address:       ":0",
			Logger:        log.New(io.Discard, "", 0),
		})
		if err != nil {
			t.Fatalf("failed creating proxy: %v", err)
		}

		baseURL, closeServer := serveOnRandomPort(t, server)
		defer closeServer()

		req, err := http.NewRequest(http.MethodPost, baseURL+"/api/echo?x=1", strings.NewReader("hello"))
		if err != nil {
			t.Fatalf("failed creating request: %v", err)
		}
		req.Header.Set("X-Client", "yes")

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("failed sending request: %v", err)
		}
		defer res.Body.Close()

		body, _ := io.ReadAll(res.Body)
		testx.Equal(t, http.StatusCreated, res.StatusCode)
		testx.Equal(t, "ok", res.Header.Get("X-Target"))
		testx.Equal(t, "proxied", string(body))

		select {
		case got := <-requests:
			testx.Equal(t, http.MethodPost, got.Method)
			testx.Equal(t, "/api/echo", got.Path)
			testx.Equal(t, "x=1", got.Query)
			testx.Equal(t, "yes", got.Header)
			testx.Equal(t, "hello", got.Body)
		case <-time.After(2 * time.Second):
			t.Fatal("timed out waiting proxied request")
		}
	})

	t.Run("rejects invalid target", func(t *testing.T) {
		server, err := serve.NewProxyServer(serve.ProxyServerOptions{
			TargetAddress: "://invalid",
			Address:       ":0",
		})

		testx.Equal(t, true, err != nil)
		testx.Equal(t, true, server == nil)
	})
}
