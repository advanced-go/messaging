# messaging

Messaging inspired to capitalize on the Go language for application development. Determining the patterns that need to be employed is critical for writing clear idiomatic Go code. This YouTube video 

## core
[Core][corepkg] provides the types for a message, message content, and a message cache.

## exchange
[Exchange][exchangepkg] provides functionality for sending messages to registered resources.

## mux
[Mux][muxpkg] provides HTTP request multiplexing, processing the request path as a URN. 

[corepkg]: <https://pkg.go.dev/github.com/advanced-go/messaging/core>
[exchangepkg]: <https://pkg.go.dev/github.com/advanced-go/messaging/exchange>
[muxpkg]: <https://pkg.go.dev/github.com/advanced-go/messaging/mux>
