package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/renatopp/go-cli"
	"github.com/renatopp/go-serve/serve"
	"github.com/renatopp/x/fmtx"
)

var logger = serve.NewDefaultLogger()

const cliDescription = `Web servers for testing and convenience.
Provides static file serving, mock servers, echo servers, and reverse proxy.
More information: http://github.com/renatopp/go-serve
`

func main() {
	cli.Name("serve")
	cli.Description(cliDescription)
	cli.AutoHelp(true)
	cli.Command("static", "Serve static files.", cmdStatic)
	cli.Command("mock", "Serve a mock server.", cmdMock)
	cli.Command("echo", "Serve an echo server.", cmdEcho)
	cli.Command("proxy", "Serve a reverse proxy.", cmdProxy)
	cli.Parse()
	cli.ShowHelp()
}

func cmdStatic() {
	cli.Name("static")
	cli.Description("Serve static files from a directory.")
	directory := cli.Pos("directory", "Directory to serve").WithDefault(".")
	address := cli.Pos("address", "Address to listen on").WithDefault(":8080")
	prefix := cli.Flag("prefix", "p", "URL prefix to serve under").WithDefault("/")
	spa := cli.FlagBool("spa", "s", "Enable fallback for SPAs, serving index.html for 404").WithDefault(false)
	cors := cli.FlagBool("cors", "c", "Enable CORS").WithDefault(false)
	verbose := cli.FlagBool("", "v", "-v logs endpoints, -vv logs endpoints, headers and body.").AsRepeatable()

	cli.Parse()

	logger.Printf("Serving static files from %s on %s", fmtx.Dim(directory.Value()), fmtx.Dim(address.Value()))
	server := serve.NewStaticServer(serve.StaticServerOptions{
		Directory:    directory.Value(),
		Address:      address.Value(),
		Prefix:       prefix.Value(),
		SpaFallback:  spa.Value(),
		EnableCors:   cors.Value(),
		VerboseLevel: len(verbose.Values()),
		Logger:       logger,
	})
	cli.FatalIf(server.ListenAndServe())
}

func cmdMock() {
	cli.Name("Mock")
	cli.Description("Serve a mock server which can return specified status code, headers and body.")
	address := cli.Pos("address", "Address to listen on").WithDefault(":8080")
	status := cli.FlagInt("status", "s", "Response status code").WithDefault(http.StatusOK)
	body := cli.Flag("body", "b", "Response body").WithDefault("")
	delay := cli.FlagDuration("delay", "d", "Response delay in milliseconds").WithDefault(0)
	headersRaw := cli.Flag("header", "H", "Response header. Repeatable: -H \"Key: Value\"").AsRepeatable()
	verbose := cli.FlagBool("", "v", "-v logs endpoints, -vv logs endpoints, headers and body.").AsRepeatable()
	cli.Parse()

	headers := make(http.Header)
	for _, raw := range headersRaw.Values() {
		parts := strings.SplitN(raw, ":", 2)
		if len(parts) != 2 {
			cli.FatalIf(fmt.Errorf("invalid header %q, expected \"Key: Value\"", raw))
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		if key == "" {
			cli.FatalIf(fmt.Errorf("invalid header %q, missing key", raw))
		}
		headers.Add(key, value)
	}

	logger.Printf("Serving mock on %s", fmtx.Dim(address.Value()))
	server := serve.NewMockServer(serve.MockServerOptions{
		Address:      address.Value(),
		Status:       status.Value(),
		Body:         body.Value(),
		Headers:      headers,
		Delay:        delay.Value(),
		VerboseLevel: len(verbose.Values()),
		Logger:       logger,
	})
	cli.FatalIf(server.ListenAndServe())
}

func cmdEcho() {
	cli.Name("Echo")
	cli.Description("Serve an echo server which responds the request data (endpoint, headers and body) as JSON.")
	address := cli.Pos("address", "Address to listen on").WithDefault(":8080")
	verbose := cli.FlagBool("", "v", "-v logs endpoints, -vv logs endpoints, headers and body.").AsRepeatable()
	cli.Parse()

	logger.Printf("Serving echo on %s", fmtx.Dim(address.Value()))
	server := serve.NewEchoServer(serve.EchoServerOptions{
		Address:      address.Value(),
		VerboseLevel: len(verbose.Values()),
		Logger:       logger,
	})
	cli.FatalIf(server.ListenAndServe())
}

func cmdProxy() {
	cli.Name("Proxy")
	cli.Description("Serve a reverse proxy which forwards requests to a target address.")
	target := cli.Pos("target_address", "Target address to proxy to").AsRequired()
	self := cli.Pos("self_address", "Address to listen on")
	verbose := cli.FlagBool("", "v", "-v logs endpoints, -vv logs endpoints, headers and body.").AsRepeatable()
	cli.Parse()

	logger.Printf("Proxying requests from %s to %s", fmtx.Dim(self.Value()), fmtx.Dim(target.Value()))
	server, err := serve.NewProxyServer(serve.ProxyServerOptions{
		TargetAddress: target.Value(),
		Address:       self.Value(),
		VerboseLevel:  len(verbose.Values()),
		Logger:        logger,
	})
	cli.FatalIf(err)
	cli.FatalIf(server.ListenAndServe())
}
