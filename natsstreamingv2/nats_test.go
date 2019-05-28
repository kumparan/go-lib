package natsstreamingv2

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/alicebob/miniredis"
	redigo "github.com/gomodule/redigo/redis"
	natsServer "github.com/nats-io/gnatsd/server"
	"github.com/nats-io/nats-streaming-server/server"
)

const (
	clusterName = "my_test_cluster"
	defaultURL  = "nats://localhost:4222"
	clientName  = "test-client"
)

func runServer(clusterName string, port int) *server.StanServer {
	opts := server.GetDefaultOptions()
	opts.ID = clusterName
	s, err := server.RunServerWithOpts(opts, &natsServer.Options{
		Host: "localhost",
		Port: port,
	})
	if err != nil {
		panic(err)
	}
	return s
}

func runAnotherServer(port int) *server.StanServer {
	return runServer(clusterName, port)
}

func captureOutput(f func()) string {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	f()
	log.SetOutput(os.Stderr)
	return buf.String()
}

func newRedisConn(url string) *redigo.Pool {
	return &redigo.Pool{
		MaxIdle:     100,
		MaxActive:   10000,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redigo.Conn, error) {
			c, err := redigo.Dial("tcp", url)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redigo.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func TestRunningWorkerAfterLostConnection(t *testing.T) {
	port := 24444
	server := runAnotherServer(port)
	defer server.Shutdown()

	m, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}
	defer m.Close()

	r := newRedisConn(m.Addr())

	var concurrent = 10
	var natsConn []*NATS
	for i := 0; i < concurrent; i++ {
		NewNATSWithCallback(
			clusterName,
			clientName+fmt.Sprintf("-%d", i),
			fmt.Sprintf("localhost:%d", port),
			func(conn *NATS) {
				natsConn = append(natsConn, conn)
			},
		)
	}

	for _, conn := range natsConn {
		conn.SetRedisConn(r)
		defer conn.Close()
	}

	server.Shutdown()

	type msg struct {
		Data string `json:"data"`
	}

	var wg sync.WaitGroup

	for _, conn := range natsConn {
		wg.Add(1)
		go func(natsConn *NATS, wg *sync.WaitGroup) {
			defer wg.Done()
			ms := &msg{
				Data: "test",
			}
			err = natsConn.Publish("test", ms)
			if err != nil {
				t.Fatal(err)
			}
		}(conn, &wg)
	}

	client := r.Get()
	defer client.Close()

	wg.Wait()
	b, err := redigo.Int(client.Do("llen", failedMessagesRedisKey))
	if err != nil {
		t.Fatal(err)
	}

	if b != 10 {
		t.Fatal("Error value must be 10")
	}

	var wg2 sync.WaitGroup
	for _, conn := range natsConn {
		wg2.Add(1)
		go func(natsConn *NATS, wg *sync.WaitGroup) {
			defer wg.Done()
			natsConn.publishFromRedis()
		}(conn, &wg2)
	}

	wg2.Wait()

	b, err = redigo.Int(client.Do("llen", failedMessagesRedisKey))
	if err != nil {
		t.Fatal(err)
	}

	if b != 10 {
		t.Fatal("Error value must be 10")
	}
}
