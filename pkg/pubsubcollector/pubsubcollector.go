package pubsubcollector

import (
	"context"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

type Collector struct {
	ctx             context.Context
	pubSub          *pubsub.PubSub
	topic           *pubsub.Topic
}

