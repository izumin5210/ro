package ro

// Model is a base definition of types.Model
type Model struct {
}

// GetKeySuffix implements the types.Model interface.
// If you does not override this function or return empty string, it will be returned errors.
func (b *Model) GetKeySuffix() string {
	return ""
}

// GetScoreMap implements the types.Model interface.
// If you does not override this function, ro.Store does not store any scores.
func (b *Model) GetScoreMap() map[string]interface{} {
	return map[string]interface{}{}
}
