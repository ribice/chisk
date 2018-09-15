package session_test

import (
	"encoding/json"
	"log"
	"strconv"
	"testing"
	"time"

	"github.com/ribice/chisk/pkg/session"

	"github.com/ribice/chisk/model"
	"github.com/ribice/chisk/pkg/redis"
	"github.com/stretchr/testify/assert"
	dockertest "gopkg.in/ory-am/dockertest.v3"
)

type item struct {
	User *chisk.AuthUser
}

func (i item) encode() string {
	b, _ := json.Marshal(i)
	return string(b)
}

func decodeItem(str string) (*item, error) {
	var i item
	err := json.Unmarshal([]byte(str), &i)
	return &i, err
}

func TestGet(t *testing.T) {
	assert := assert.New(t)

	cases := []struct {
		name     string
		getToken string
		setToken string
		setData  interface{}
		wantErr  bool
		wantData *chisk.AuthUser
	}{
		{
			name:     "Different set and get token",
			setToken: "token1",
			getToken: "token2",
			wantErr:  true,
		},
		{
			name:     "Invalid item",
			setToken: "token1",
			getToken: "token1",
			setData:  "simpleString",
			wantErr:  true,
		},
		{
			name:     "Success",
			setToken: "token1",
			getToken: "token1",
			setData: item{
				User: &chisk.AuthUser{
					ID:          "userid",
					Email:       "johndoe@mail.com",
					DisplayName: "johndoe",
				},
			}.encode(),
			wantData: &chisk.AuthUser{
				ID:          "userid",
				Email:       "johndoe@mail.com",
				DisplayName: "johndoe",
			},
		},
	}

	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("Could not connect to docker: %s", err)
	}

	resource, err := pool.Run("redis", "4.0.11", nil)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	port, _ := strconv.Atoi(resource.GetPort("6379/tcp"))

	rclient, err := redis.New("localhost", "", port)
	if err != nil {
		t.Fatal(err)
	}
	sessSvc := session.New(rclient, 1)

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {

			err := rclient.Set(tt.setToken, tt.setData, 1*time.Hour).Err()
			assert.NoError(err)

			au, err := sessSvc.Get(tt.getToken)
			assert.Equal(tt.wantErr, err != nil)
			assert.Equal(tt.wantData, au)

		})
	}

	pool.Purge(resource)

}

func TestPut(t *testing.T) {
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

	rclient, err := redis.New("localhost", "", port)
	if err != nil {
		t.Fatal(err)
	}

	sessSvc := session.New(rclient, 1)

	err = sessSvc.Put(&chisk.User{
		Base: chisk.Base{
			ID: "userid",
		},
		Email:       "johndoe@mail.com",
		DisplayName: "johndoe",
		Token:       "usertoken",
	})
	assert.NoError(err)

	strUser, err := rclient.Get("usertoken").Result()
	assert.NoError(err)

	item, err := decodeItem(strUser)
	assert.NoError(err)

	assert.Equal(&chisk.AuthUser{
		ID:          "userid",
		Email:       "johndoe@mail.com",
		DisplayName: "johndoe",
	}, item.User)

	pool.Purge(resource)

}

func TestDelete(t *testing.T) {
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

	rclient, err := redis.New("localhost", "", port)
	if err != nil {
		t.Fatal(err)
	}

	sessSvc := session.New(rclient, 1)

	err = rclient.Set("usertoken", "value", 1*time.Minute).Err()
	assert.NoError(err)

	err = sessSvc.Delete("usertoken")
	assert.NoError(err)

	err = rclient.Get("usertoken").Err()
	assert.Error(err)

	pool.Purge(resource)
}

func TestUpdate(t *testing.T) {
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

	rclient, err := redis.New("localhost", "", port)
	if err != nil {
		t.Fatal(err)
	}

	sessSvc := session.New(rclient, 1)

	user := &chisk.User{
		DisplayName: "johndoe",
	}

	err = sessSvc.Update(user)
	assert.Error(err)

	user.Token = "usertoken"

	err = sessSvc.Update(user)
	assert.Error(err)

	i := item{
		User: user.AuthUser(),
	}

	err = rclient.Set("usertoken", i.encode(), 1*time.Minute).Err()
	assert.NoError(err)

	err = sessSvc.Update(user)
	assert.NoError(err)

	strUser, err := rclient.Get("usertoken").Result()
	assert.NoError(err)

	item, err := decodeItem(strUser)
	assert.NoError(err)

	assert.Equal(&chisk.AuthUser{
		DisplayName: "johndoe",
	}, item.User)

	pool.Purge(resource)
}
