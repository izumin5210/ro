package rotesting

import (
	"fmt"
	"log"

	"github.com/gomodule/redigo/redis"
	dockertest "gopkg.in/ory-am/dockertest.v3"
)

// Pool contains dockertest and redis connection pool
type Pool struct {
	redisPool  *redis.Pool
	dockerPool *dockertest.Pool
	dockerRes  *dockertest.Resource
}

// MustCreate creates new pool object
func MustCreate() *Pool {
	p := &Pool{}

	var err error
	p.dockerPool, err = dockertest.NewPool("")
	if err != nil {
		log.Fatalf("could not connect to docker: %s", err)
	}

	p.dockerRes, err = p.dockerPool.Run("redis", "4.0.2-alpine", nil)
	if err != nil {
		log.Fatalf("could not start resource: %s", err)
	}

	p.redisPool = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.DialURL(fmt.Sprintf("redis://localhost:%s", p.dockerRes.GetPort("6379/tcp")))
		},
	}

	if err = p.dockerPool.Retry(func() error {
		conn := p.Get()
		defer conn.Close()
		_, err := conn.Do("PING")

		return err
	}); err != nil {
		log.Fatalf("could not connect to docker: %s", err)
	}

	return p
}

// Get gets a connection with redis
func (p *Pool) Get() redis.Conn {
	return p.redisPool.Get()
}

// Cleanup remove all data in redis
func (p *Pool) Cleanup() error {
	conn := p.Get()
	defer conn.Close()
	_, err := conn.Do("FLUSHALL")
	return err
}

// MustClose closes redis connection pool and dockertest pool
func (p *Pool) MustClose() {
	var errs []error
	if err := p.Cleanup(); err != nil {
		errs = append(errs, err)
	}
	if err := p.redisPool.Close(); err != nil {
		errs = append(errs, err)
	}
	if err := p.dockerPool.Purge(p.dockerRes); err != nil {
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		log.Fatalf("unexpected error: %v", errs[0])
	}
}
