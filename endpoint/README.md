# endpoint

Both human and machine readable endpoint format.

## Format

```
<protocol>:<value>+
```

For example:

- File: `file:/var/cache/data`
- HTTP: `http:example.com:80`
- HTTPS: `https:example.com:443`

## Quick Start

```go
ep, err := endpoint.Parse("https:example.com")
if err != nil {
	log.Fatal("parse: ", err)
}

switch ep.Protocol() {
case ProtocolHTTP:
    url, host, port := ep.HTTP()
    log.Println("url: ", url)
    log.Println("host: ", host)
    log.Println("port: ", port)
case ProtocolHTTPS:
    url, host, port := ep.HTTPS()
    log.Println("url: ", url)
    log.Println("host: ", host)
    log.Println("port: ", port)
case ProtocolFile:
    path := ep.File()
    log.Println("path: ", path)
default:
    panic("unsupported protocol")
}
```
