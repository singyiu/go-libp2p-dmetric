# go-libp2p-dmetric

## Prerequisite
* go 1.18+

## Examples
* [dCounter](examples/dcounter/README.md)

## How to use
### (A) To publish metric in a libp2p node
```
import (
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
    "github.com/singyiu/go-libp2p-dmetric/pkg/dmetric"
    "github.com/singyiu/go-libp2p-dmetric/pkg/pubsubpublisher"
)

// create a new libp2p Host that listens on a random TCP port
h, err := libp2p.New(libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"))

// create a new PubSub service using the GossipSub router
ps, err := pubsub.NewGossipSub(ctx, h)
topic, err := ps.Join(DiscoveryServiceTag)

// setup local mDNS discovery
err := setupDiscovery(h)

// create a pubsub publisher that would check and publish if any metric should be published at regular interval
publisher, err := pubsubpublisher.NewIntervalPublisher(ctx, topic, MetricPublisherInterval)

// create metric
counter01 := dmetric.NewCounter(h.ID(), "testCounter01", "testDesc01", 0, labelPairs01)

// register the metric with the publisher		
publisher.RegisterPublishableObj(counter01)
```

### (B) To collect metric from libp2p pubsub
```
import (
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
    "github.com/singyiu/go-libp2p-dmetric/pkg/dmetric"
    "github.com/singyiu/go-libp2p-dmetric/pkg/pubsubpublisher"
    "github.com/singyiu/go-libp2p-dmetric/pkg/rxpubsub"
)

// create a new libp2p Host that listens on a random TCP port
h, err := libp2p.New(libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"))

// create a new PubSub service using the GossipSub router
ps, err := pubsub.NewGossipSub(ctx, h)
topic, err := ps.Join(DiscoveryServiceTag)

// setup local mDNS discovery
err := setupDiscovery(h)

// create dMetric registerer
reg := dmetric.NewRegisterer()

// create a pubsub message producer that would subscribe to the target topic
messageProducer, err := rxpubsub.GetMessageProducerFromTopic(ctx, h.ID(), topic)

// setup pipeline to transform the pubsub message into a dMetric message
// and publish it to the Prometheus endpoint
dMessageCh := messageProducer.
			Map(rxpubsub.MapFuncPubSubMsgToObj[dmetric.Message]).
			Map(dmetric.GetSideEffectPublishMessageToPrometheus(reg)).
			OnErrorReturn(rx.GetErrFuncLogError("dMessageCh")).
			Observe(rxgo.WithErrorStrategy(rxgo.ContinueOnError))

// serve the /metric endpoint for Prometheus to collect the metrics
go prometheushelper.RunServer(ctx, PrometheusServerAddressStr)
```