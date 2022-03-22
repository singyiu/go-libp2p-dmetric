package dmetric

import (
	"encoding/json"
	"github.com/samber/lo"
)

type LabelPair struct {
	Name   string `json:"name"`
	StrVal string `json:"strVal"`
}

type LabelPairs []LabelPair

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

func (lps LabelPairs) String() string {
	bytes, _ := json.Marshal(lps)
	return string(bytes)
}

func (lps LabelPairs) ToMap() map[string]string {
	output := make(map[string]string)
	for _, lp := range lps {
		output[lp.Name] = lp.StrVal
	}
	return output
}
