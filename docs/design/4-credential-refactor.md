---
author: Xuanwo <github@xuanwo.io>
status: finished
updated_at: 2019-12-26
updates: [3](./3-support-service-init-via-config-string.md)
---

# Proposal: Credential refactor

## Background

Credential is the most import part to connect a Storage service. 

Current implement use a `credential.Provider` interface to output value:

```go
type Provider interface {
    Value() Value
}

type Value struct {
    AccessKey string
    SecretKey string
}
```

We only implement `static` protocol for Provider.

This implement inspired from aws: different credence are different ways to retrieve access key and secret key. However, this idea not works for every service. For Google: different credence are different ways to retrieve oauth2 tokens.

As a unified storage layer, we need to find a way to eliminate those inconsistencies.

## Proposal

Refactor `credential.Provider` to following:

```go
type Provider struct {
	protocol string
	args     []string
}

func (p *Provider) Protocol() string {
	return p.protocol
}

func (p *Provider) Value() []string {
	return p.args
}
```

User need to check credential's protocol before use them, and credential detail will be stored in `p.args` which can be retrieved via `p.Value()`.

Service who support only one credential protocol could be:

```go
credProtocol, cred := opt.Credential.Protocol(), opt.Credential.Value()
if credProtocol != credential.ProtocolHmac {
    return nil, fmt.Errorf(errorMessage, s, credential.ErrUnsupportedProtocol)
}
// Init service related config via credential values.
cfg, err := config.New(cred[0], cred[1])
```

Service who support more than one protocol could be:

```go
cfg := aws.NewConfig()

credProtocol, cred := opt.Credential.Protocol(), opt.Credential.Value()
switch credProtocol {
case credential.ProtocolHmac:
    cfg = cfg.WithCredentials(credentials.NewStaticCredentials(cred[0], cred[1], ""))
case credential.ProtocolEnv:
    cfg = cfg.WithCredentials(credentials.NewEnvCredentials())
default:
    return nil, fmt.Errorf(errorMessage, s, credential.ErrUnsupportedProtocol)
}
```

To implement a new protocol, develop should do following things.

- Add const in the format `Protocol<Name>` like `ProtocolEnv`
- Add comments for `Protocol<Name>`, describe clearly that how to use the value
- Implement init function `New<Name>(value ...string) (*Provider, error)` and `MustNew<Name>(value ...string) (*Provider, error)`
- Add `Protocol<Name>` into `credential.Parse` switch case
- Add unit test cases

## Rationale

AWS use a provider to get access key and secret key from static, file and env.

```go
type Credentials struct {
	creds        Value
	forceRefresh bool

	m sync.RWMutex

	provider Provider
}

type Value struct {
	// AWS Access key ID
	AccessKeyID string

	// AWS Secret Access Key
	SecretAccessKey string

	// AWS Session Token
	SessionToken string

	// Provider used to get credentials
	ProviderName string
}

type Provider interface {
	// Retrieve returns nil if it successfully retrieved the value.
	// Error is returned if the value were not obtainable, or empty.
	Retrieve() (Value, error)

	// IsExpired returns if the credentials are no longer valid, and need
	// to be retrieved.
	IsExpired() bool
}
```

Google has the same provider design but for oauth2 token:

```go
type Credentials struct {
    ProjectID   string // may be empty
    TokenSource oauth2.TokenSource

    // JSON contains the raw bytes from a JSON credentials file.
    // This field may be nil if authentication is provided by the
    // environment and not with a credentials file, e.g. when code is
    // running on Google Cloud Platform.
    JSON []byte
}
```

## Compatibility

Changes only introduce into `credential` package, credential's config string will be affected, no other changes in public interface.

## Implementation

Most of the work would be done by the author of this proposal.

In order to avoid misunderstanding between protocol `access key / secret key with hmac` and `api key`(they are both static keys), we rename current `static` protocol to `hmac`. No so inaccurate, but more clear.
