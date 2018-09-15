package redis

import (
	"fmt"
	"strconv"

	"github.com/go-redis/redis"
)

// New instantiates new Redis client
func New(addr, pass string, port int) (*redis.Client, error) {
	opts := redis.Options{
		Addr: addr + ":" + strconv.Itoa(port),
	}
	if pass != "" {
		opts.Password = pass
	}

	client := redis.NewClient(&opts)

	_, err := client.Ping().Result()
	if err != nil {
		return nil, fmt.Errorf("Cannot connect to Redis Addr %v, Port %v Reason %v", addr, port, err)
	}

	return client, nil
}
