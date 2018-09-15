package uid_test

import (
	"testing"

	"github.com/ribice/chisk/pkg/uid"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	id := uid.New()
	assert.Len(t, id, 20)
}
