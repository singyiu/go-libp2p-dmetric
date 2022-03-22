package dmetric

import (
	"encoding/json"
	"github.com/libp2p/go-libp2p-core/peer"
	"sync/atomic"
	"time"
)

// Counter dMetric uint counter struct
type Counter struct {
	SourceId         peer.ID
	Name             string
	Desc             string
	uintVal          uint64
	LabelPairs       []LabelPair
	LastPublishedVal uint64
	LastPublishedAt  time.Time
}

// NewCounter create a new counter
func NewCounter(sourceId peer.ID, name string, desc string, val uint64, labelPairs []LabelPair) *Counter {
	obj := Counter{
		SourceId:   sourceId,
		Name:       name,
		Desc:       desc,
		uintVal:    val,
		LabelPairs: labelPairs,
	}
	return &obj
}

// Inc increase the counter value by 1
func (c *Counter) Inc() {
	atomic.AddUint64(&c.uintVal, 1)
}

// GetVal return the value of the counter
func (c *Counter) GetVal() uint64 {
	ival := atomic.LoadUint64(&c.uintVal)
	return ival
}

// ToJsonBytes for Marshalable interface
// return json.Marshal of Message
func (c *Counter) ToJsonBytes() ([]byte, error) {
	msg := Message{
		SourceId:   c.SourceId,
		Type:       MetricTypeCounter,
		Name:       c.Name,
		LabelPairs: c.LabelPairs,
		UIntVal:    c.uintVal,
	}
	return json.Marshal(msg)
}

// ShouldBePublished for Publishable interface
// return true if LastPublishedVal is not the same as uintVal
func (c *Counter) ShouldBePublished() bool {
	return c.LastPublishedVal != c.uintVal
}

// OnPublished for Publishable interface
// update LastPublishedVal upon published
func (c *Counter) OnPublished() {
	c.LastPublishedVal = c.uintVal
	c.LastPublishedAt = time.Now()
}

// CounterVec =============================

// CounterVec collection of counters
type CounterVec struct {
	SourceId   peer.ID
	Name       string
	Desc       string
	CounterMap map[string]*Counter
}

// NewCounterVec create a new CounterVec
func NewCounterVec(sourceId peer.ID, name string, desc string) *CounterVec {
	cv := CounterVec{
		SourceId:   sourceId,
		Name:       name,
		Desc:       desc,
		CounterMap: make(map[string]*Counter),
	}
	return &cv
}

// Inc increase the value of the target counter, or create a new one if needed
func (cv *CounterVec) Inc(labelPairs []LabelPair) {
	labelPairsStr := LabelPairs(labelPairs).String()
	c, ok := cv.CounterMap[labelPairsStr]
	if !ok {
		c = NewCounter(cv.SourceId, cv.Name, cv.Desc, 0, labelPairs)
		cv.CounterMap[labelPairsStr] = c
	}
	c.Inc()
}

// GetValueOf return the value of the target counter
func (cv *CounterVec) GetValueOf(labelPairs []LabelPair) uint64 {
	c, ok := cv.CounterMap[LabelPairs(labelPairs).String()]
	if !ok {
		return 0
	}
	return c.GetVal()
}
