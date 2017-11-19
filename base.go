package ro

// Base is a base definition of types.Model
type Base struct {
}

// GetKeyPrefix implements the types.Model interface.
// If you does not override this function, it will be used a type name in default.
func (b *Base) GetKeyPrefix() string {
	return ""
}

// GetKeySuffix implements the types.Model interface.
// If you does not override this function or return empty string, it will be returned errors.
func (b *Base) GetKeySuffix() string {
	return ""
}
