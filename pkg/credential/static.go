package credential

const (
	// ProtocolStatic is the static credential protocol.
	ProtocolStatic = "static"
)

// Static will hold static credential.
type Static struct {
	accessKey string
	secretKey string
}

// NewStatic will create a new static credential.
func NewStatic(accessKey, secretKey string) Static {
	return Static{
		accessKey: accessKey,
		secretKey: secretKey,
	}
}

// Value implements Provider interface.
func (s Static) Value() Value {
	return Value{
		AccessKey: s.accessKey,
		SecretKey: s.secretKey,
	}
}
