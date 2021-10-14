# credential

Both human and machine-readable credential format.

## Format

```
<protocol>:<value>+
```

For example:

- hmac: `hmac:access_key:secret_key`
- apikey: `apikey:apikey`
- file: `file:/path/to/config/file`
- basic: `basic:user:password`

## Quick Start

```go
cred, err := credential.Parse("hmac:access_key:secret_key)
if err != nil {
    log.Fatal("parse: ", err)
}

switch cred.Protocol() {
case ProtocolHmac:
    ak, sk := cred.Hmac()
    log.Println("access_key: ", ak)
    log.Println("secret_key: ", sk)
case ProtocolAPIKey:
    apikey := cred.APIKey()
    log.Println("apikey: ", apikey)
case ProtocolFile:
    path := cred.File()
    log.Println("path: ", path)
case ProtocolEnv:
    log.Println("use env value")
case ProtocolBase64:
    content := cred.Base64()
    log.Println("base64: ", content)
case ProtocolBasic:
    user, password := cred.Basic()
    log.Println("user: ", user)
    log.Println("password: ", password)
default:
    panic("unsupported protocol")
}
```
