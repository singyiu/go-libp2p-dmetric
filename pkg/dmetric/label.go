package dmetric

import (
	"encoding/json"
	"github.com/samber/lo"
)

type LabelPair struct {
	Name   string `json:"name"`
	StrVal string `json:"strVal"`
}

// GetLabelPairsFromLabelMap return []LabelPair with a sorted map
func GetLabelPairsFromLabelMap(m map[string]string) []LabelPair {
	var output []LabelPair
	keys := lo.Keys(m)
	for _, k := range keys {
		lp := LabelPair{
			Name:   k,
			StrVal: m[k],
		}
		output = append(output, lp)
	}
	return output
}

func GetLabelIdStrFromLabelPairs(lp []LabelPair) string {
	bytes, _ := json.Marshal(lp)
	return string(bytes)
}
