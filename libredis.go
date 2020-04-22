package libredis

import (
	"context"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	"go.elastic.co/apm/module/apmredigo"
)

// Options is ..
type Options struct {
	Host        string
	Port        int
	Password    string
	MaxIdle     int
	MaxActive   int
	IdleTimeout int
	Enabled     bool
}

var pool *redis.Pool

// Connect is ...
func Connect(options Options) *redis.Pool {
	if pool == nil {
		pool = &redis.Pool{
			MaxIdle:     options.MaxIdle,
			MaxActive:   options.MaxActive,
			IdleTimeout: time.Duration(options.IdleTimeout) * time.Second,
			Dial: func() (redis.Conn, error) {
				address := fmt.Sprintf("%s:%d", options.Host, options.Port)
				c, err := redis.Dial("tcp", address)
				if err != nil {
					return nil, err
				}

				// Do authentication process if password not empty.
				if options.Password != "" {
					if _, err := c.Do("AUTH", options.Password); err != nil {
						c.Close()
						return nil, err
					}
				}

				return c, nil
			},
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				if time.Since(t) < time.Minute {
					return nil
				}

				_, err := c.Do("PING")
				return err
			},
			Wait:            true,
			MaxConnLifetime: 15 * time.Minute,
		}

		return pool
	}

	return pool
}

// ConnectWithAMP is ...
func ConnectWithAMP(options Options) *redis.Pool {
	if pool == nil {
		pool = &redis.Pool{
			MaxIdle:     options.MaxIdle,
			MaxActive:   options.MaxActive,
			IdleTimeout: time.Duration(options.IdleTimeout) * time.Second,
			Dial: func() (redis.Conn, error) {
				address := fmt.Sprintf("%s:%d", options.Host, options.Port)
				c, err := redis.Dial("tcp", address)
				apmredigo.Wrap(c).WithContext(context.Background())
				if err != nil {
					return nil, err
				}

				// Do authentication process if password not empty.
				if options.Password != "" {
					if _, err := c.Do("AUTH", options.Password); err != nil {
						c.Close()
						return nil, err
					}
				}

				return c, nil
			},
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				if time.Since(t) < time.Minute {
					return nil
				}

				_, err := c.Do("PING")
				return err
			},
			Wait:            true,
			MaxConnLifetime: 15 * time.Minute,
		}

		return pool
	}

	return pool
}
