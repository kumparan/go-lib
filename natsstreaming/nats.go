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
		ID     int64  `json:"id"`
		UserID int64  `json:"user_id"`
		Body   string `json:"body,omitempty"`
		Time   string `json:"time"`
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

// QueueSubscribe :nodoc:
func (n *NATS) QueueSubscribe(subject, qgroup string, cb stan.MsgHandler, opts ...stan.SubscriptionOption) (stan.Subscription, error) {
	if n.testing {
		return nil, nil
	}

	return n.conn.QueueSubscribe(subject, qgroup, cb, opts...)
}
