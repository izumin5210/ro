package rotesting

import (
	"fmt"
)

// Post is a test object
type Post struct {
	ID        uint64 `redis:"id"`
	Title     string `redis:"title"`
	Body      string `redis:"body"`
	UpdatedAt int64  `redis:"updated_at"`
}

// GetKeySuffix implements the types.Model interface
func (p *Post) GetKeySuffix() string {
	return fmt.Sprint(p.ID)
}

// GetScoreMap implements the types.Model interface
func (p *Post) GetScoreMap() map[string]interface{} {
	return map[string]interface{}{
		"id":     p.ID,
		"recent": p.UpdatedAt,
	}
}
