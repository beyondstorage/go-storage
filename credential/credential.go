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
	// ProtocolBasic will hold user and password credential.
	//
	// value = [user, password]
	ProtocolBasic = "basic"
)

// Credential will provide credential protocol and values.
type Credential struct {
	protocol string
	args     []string
}

// Protocol provides current credential's protocol.
func (p Credential) Protocol() string {
	return p.protocol
}

// Value provides current credential's value in string array.
func (p Credential) Value() []string {
	return p.args
}

// Value provides current credential's value in string array.
func (p Credential) String() string {
	if len(p.args) == 0 {
		return p.protocol
	}
	return p.protocol + ":" + strings.Join(p.args, ":")
}

func (p Credential) Hmac() (accessKey, secretKey string) {
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

func (p Credential) APIKey() (apiKey string) {
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

func (p Credential) File() (path string) {
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

func (p Credential) Base64() (value string) {
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

func (p Credential) Basic() (user, password string) {
	if p.protocol != ProtocolBasic {
		panic(Error{
			Op:       "basic",
			Err:      ErrInvalidValue,
			Protocol: p.protocol,
			Values:   p.args,
		})
	}
	return p.args[0], p.args[1]
}

// Parse will parse config string to create a credential Credential.
func Parse(cfg string) (Credential, error) {
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
	case ProtocolBasic:
		return NewBasic(s[1], s[2]), nil
	default:
		return Credential{}, &Error{"parse", ErrUnsupportedProtocol, s[0], nil}
	}
}

// NewHmac create a hmac provider.
func NewHmac(accessKey, secretKey string) Credential {
	return Credential{ProtocolHmac, []string{accessKey, secretKey}}
}

// NewAPIKey create a api key provider.
func NewAPIKey(apiKey string) Credential {
	return Credential{ProtocolAPIKey, []string{apiKey}}
}

// NewFile create a file provider.
func NewFile(filePath string) Credential {
	return Credential{ProtocolFile, []string{filePath}}
}

// NewEnv create a env provider.
func NewEnv() Credential {
	return Credential{ProtocolEnv, nil}
}

// NewBase64 create a base64 provider.
func NewBase64(value string) Credential {
	return Credential{ProtocolBase64, []string{value}}
}

// NewBasic create a basic provider.
func NewBasic(user, password string) Credential {
	return Credential{ProtocolBasic, []string{user, password}}
}
