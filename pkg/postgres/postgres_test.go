package pgsql_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/ribice/chisk/model"

	"github.com/stretchr/testify/assert"

	"github.com/go-pg/pg"

	"github.com/ory/dockertest"
	"github.com/ribice/chisk/pkg/postgres"
)

func TestNew(t *testing.T) {
	assert := assert.New(t)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	resource, err := pool.Run("postgres", "9.6.5-alpine", []string{"POSTGRES_PASSWORD=secret", "POSTGRES_DB=test"})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	var db *pg.DB
	hostPort := resource.GetHostPort("5432/tcp")
	dbConnStr := fmt.Sprintf("postgres://postgres:secret@%s/test?sslmode=disable", hostPort)

	_, err = pgsql.New("invalidPSN", false, 0)
	assert.Error(err)

	fmt.Println(dbConnStr)

	if err = pool.Retry(func() error {
		db, err = pgsql.New(dbConnStr, true, 1)
		return err
	}); err != nil {
		t.Fatalf("Could not connect to docker: %s", err)
	}

	assert.NotNil(db)

	err = db.Select(&chisk.AuthUser{})
	assert.Error(err)
}
