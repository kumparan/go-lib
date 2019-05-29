package natsstreamingv2

import (
	"encoding/json"
	"sync"
	"time"

	redigo "github.com/gomodule/redigo/redis"
	"github.com/jasonlvhit/gocron"
	log "github.com/sirupsen/logrus"

	"github.com/kumparan/go-lib/utils"
	stan "github.com/nats-io/go-nats-streaming"
)

type (
	// EventType :nodoc:
	EventType string

	// NatsCallback :nodoc:
	NatsCallback func(conn *NATS)

	// NATS :nodoc:
	NATS struct {
		conn      stan.Conn
		redisConn *redigo.Pool
		testing   bool

		stopCh      chan struct{}
		reconnectCh chan struct{}
		wg          *sync.WaitGroup

		info              *natsInfo
		reconnectInterval time.Duration
		isWorkerRun       bool
		workerLock        *sync.Mutex
	}

	// NatsMessage :nodoc:
	NatsMessage struct {
		ID     int64     `json:"id"`
		UserID int64     `json:"user_id"`
		Type   EventType `json:"type"`
		Body   string    `json:"body,omitempty"`
		Time   string    `json:"time"`
	}

	NatsMessageWithSubject struct {
		Subject string `json:"subject"`
		Message []byte `json:"message"`
	}

	// natsInfo contains informations that will be use to reconnecting to nats streaming
	natsInfo struct {
		clusterID string
		clientID  string
		url       string
		opts      []stan.Option
		callback  NatsCallback
	}
)

const (
	failedMessagesRedisKey   = "nats:failed-messages"
	deadMessageRedisKey      = "nats:dead-messages"
	defaultReconnectInterval = 500 * time.Millisecond
)

// NewNATSWithCallback IMPORTANT! Not to send any stan.NatsURL or stan.SetConnectionLostHandler as options
func NewNATSWithCallback(clusterID, clientID, url string, fn NatsCallback, options ...stan.Option) {
	nc := &NATS{
		reconnectCh:       make(chan struct{}, 1),
		stopCh:            make(chan struct{}),
		wg:                new(sync.WaitGroup),
		reconnectInterval: defaultReconnectInterval,
		isWorkerRun:       false,
		workerLock:        new(sync.Mutex),
	}

	options = append(options, stan.SetConnectionLostHandler(func(conn stan.Conn, reason error) {
		log.Errorf("Nats connection lost! Reason: %v", reason)
		nc.reconnectCh <- struct{}{}
	}))

	nc.info = &natsInfo{
		url:       url,
		clusterID: clusterID,
		clientID:  clientID,
		callback:  fn,
		opts:      options,
	}

	conn, err := connect(clusterID, clientID, url, options...)
	if err != nil {
		log.WithFields(log.Fields{
			"clusterID": clusterID,
			"clientID":  clientID,
		}).Fatal(err)
	}

	log.WithFields(log.Fields{
		"clusterID": clusterID,
		"clientID":  clientID,
	}).Info("Nats connection made...")

	nc.SetConn(conn)
	// Run callback function
	nc.runCallback()
}

// connect to nats streaming
func connect(clusterID, clientID, url string, options ...stan.Option) (stan.Conn, error) {
	options = append(options, stan.NatsURL(url))
	nc, err := stan.Connect(clusterID, clientID, options...)
	if err != nil {
		return nil, err
	}
	return nc, nil
}

// NewTestNATS :nodoc:
func NewTestNATS() *NATS {
	return &NATS{testing: true}
}

// SetConn :nodoc:
func (n *NATS) SetConn(conn stan.Conn) {
	n.conn = conn
}

// SetRedisConn :nodoc:
func (n *NATS) SetRedisConn(conn *redigo.Pool) {
	n.redisConn = conn
}

// SetReconnectInterval :nodoc:
func (n *NATS) SetReconnectInterval(d time.Duration) {
	n.reconnectInterval = d
}

