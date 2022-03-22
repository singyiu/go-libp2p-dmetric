package pubsubpublisher

import (
	"context"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/reactivex/rxgo/v2"
	"github.com/singyiu/go-libp2p-dmetric/pkg/common"
	"github.com/singyiu/go-libp2p-dmetric/pkg/rx"
	"time"
)

type Publisher struct {
	ctx             context.Context
	topic           *pubsub.Topic
	publishableObjs []common.Publishable
}

func NewIntervalPublisher(ctx context.Context, topic *pubsub.Topic, publishInterval time.Duration) (*Publisher, error) {
	p := Publisher{
		ctx:   ctx,
		topic: topic,
	}
	go p.StartPublishIntervalLoop(ctx, publishInterval)
	return &p, nil
}

func (p *Publisher) RegisterPublishableObj(obj common.Publishable) {
	p.publishableObjs = append(p.publishableObjs, obj)
}

func GetMapFuncAnyToPublishResults(p *Publisher) func(context.Context, interface{}) (interface{}, error) {
	return func(ctx context.Context, i interface{}) (interface{}, error) {
		var results []string
		for _, obj := range p.publishableObjs {
			if obj.ShouldBePublished() {
				bytes, err := obj.ToJsonBytes()
				if err != nil {
					return results, err
				}
				err = p.topic.Publish(ctx, bytes)
				if err != nil {
					return results, err
				}
				obj.OnPublished()
				results = append(results, string(bytes))
			}
		}
		return results, nil
	}
}

func (p *Publisher) StartPublishIntervalLoop(ctx context.Context, publishInterval time.Duration) {
	ch := rxgo.Interval(rxgo.WithDuration(publishInterval)).
		Map(GetMapFuncAnyToPublishResults(p)).
		Map(rx.GetSideEffectLog("PublishResults")).
		OnErrorReturn(rx.GetErrFuncLogError("StartPublishIntervalLoop")).
		Observe(rxgo.WithErrorStrategy(rxgo.ContinueOnError))

	for {
		select {
		case _, ok := <-ch:
			if !ok {
				panic("StartPublishIntervalLoop ch closed")
			}
		case <-ctx.Done():
			return
		}
	}
}
