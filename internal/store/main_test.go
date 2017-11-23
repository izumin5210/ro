package store

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/garyburd/redigo/redis"
	dockertest "gopkg.in/ory-am/dockertest.v3"
)

type TestPost struct {
	ID        uint64 `redis:"id"`
	Title     string `redis:"title"`
	Body      string `redis:"body"`
	UpdatedAt int64  `redis:"updated_at"`
}

func (p *TestPost) GetKeySuffix() string {
	return fmt.Sprint(p.ID)
}

var redisPool *redis.Pool

func TestMain(m *testing.M) {
	var err error
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	resource, err := pool.Run("redis", "4.0.2-alpine", nil)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	redisPool = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.DialURL(fmt.Sprintf("redis://localhost:%s", resource.GetPort("6379/tcp")))
		},
	}

	if err = pool.Retry(func() error {
		conn := redisPool.Get()
		defer conn.Close()
		_, err := conn.Do("PING")

		return err
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	exitCode := m.Run()

	err = redisPool.Close()
	if err != nil {
		log.Fatalf("Failed to close redis pool: %s", err)
	}
	err = pool.Purge(resource)
	if err != nil {
		log.Fatalf("Failed to purge docker pool: %s", err)
	}

	os.Exit(exitCode)
}

func teardown(t *testing.T) {
	conn := redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("FLUSHALL")
	if err != nil {
		log.Fatalf("Failed to flush redis: %s", err)
	}
}