// Close NatsConnection :nodoc:
func (n *NATS) Close() error {
	close(n.stopCh)
	n.wg.Wait()
	close(n.reconnectCh)

	if n.conn != nil {
		err := n.conn.Close()
		if err != nil {
			return err
		}
	}

	if n.redisConn != nil {
		err := n.redisConn.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

// Run :nodoc:
func (n *NATS) Run() {
	if n.redisConn != nil {
		s := gocron.NewScheduler()
		s.Every(1).Minutes().Do(n.publishFromRedis)

		n.wg.Add(1)
		go func() {
			defer n.wg.Done()
			c := s.Start()
			select {
			case <-n.stopCh:
				close(c)
			}
		}()
	}

	n.wg.Add(1)
	go n.reconnectWorker()
}

// Publish :nodoc:
func (n *NATS) Publish(subject string, v interface{}) error {
	if n.testing {
		return nil
	}

	if n.conn != nil {
		err := n.conn.Publish(subject, utils.ToByte(v))
		if err == nil {
			return nil
		}
	}

	// Push to redis if failed
	if n.redisConn != nil {
		client := n.redisConn.Get()
		defer client.Close()
		i, err := redigo.Int(client.Do("RPUSH", failedMessagesRedisKey, utils.ToByte(NatsMessageWithSubject{
			Subject: subject,
			Message: utils.ToByte(v),
		})))
		if err != nil {
			log.WithFields(log.Fields{
				"subject": subject,
				"msg":     utils.ToByte(v),
			}).Error(err)
		}

		if i%100 == 0 {
			log.Errorf("total message that handled by redis : %+v", i)
		}
	} else {
		log.Fatal("publish to nats failed and no redis conection available to handle it")
	}
	return nil
}

// QueueSubscribe :nodoc:
func (n *NATS) QueueSubscribe(subject, qgroup string, cb stan.MsgHandler, opts ...stan.SubscriptionOption) (stan.Subscription, error) {
	if n.testing {
		return nil, nil
	}

	return n.conn.QueueSubscribe(subject, qgroup, cb, opts...)
}

// Subscribe :nodoc:
func (n *NATS) Subscribe(subject string, cb stan.MsgHandler, opts ...stan.SubscriptionOption) (stan.Subscription, error) {
	if n.testing {
		return nil, nil
	}
	return n.conn.Subscribe(subject, cb, opts...)
}

func (n *NATS) runCallback() {
	if n.info != nil && n.info.callback != nil {
		n.info.callback(n)
	}
}

func (n *NATS) publishFromRedis() {
	log.Info("running publish from redis...")
	n.workerLock.Lock()
	if n.isWorkerRun {
		n.workerLock.Unlock()
		return
	}
	n.isWorkerRun = true
	n.workerLock.Unlock()

	defer func() {
		n.workerLock.Lock()
		n.isWorkerRun = false
		n.workerLock.Unlock()
	}()

	if n.conn == nil {
		log.Error("abort due to connection problem")
		return
	}

	client := n.redisConn.Get()
	defer client.Close()

	for {
		b, err := redigo.Bytes(client.Do("LPOP", failedMessagesRedisKey))
		if err != nil && err != redigo.ErrNil {
			log.Error(err)
			return
		}

		if len(b) == 0 {
			return
		}

		msg := new(NatsMessageWithSubject)
		err = json.Unmarshal(b, msg)
		if err == nil {
			if n.conn != nil {
				err = n.conn.Publish(msg.Subject, msg.Message)
				if err == nil {
					continue
				}
			}
		}

		if err != nil {
			log.WithFields(log.Fields{
				"from":  "publishFromRedis",
				"value": string(b),
			}).Error(err)
		}

		_, errRedis := client.Do("RPUSH", deadMessageRedisKey, b)
		if errRedis != nil {
			log.WithFields(log.Fields{
				"from":  "publishFromRedis",
				"value": string(b),
			}).Error(errRedis)
			return
		}

		if err == stan.ErrConnectionClosed {
			log.WithFields(log.Fields{
				"from": "publishFromRedis",
			}).Error("Error nats connection closed")
			return
		}
	}
}

func (n *NATS) reconnectWorker() {
	defer n.wg.Done()
	for {
		select {
		case <-n.stopCh:
			return
		case <-n.reconnectCh:
			conn, err := connect(n.info.clusterID, n.info.clientID, n.info.url, n.info.opts...)
			if err != nil {
				log.Error("failed to reconnect")
				time.Sleep(n.reconnectInterval)

				// send to channel safely
				select {
				case n.reconnectCh <- struct{}{}:
				}
				continue
			}
			n.SetConn(conn)
			n.runCallback()
		}
	}
}
