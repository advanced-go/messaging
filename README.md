# messaging

Messaging inspired to capitalize on the Go language for application development. Determining the patterns that need to be employed is critical for writing clear idiomatic Go code. This YouTube video [Edward Muller - Go Anti-Patterns][emuller], does an excellent job of framing idiomatic go. 
[Robert Griesemer - The Evolution of Go][rgriesemer], @ 4:00 minutes, a


## core
[Core][corepkg] provides functionality for processing an Http request/response. Exchange functionality is provied via a templated function, utilizing
template paramters for error processing, deserialization type, and the function for processing the http.Client.Do():

## exchange
[Exchange][exchangepkg] provides functionality for processing an Http request/response. Exchange functionality is provied via a templated function, utilizing
template paramters for error processing, deserialization type, and the function for processing the http.Client.Do():

## mux
[Mux][muxpkg] provides functionality for processing an Http request/response. Exchange functionality is provied via a templated function, utilizing
template paramters for error processing, deserialization type, and the function for processing the http.Client.Do():

[corepkg]: <https://pkg.go.dev/github.com/advanced-go/messaging/core>
[exchangepkg]: <https://pkg.go.dev/github.com/advanced-go/messaging/exchange>
[muxpkg]: <https://pkg.go.dev/github.com/advanced-go/messaging/mux>
