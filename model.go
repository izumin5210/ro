package ro

// Model is a base definition of types.Model
type Model struct {
}

// GetKeySuffix implements the types.Model interface.
// If you does not override this function or return empty string, it will be returned errors.
func (b *Model) GetKeySuffix() string {
	return ""
}
