package mock

import (
	"github.com/ribice/chisk/model"
)

// Session mock
type Session struct {
	GetFn func(string) (*chisk.AuthUser, error)
}

// Get mock
func (s *Session) Get(token string) (*chisk.AuthUser, error) {
	return s.GetFn(token)
}
