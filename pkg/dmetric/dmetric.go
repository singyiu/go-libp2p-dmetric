package dmetric

import "github.com/libp2p/go-libp2p-core/peer"

type MetricType uint64

const (
	MetricTypeUndefined MetricType = iota
	MetricTypeCounter
)

type Message struct {
	SourceId   peer.ID     `json:"sourceId"`
	Type       MetricType  `json:"type"`
	Name       string      `json:"name"`
	LabelPairs []LabelPair `json:"labelPairs"`
	UIntVal    uint64      `json:"uintVal"`
}

func (m Message) GetLabelPairsStr() string {
	return LabelPairs(m.LabelPairs).String()
}

func (m Message) GetLabelPairsMap() map[string]string {
	output := LabelPairs(m.LabelPairs).ToMap()
	output["sourceId"] = m.SourceId.String()
	return output
}
