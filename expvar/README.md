# expvar

This is a simple server that exposes the default variables from the standard
library's [expvar](https://godoc.org/expvar), which expose invocation command
line, and memory statistics. It also registers a custom expvar.

Although not the main purpose of this example, it also demonstrates how to
perform a graceful shutdown of an HTTP server as added in Go 1.8. A
`http.Server` variable must be declared in order to later call `Shutdown` on
it, which is not accessible when using the default HTTP server.

## Notes
You can change the address on which the HTTP server listens by setting the
environment variable `ADDR`, in standard Go **address:port** format (e.g.,
"localhost:5777", ":6067").
