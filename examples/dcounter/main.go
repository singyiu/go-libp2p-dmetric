package main

import (
	"context"
	"github.com/reactivex/rxgo/v2"
	"github.com/rs/zerolog/log"
	"github.com/singyiu/go-libp2p-dmetric/pkg/dmetric"
	"github.com/singyiu/go-libp2p-dmetric/pkg/rxos"
	"time"
)

func GetSideEffectIncCounterAndLog(cv *dmetric.CounterVec) func(context.Context, interface{}) (interface{}, error) {
	return func(_ context.Context, i interface{}) (interface{}, error) {
		labelMap := map[string]string{"label01":"value01"}
		cv.Inc(labelMap)
		log.Info().Uint64(cv.Name, cv.GetValueOf(labelMap)).Msg("")
		return i, nil
	}
}

func start(ctx context.Context) {
	counterVec01 := dmetric.NewCounterVec("testCounterVec01", "testDesc01")

	// increase counter at regular interval
	updateCounterCh := rxgo.Interval(rxgo.WithDuration(time.Second * 10)).
		Map(GetSideEffectIncCounterAndLog(counterVec01)).
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
	<-rxos.OsSignalTermProducer().Observe() // wait for termination signal
	cancel()
}