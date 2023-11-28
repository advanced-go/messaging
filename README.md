# messaging

Messaging provides functionality for resources to communicate. Resources may be in process, or out of process, accessible via HTTP. 

## core
[Core][corepkg] provides types for a message, message content, and a message cache.

## exchange
[Exchange][exchangepkg] provides functionality for sending messages to registered resources.

## mux
[Mux][muxpkg] provides HTTP request multiplexing, treating the request path as a URN. 

[corepkg]: <https://pkg.go.dev/github.com/advanced-go/messaging/core>
[exchangepkg]: <https://pkg.go.dev/github.com/advanced-go/messaging/exchange>
[muxpkg]: <https://pkg.go.dev/github.com/advanced-go/messaging/mux>
