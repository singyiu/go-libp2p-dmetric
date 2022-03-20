package main

// taking references from libp2p pubsub example

import (
	"context"
	"fmt"
	"github.com/reactivex/rxgo/v2"
	"github.com/rs/zerolog/log"
	"github.com/singyiu/go-libp2p-dmetric/pkg/dmetric"
	"github.com/singyiu/go-libp2p-dmetric/pkg/pubsubpublisher"
	"github.com/singyiu/go-libp2p-dmetric/pkg/rx"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
)

// DiscoveryInterval is how often we re-publish our mDNS records.
const DiscoveryInterval = time.Hour

// DiscoveryServiceTag is used in our mDNS advertisements to discover other chat peers.
const DiscoveryServiceTag = "go-libp2p-dmetric-example-dcounter"

// MetricPublisherInterval interval for checking if the registered metrics need to be published
const MetricPublisherInterval = time.Minute

// CounterIncInterval interval for increasing the test counter val
const CounterIncInterval = time.Second * 10

// discoveryNotifee gets notified when we find a new peer via mDNS discovery
type discoveryNotifee struct {
	h host.Host
}

// HandlePeerFound connects to peers discovered via mDNS. Once they're connected,
// the PubSub system will automatically start interacting with them if they also
// support PubSub.
func (n *discoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
	fmt.Printf("discovered new peer %s\n", pi.ID.Pretty())
	err := n.h.Connect(context.Background(), pi)
	if err != nil {
		fmt.Printf("error connecting to peer %s: %s\n", pi.ID.Pretty(), err)
	}
}

// setupDiscovery creates an mDNS discovery service and attaches it to the libp2p Host.
// This lets us automatically discover peers on the same LAN and connect to them.
func setupDiscovery(h host.Host) error {
	// setup mDNS discovery to find local peers
	s := mdns.NewMdnsService(h, DiscoveryServiceTag, &discoveryNotifee{h: h})
	return s.Start()
}

func GetMapFuncAnyToIncreasedCounterVal(c *dmetric.Counter) func(context.Context, interface{}) (interface{}, error) {
	return func(_ context.Context, i interface{}) (interface{}, error) {
		c.Inc()
		return c.GetVal(), nil
	}
}

// start the processing loop
func start(ctx context.Context) {
	// create a new libp2p Host that listens on a random TCP port
	h, err := libp2p.New(libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"))
	if err != nil {
		panic(err)
	}

	// create a new PubSub service using the GossipSub router
	ps, err := pubsub.NewGossipSub(ctx, h)
	if err != nil {
		panic(err)
	}

	// setup local mDNS discovery
	if err := setupDiscovery(h); err != nil {
		panic(err)
	}

	// create a pubsub publisher that would check and publish if any metric should be published at regular interval
	publisher, err := pubsubpublisher.NewIntervalPublisher(ctx, ps, "", MetricPublisherInterval)
	if err != nil {
		panic(err)
	}

	// create a test counter and register it with the publisher
	counter01 := dmetric.NewCounter("testCounter01", "testDesc01", 0, map[string]string{"label01":"labelValue01"})
	publisher.RegisterPublishableObj(counter01)

	// increase counter at regular interval
	log.Info().Msgf("Increase counter at regular interval %+v", CounterIncInterval)
	updateCounterCh := rxgo.Interval(rxgo.WithDuration(CounterIncInterval)).
		Map(GetMapFuncAnyToIncreasedCounterVal(counter01)).
		Map(rx.GetSideEffectLog("IncreasedCounterVal")).
		OnErrorReturn(rx.GetErrFuncLogError("updateCounterCh")).
		Observe(rxgo.WithErrorStrategy(rxgo.ContinueOnError))

	for {
		select {
		case _, ok := <-updateCounterCh:
			if !ok {
				log.Fatal().Stack().Msg("updateCounterCh closed")
			}
		case <-ctx.Done():
			return
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go start(ctx)
	<-rx.OsSignalTermProducer().Observe() // wait for termination signal
	cancel()
}
