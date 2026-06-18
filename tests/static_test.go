package main

import (
	"log"
	"testing"

	"github.com/renatopp/go-serve/serve"
	"github.com/renatopp/x/httpx"
	"github.com/renatopp/x/testx"
)

func TestStaticModeCLI(t *testing.T) {

	t.Run("serves files from directory", func(t *testing.T) {
		server := serve.NewStaticServer(serve.StaticServerOptions{
			Address:     ":0",
			Directory:   "./files",
			Logger:      log.New(nil, "", 0),
			Prefix:      "//sample/",
			SpaFallback: false,
			EnableCors:  false,
		})

		baseURL, closeServer := serveOnRandomPort(t, server)
		defer closeServer()

		res := httpx.Fetch("GET", baseURL+"/sample/index.html")
		testx.Equal(t, 200, res.StatusCode)
		testx.Equal(t, "index mock\n", res.Text())

		res2 := httpx.Fetch("GET", baseURL+"/sample/sample.md")
		testx.Equal(t, 200, res2.StatusCode)
		testx.Equal(t, "sample md\n", res2.Text())

		res3 := httpx.Fetch("GET", baseURL+"/sample/asdf")
		testx.Equal(t, 404, res3.StatusCode)
	})

	t.Run("supports spa fallback", func(t *testing.T) {
		server := serve.NewStaticServer(serve.StaticServerOptions{
			Address:     ":0",
			Directory:   "./files",
			Logger:      log.New(nil, "", 0),
			Prefix:      "//sample/",
			SpaFallback: true,
			EnableCors:  false,
		})

		baseURL, closeServer := serveOnRandomPort(t, server)
		defer closeServer()

		res := httpx.Fetch("GET", baseURL+"/sample/asdf")
		testx.Equal(t, 200, res.StatusCode)
		testx.Equal(t, "index mock\n", res.Text())

		res2 := httpx.Fetch("GET", baseURL+"/sample/sample.md")
		testx.Equal(t, 200, res2.StatusCode)
		testx.Equal(t, "sample md\n", res2.Text())
	})

	t.Run("supports optional cors headers", func(t *testing.T) {
		server := serve.NewStaticServer(serve.StaticServerOptions{
			Address:     ":0",
			Directory:   "./files",
			Logger:      log.New(nil, "", 0),
			Prefix:      "//sample/",
			SpaFallback: false,
			EnableCors:  true,
		})

		baseURL, closeServer := serveOnRandomPort(t, server)
		defer closeServer()

		res := httpx.Fetch("GET", baseURL+"/sample/")
		testx.Equal(t, 200, res.StatusCode)
		testx.Equal(t, "index mock\n", res.Text())
		testx.Equal(t, "*", res.Header("Access-Control-Allow-Origin"))
	})
}
