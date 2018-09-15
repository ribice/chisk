package session

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ribice/chisk/model"

	"github.com/go-redis/redis"
)

// New creates new redis session store
func New(c *redis.Client, d int) *Service {
	return &Service{client: c, duration: time.Duration(d) * time.Hour}
}

// Service represents in memory session store
type Service struct {
	client   *redis.Client
	duration time.Duration
}

type item struct {
	User *chisk.AuthUser
}

func (i *item) encode() string {
	b, _ := json.Marshal(i)
	return string(b)
}

func (*Service) decodeItem(str string) (*item, error) {
	var i item
	err := json.Unmarshal([]byte(str), &i)
	return &i, err
}

// Get tries to fetch existing session for given jwt token
func (s *Service) Get(token string) (*chisk.AuthUser, error) {
	val, err := s.client.Get(token).Result()
	if err != nil {
		return nil, err
	}

	i, err := s.decodeItem(val)
	if err != nil {
		return nil, err
	}

	return i.User, nil
}

// Put saves new user session based on key
func (s *Service) Put(user *chisk.User) error {
	i := item{
		User: user.AuthUser(),
	}

	return s.client.Set(user.Token, i.encode(), s.duration).Err()
}

// Delete deletes session based on jwt token key
func (s *Service) Delete(token string) error {
	return s.client.Del(token).Err()
}

// Update updates current user's session
func (s *Service) Update(user *chisk.User) error {
	// if User never logged in
	if user.Token == "" {
		return fmt.Errorf("missing token")
	}

	i := item{
		User: user.AuthUser(),
	}

	return s.client.GetSet(user.Token, i.encode()).Err()
}
