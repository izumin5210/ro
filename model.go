package ro

// Model is a base definition of types.Model
type Model struct {
}

// GetKeyPrefix implements the types.Model interface.
// If you does not override this function, it will be used a type name in default.
func (b *Model) GetKeyPrefix() string {
	return ""
}

// GetKeySuffix implements the types.Model interface.
// If you does not override this function or return empty string, it will be returned errors.
func (b *Model) GetKeySuffix() string {
	return ""
}
