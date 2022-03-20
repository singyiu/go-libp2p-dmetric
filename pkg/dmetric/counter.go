package dmetric

import (
	"sync/atomic"
	"time"
)

type Counter struct {
	Name string
	Desc string
	val  uint64
	LabelMap map[string]string
	LastPublishedVal uint64
	LastPublishedAt time.Time
}

func NewCounter(name string, desc string, val uint64, labelMap map[string]string) *Counter {
	obj := Counter{
		Name:     name,
		Desc:     desc,
		val:   val,
		LabelMap: labelMap,
	}
	return &obj
}

func (c *Counter) Inc() {
	atomic.AddUint64(&c.val, 1)
}

func (c *Counter) GetVal() uint64 {
	ival := atomic.LoadUint64(&c.val)
	return ival
}

type CounterVec struct {
	Name string
	Desc string
	CounterMap map[LabelIdStr]*Counter
}

func NewCounterVec(name string, desc string) *CounterVec {
	cv := CounterVec{
		Name:     name,
		Desc:     desc,
		CounterMap: make(map[LabelIdStr]*Counter),
	}
	return &cv
}

func (cv *CounterVec) Inc(labelMap map[string]string) {
	labelIdStr := GetLabelIdStrFromMap(labelMap)
	c, ok := cv.CounterMap[labelIdStr]
	if !ok {
		c = NewCounter(cv.Name, cv.Desc, 0, labelMap)
		cv.CounterMap[labelIdStr] = c
	}
	c.Inc()
}

func (cv *CounterVec) GetValueOf(labelMap map[string]string) uint64 {
	labelIdStr := GetLabelIdStrFromMap(labelMap)
	c, ok := cv.CounterMap[labelIdStr]
	if !ok {
		return 0
	}
	return c.GetVal()
}