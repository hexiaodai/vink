package common

import (
	"encoding/json"
	"strings"

	corev1 "k8s.io/api/core/v1"
)

const (
	searchPartsSep = ';'
)

func SearchParts(search string) []string {
	parts := []string{}
	start := 0
	l := len(search)
	trimPart := func(origin string) string {
		return strings.ReplaceAll(origin, ";;", ";")
	}
	for i := 0; i < l-1; i++ {
		if search[i] == searchPartsSep {
			if search[i+1] == searchPartsSep {
				// double ;, skip
				i++
				continue
			} else {
				// need split
				parts = append(parts, trimPart(search[start:i]))
				start = i + 1
			}
		}
	}
	parts = append(parts, trimPart(search[start:]))
	ps := []string{}
	for _, p := range parts {
		if p != "" {
			ps = append(ps, p)
		}
	}
	return ps
}

func (r *Resources) AsKubeResources() corev1.ResourceRequirements {
	rr := corev1.ResourceRequirements{}
	bs, _ := json.Marshal(r)
	json.Unmarshal(bs, &rr)
	return rr
}

func (r *Resources) UnmarshalFromKubeResources(rr corev1.ResourceRequirements) {
	bs, _ := json.Marshal(rr)
	json.Unmarshal(bs, r)
}

func (a *Affinity) AsKubeAffinity() *corev1.Affinity {
	af := &corev1.Affinity{}
	bs, _ := json.Marshal(a)
	json.Unmarshal(bs, af)
	return af
}

func (a *Affinity) UnmarshalFromKubeAffinity(af *corev1.Affinity) {
	bs, _ := json.Marshal(af)
	json.Unmarshal(bs, a)
}
