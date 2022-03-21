package rxpubsub

import (
	"context"
	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/reactivex/rxgo/v2"
	"github.com/rs/zerolog/log"
)

func LoopSubscribeToTopicAndPublishToChannel(ctx context.Context, hostId peer.ID, sub *pubsub.Subscription, outputCh chan rxgo.Item) {
	for {
		msg, err := sub.Next(ctx)
		if err != nil {
			log.Error().Err(err).Msg("LoopSubscribeToTopicAndPublishToChannel sun.Next failed")
			close(outputCh)
			return
		}
		// only forward messages delivered by others
		if msg.ReceivedFrom == hostId {
			continue
		}
		outputCh <- rxgo.Of(msg)
	}
}

func GetMessageProducerFromTopic(ctx context.Context, hostId peer.ID, topic *pubsub.Topic) (rxgo.Observable, error) {
	outputCh := make(chan rxgo.Item)
	sub, err := topic.Subscribe()
	if err != nil {
		return nil, err
	}
	go LoopSubscribeToTopicAndPublishToChannel(ctx, hostId, sub, outputCh)
	return rxgo.FromChannel(outputCh), nil
}