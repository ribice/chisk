package chisk

import (
	"time"

	"github.com/go-pg/pg/orm"
	"github.com/rs/xid"
)

// Base contains common table fields
type Base struct {
	ID        string     `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// BeforeInsert hooks into insert operations
func (b *Base) BeforeInsert(_ orm.DB) error {
	b.ID = xid.New().String()
	now := time.Now()
	b.CreatedAt = now
	b.UpdatedAt = now
	return nil
}

// BeforeUpdate hooks into update operations
func (b *Base) BeforeUpdate(_ orm.DB) error {
	b.UpdatedAt = time.Now()
	return nil
}
