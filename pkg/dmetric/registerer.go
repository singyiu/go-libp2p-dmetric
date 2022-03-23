package dmetric

import (
	"context"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"sync/atomic"
)

type Registerer struct {
	CounterMsgMap map[string]*Message
}

func NewRegisterer() *Registerer {
	r := Registerer{
		CounterMsgMap: make(map[string]*Message),
	}
	return &r
}

func GetSideEffectPublishMessageToPrometheus(reg *Registerer) func(context.Context, interface{}) (interface{}, error) {
	return func(_ context.Context, i interface{}) (interface{}, error) {
		msg, ok := i.(*Message)
		if !ok {
			return nil, fmt.Errorf("input not *Message %+v", i)
		}

		labelPairStr := msg.GetLabelPairsStr()
		switch msg.Type {
		case MetricTypeCounter:
			counterMsg, ok := reg.CounterMsgMap[labelPairStr]
			if ok {
				//update counter value
				atomic.StoreUint64(&counterMsg.UIntVal, msg.UIntVal)
			} else {
				//register new counter
				reg.CounterMsgMap[labelPairStr] = msg
				promauto.NewCounterFunc(prometheus.CounterOpts{
					Name:        msg.Name,
					Help:        "dMetric counter",
					ConstLabels: msg.GetLabelPairsMap(),
				}, func() float64 {
					return float64(atomic.LoadUint64(&msg.UIntVal))
				})
			}
		}

		return msg, nil
	}
}
