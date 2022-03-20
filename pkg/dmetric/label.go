package dmetric

import (
	"github.com/samber/lo"
	"sort"
	"strings"
)

type LabelIdStr string

func GetLabelIdStrFromMap(m map[string]string) LabelIdStr {
	if len(m) == 0 {
		return ""
	}

	keys := lo.Keys(m)
	sort.Strings(keys)
	var sb strings.Builder
	for _, k := range keys {
		sb.WriteString(k)
		sb.WriteString(":")
		sb.WriteString(m[k])
		sb.WriteString(",")
	}
	return LabelIdStr(sb.String())
}
