package dmetric

type MetricType uint64

const (
	MetricTypeUndefined MetricType = iota
	MetricTypeCounter
)

type DMetricMessage struct {
	SourceId string     `json:"sourceId"`
	Type     MetricType `json:"type"`
	Name     string     `json:"name"`
	LabelId  LabelIdStr `json:"labelId"`
	UIntVal  uint64     `json:"uintVal"`
}
