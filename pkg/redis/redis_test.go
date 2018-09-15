package redis_test

import (
	"log"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"gopkg.in/ory-am/dockertest.v3"

	"github.com/ribice/chisk/pkg/redis"
)

func TestNew(t *testing.T) {
	assert := assert.New(t)
	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("Could not connect to docker: %s", err)
	}

	resource, err := pool.Run("redis", "4.0.11", nil)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	port, _ := strconv.Atoi(resource.GetPort("6379/tcp"))

	rc, err := redis.New("localhost", "pass", 6379)
	assert.Error(err)
	assert.Nil(rc)

	rc, err = redis.New("localhost", "", port)
	assert.NoError(err)
	assert.NotNil(rc)

	pool.Purge(resource)

}
