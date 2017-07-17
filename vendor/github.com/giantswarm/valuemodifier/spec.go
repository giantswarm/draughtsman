package valuemodifier

// ValueModifier implements some modification mechanism for values being
// provided. This can e.g. be an implementation to encrypt a given value using
// GPG encryption standards or encode a given value using base64 encoding
// standard.
type ValueModifier interface {
	Modify(value []byte) ([]byte, error)
}
