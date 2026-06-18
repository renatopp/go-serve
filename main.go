package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/renatopp/go-cli"
	"github.com/renatopp/go-serve/servers"
	"github.com/renatopp/x/fmtx"
)

func main() {
	cli.Name("serve")
	cli.Description(`Quick web servers for testing and convenience.`)
	cli.AutoHelp(true)
	cli.Command("static", "Serve static files.", cmdStatic)
	cli.Command("mock", "Serve a mock server.", cmdMock)
	cli.Command("echo", "Serve an echo server.", cmdEcho)
	cli.Command("proxy", "Serve a reverse proxy.", cmdProxy)
	cli.Parse()
	cli.ShowHelp()
}

func cmdStatic() {
	directory := cli.Pos("directory", "Directory to serve").WithDefault(".")
	address := cli.Pos("address", "Address to listen on").WithDefault(":8080")
	prefix := cli.Flag("prefix", "p", "URL prefix to serve under").WithDefault("/")
	spa := cli.FlagBool("spa", "s", "Enable fallback for SPAs, serving index.html for 404").WithDefault(false)
	cors := cli.FlagBool("cors", "c", "Enable CORS").WithDefault(false)
	verbose := cli.FlagBool("", "v", "Enable request logging. -v for endpoints, -vv for headers and body.").AsRepeatable()

	cli.Parse()

	cli.Print("Serving static files from %s on %s\n", fmtx.Dim(directory.Value()), fmtx.Dim(address.Value()))
	server := servers.NewStaticServer(servers.StaticServerOptions{
		Directory:    directory.Value(),
		Address:      address.Value(),
		Prefix:       prefix.Value(),
		SpaFallback:  spa.Value(),
		EnableCors:   cors.Value(),
		VerboseLevel: len(verbose.Values()),
	})
	cli.FatalIf(server.ListenAndServe())
}

func cmdMock() {
	address := cli.Pos("address", "Address to listen on").WithDefault(":8080")
	status := cli.FlagInt("status", "s", "Response status code").WithDefault(http.StatusOK)
	body := cli.Flag("body", "b", "Response body").WithDefault("")
	delay := cli.FlagInt("delay", "d", "Response delay in milliseconds").WithDefault(0)
	headersRaw := cli.Flag("header", "H", "Response header. Repeatable: -H \"Key: Value\"").AsRepeatable()
	verbose := cli.FlagBool("", "v", "Enable request logging. -v for endpoints, -vv for headers and body.").AsRepeatable()
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

	cli.Print("Serving mock on %s\n", fmtx.Dim(address.Value()))
	server := servers.NewMockServer(servers.MockServerOptions{
		Address:      address.Value(),
		Status:       status.Value(),
		Body:         body.Value(),
		Headers:      headers,
		Delay:        time.Duration(delay.Value()) * time.Millisecond,
		VerboseLevel: len(verbose.Values()),
	})
	cli.FatalIf(server.ListenAndServe())
}

func cmdEcho() {
	address := cli.Pos("address", "Address to listen on").WithDefault(":8080")
	verbose := cli.FlagBool("", "v", "Enable request logging. -v for endpoints, -vv for headers and body.").AsRepeatable()
	cli.Parse()

	cli.Print("Serving echo on %s\n", fmtx.Dim(address.Value()))
	server := servers.NewEchoServer(servers.EchoServerOptions{
		Address:      address.Value(),
		VerboseLevel: len(verbose.Values()),
	})
	cli.FatalIf(server.ListenAndServe())
}

func cmdProxy() {
	targetAddress := cli.Pos("target_address", "Target address to proxy to")
	address := cli.Pos("self_address", "Address to listen on")
	verbose := cli.FlagBool("", "v", "Enable request logging. -v for endpoints, -vv for headers and body.").AsRepeatable()
	cli.Parse()

	cli.Print("Proxying requests from %s to %s\n", fmtx.Dim(address.Value()), fmtx.Dim(targetAddress.Value()))
	server, err := servers.NewProxyServer(servers.ProxyServerOptions{
		TargetAddress: targetAddress.Value(),
		Address:       address.Value(),
		VerboseLevel:  len(verbose.Values()),
	})
	cli.FatalIf(err)
	cli.FatalIf(server.ListenAndServe())
}
