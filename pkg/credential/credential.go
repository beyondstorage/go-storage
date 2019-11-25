package credential

// Provider can provide all authenticate needed info.
type Provider interface {
	Value() Value
}

// Value is the credential value.
type Value struct {
	AccessKey string
	SecretKey string
}
