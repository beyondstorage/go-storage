package credential

import (
	"strings"
)

const (
	// ProtocolHmac will hold access key and secret key credential.
	//
	// HMAC means hash-based message authentication code, it may be inaccurate to represent credential
	// protocol ak/sk(access key + secret key with hmac), but it's simple and no confuse with other
	// protocol, so just keep this.
	//
	// value = [Access Key, Secret Key]
	ProtocolHmac = "hmac"
	// ProtocolAPIKey will hold api key credential.
	//
	// value = [API Key]
	ProtocolAPIKey = "apikey"
	// ProtocolFile will hold file credential.
	//
	// value = [File Path], service decide how to use this file
	ProtocolFile = "file"
	// ProtocolEnv will represent credential from env.
	//
	// value = [], service retrieves credential value from env.
	ProtocolEnv = "env"
	// ProtocolBase64 will represents credential binary data in base64
	//
	// Storage service like gcs will take token files as input, we provide base64 protocol so that user
	// can pass token binary data directly.
	ProtocolBase64 = "base64"
)

// Provider will provide credential protocol and values.
type Provider struct {
	protocol string
	args     []string
}

// Protocol provides current credential's protocol.
func (p Provider) Protocol() string {
	return p.protocol
}

// Value provides current credential's value in string array.
func (p Provider) Value() []string {
	return p.args
}

// Value provides current credential's value in string array.
func (p Provider) String() string {
	if len(p.args) == 0 {
		return p.protocol
	}
	return p.protocol + ":" + strings.Join(p.args, ":")
}

func (p Provider) Hmac() (accessKey, secretKey string) {
	if p.protocol != ProtocolHmac {
		panic(Error{
			Op:       "hmac",
			Err:      ErrInvalidValue,
			Protocol: p.protocol,
			Values:   p.args,
		})
	}
	return p.args[0], p.args[1]
}

func (p Provider) APIKey() (apiKey string) {
	if p.protocol != ProtocolAPIKey {
		panic(Error{
			Op:       "api_key",
			Err:      ErrInvalidValue,
			Protocol: p.protocol,
			Values:   p.args,
		})
	}
	return p.args[0]
}

func (p Provider) File() (path string) {
	if p.protocol != ProtocolFile {
		panic(Error{
			Op:       "file",
			Err:      ErrInvalidValue,
			Protocol: p.protocol,
			Values:   p.args,
		})
	}
	return p.args[0]
}

func (p Provider) Base64() (value string) {
	if p.protocol != ProtocolBase64 {
		panic(Error{
			Op:       "base64",
			Err:      ErrInvalidValue,
			Protocol: p.protocol,
			Values:   p.args,
		})
	}
	return p.args[0]
}

// Parse will parse config string to create a credential Provider.
func Parse(cfg string) (Provider, error) {
	s := strings.Split(cfg, ":")

	switch s[0] {
	case ProtocolHmac:
		return NewHmac(s[1], s[2]), nil
	case ProtocolAPIKey:
		return NewAPIKey(s[1]), nil
	case ProtocolFile:
		return NewFile(s[1]), nil
	case ProtocolEnv:
		return NewEnv(), nil
	case ProtocolBase64:
		return NewBase64(s[1]), nil
	default:
		return Provider{}, &Error{"parse", ErrUnsupportedProtocol, s[0], nil}
	}
}

// NewHmac create a hmac provider.
func NewHmac(accessKey, secretKey string) Provider {
	return Provider{ProtocolHmac, []string{accessKey, secretKey}}
}

// NewAPIKey create a api key provider.
func NewAPIKey(apiKey string) Provider {
	return Provider{ProtocolAPIKey, []string{apiKey}}
}

// NewFile create a file provider.
func NewFile(filePath string) Provider {
	return Provider{ProtocolFile, []string{filePath}}
}

// NewEnv create a env provider.
func NewEnv() Provider {
	return Provider{ProtocolEnv, nil}
}

// NewBase64 create a base64 provider.
func NewBase64(value string) Provider {
	return Provider{ProtocolBase64, []string{value}}
}
