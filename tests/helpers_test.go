package main

import (
	"fmt"
	"net"
	"net/http"
	"testing"
)

func serveOnRandomPort(t *testing.T, server *http.Server) (string, func()) {
	t.Helper()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to open listener: %v", err)
	}

	go server.Serve(listener)

	return fmt.Sprintf("http://%s", listener.Addr().String()), func() {
		_ = server.Close()
	}
}
