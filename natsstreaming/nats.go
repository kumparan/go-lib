package natsstreaming

import (
	"github.com/kumparan/go-lib/utils"
	"github.com/nats-io/go-nats-streaming"
)

type (
	NATS struct {
		conn    stan.Conn
		testing bool
	}

	// NatsMessage :nodoc:
	NatsMessage struct {
		ID     int64
		Type   string
		UserID int64
		Body   string
		Time   string
	}
)

// NewNATS :nodoc:
func NewNATS(clusterID, clientID, url string) (*NATS, error) {
	nc, err := stan.Connect(clusterID, clientID, stan.NatsURL(url))
	if err != nil {
		return nil, err
	}

	return &NATS{conn: nc}, nil
}

// NewTestNATS :nodoc:
func NewTestNATS() *NATS {
	return &NATS{testing: true}
}

// Close NatsConnection :nodoc:
func (n *NATS) Close() {
	n.conn.Close()
}

// Publish :nodoc:
func (n *NATS) Publish(subject string, v interface{}) error {
	if n.testing {
		return nil
	}

	return n.conn.Publish(subject, utils.ToByte(v))
}
