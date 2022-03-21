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
	LabelMap         map[string]string
	LastPublishedVal uint64
	LastPublishedAt  time.Time
}

// NewCounter create a new counter
func NewCounter(sourceId peer.ID, name string, desc string, val uint64, labelMap map[string]string) *Counter {
	obj := Counter{
		SourceId: sourceId,
		Name:     name,
		Desc:     desc,
		uintVal:  val,
		LabelMap: labelMap,
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
		SourceId: c.SourceId,
		Type:     MetricTypeCounter,
		Name:     c.Name,
		LabelId:  GetLabelIdStrFromMap(c.LabelMap),
		UIntVal:  c.uintVal,
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

// CounterVec collection of counter with different labelIdStrs
type CounterVec struct {
	SourceId   peer.ID
	Name       string
	Desc       string
	CounterMap map[LabelIdStr]*Counter
}

// NewCounterVec create a new CounterVec
func NewCounterVec(sourceId peer.ID, name string, desc string) *CounterVec {
	cv := CounterVec{
		SourceId:   sourceId,
		Name:       name,
		Desc:       desc,
		CounterMap: make(map[LabelIdStr]*Counter),
	}
	return &cv
}

// Inc increase the value of the target counter, or create a new one if needed
func (cv *CounterVec) Inc(labelMap map[string]string) {
	labelIdStr := GetLabelIdStrFromMap(labelMap)
	c, ok := cv.CounterMap[labelIdStr]
	if !ok {
		c = NewCounter(cv.SourceId, cv.Name, cv.Desc, 0, labelMap)
		cv.CounterMap[labelIdStr] = c
	}
	c.Inc()
}

// GetValueOf return the value of the target counter
func (cv *CounterVec) GetValueOf(labelMap map[string]string) uint64 {
	labelIdStr := GetLabelIdStrFromMap(labelMap)
	c, ok := cv.CounterMap[labelIdStr]
	if !ok {
		return 0
	}
	return c.GetVal()
}
