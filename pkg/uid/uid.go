package uid

import (
	"github.com/rs/xid"
)

// New returns new unique ID
func New() string {
	return xid.New().String()
}
