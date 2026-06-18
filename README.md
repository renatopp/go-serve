# go-serve

Web servers for testing and convenience.

`go-serve` is a CLI tool and library that provides convenient web servers functionality, such as serving static files, proxying requests, and more.

## Installation

To install `serve` CLI tool, just run:

```bash
go install github.com/renatopp/go-serve/cmd/serve@latest
```

To install as a library, use:

```bash
go get github.com/renatopp/go-serve/serve
```

Then, import the package `github.com/renatopp/go-serve` in your Go code.

## CLI Usage

Currently, there are the following server types available:

- [Static](#static)
- [Mock](#mock)
- [Echo](#echo)
- [Proxy](#proxy)

### Static

- Serve static files from current directory:

  `serve static`

- Serve static files from a hello/ directory:

  `serve static hello/`

- Serve static files at port 3030:

  `serve static . :3030`

- Serve static files for SPA with CORS enabled:

  `serve static --spa --cors`

- Serve static files with request logging:

  `serve static -vv`

### Mock

- Serve a mock server:

  `serve mock`

- Serve a mock serve at port 3030:

  `serve mock :3030`

- Serve a mock server with custom response:

  `serve mock --status=201 --header='Content-Type: text/plain' --body="Hello, World!"`

- Serve a mock server with delayed response:

  `serve mock --delay=5s`

- Serve a mock server with request logging:

  `serve mock -vv`

### Echo

- Serve an echo server:

  `serve echo`

- Serve an echo server at port 3030:

  `serve echo :3030`

- Serve an echo server with request logging:

  `serve echo -vv`

### Proxy

- Serve a proxy server forwarding to https://example.com:

  `serve proxy https://example.com`

- Serve a proxy server at port 3030:

  `serve proxy https://example.com :3030`

- Serve a proxy server with request logging:

  `serve proxy https://example.com -vv`

## Library Usage

To use `go-serve` as a library, import `github.com/renatopp/go-serve/serve` in your Go code.

```go
import "github.com/renatopp/go-serve/serve"

func main() {
  server := serve.NewStaticServer(serve.StaticServerOptions{
    Directory: "."
  })
  server.ListenAndServe()
}
```
