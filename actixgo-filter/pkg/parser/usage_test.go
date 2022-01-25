package parser

import (
	"reflect"
	"testing"
)

func TestBuildFields(t *testing.T) {
	tests := []struct {
		name      string
		given     map[string]int
		format    string
		namespace string
		expected  []Field
	}{
		{
			"headers transform",
			map[string]int{
				"$.Accept":                  1,
				"$.Accept-Encoding":         1,
				"$.Cache-Control":           1,
				"$.Content-Length":          1,
				"$.Content-Type":            1,
				"$.Host":                    1,
				"$.Postman-Token":           1,
				"$.User-Agent":              1,
				"$.X-B3-Parentspanid":       1,
				"$.X-B3-Sampled":            1,
				"$.X-B3-Spanid":             1,
				"$.X-B3-Traceid":            1,
				"$.X-Envoy-Attempt-Count":   1,
				"$.X-Envoy-Internal":        1,
				"$.X-Forwarded-Client-Cert": 1,
			},
			"properties",
			"header",
			[]Field{
				{"properties", "header", "$.Accept", 1},
				{"properties", "header", "$.Accept-Encoding", 1},
				{"properties", "header", "$.Cache-Control", 1},
				{"properties", "header", "$.Content-Length", 1},
				{"properties", "header", "$.Content-Type", 1},
				{"properties", "header", "$.Host", 1},
				{"properties", "header", "$.Postman-Token", 1},
				{"properties", "header", "$.User-Agent", 1},
				{"properties", "header", "$.X-B3-Parentspanid", 1},
				{"properties", "header", "$.X-B3-Sampled", 1},
				{"properties", "header", "$.X-B3-Spanid", 1},
				{"properties", "header", "$.X-B3-Traceid", 1},
				{"properties", "header", "$.X-Envoy-Attempt-Count", 1},
				{"properties", "header", "$.X-Envoy-Internal", 1},
				{"properties", "header", "$.X-Forwarded-Client-Cert", 1},
			},
		},
		{
			"body transform",
			map[string]int{
				"$.$and.[*].subject_type.$in.[*]":             1,
				"$.$and.[*].created_at.$gte":                  1,
				"$.$and.[*].created_at.$lte":                  1,
				"$.$and.[*].status.$in.[*]":                   1,
				"$.$and.[*].status.$in.[*].crazy":             1,
				"$.$and.[*].$or.[*].subject.identity.$in.[*]": 1,
				"$.$and.[*].$or.[*].subject.id.$in.[*]":       1,
				"$.$and.[*].$or.[*].type.$in.[*]":             1,
			},
			"json",
			"body",
			[]Field{
				{"json", "body", "$.$and.[*].subject_type.$in.[*]", 1},
				{"json", "body", "$.$and.[*].created_at.$gte", 1},
				{"json", "body", "$.$and.[*].created_at.$lte", 1},
				{"json", "body", "$.$and.[*].status.$in.[*]", 1},
				{"json", "body", "$.$and.[*].status.$in.[*].crazy", 1},
				{"json", "body", "$.$and.[*].$or.[*].subject.identity.$in.[*]", 1},
				{"json", "body", "$.$and.[*].$or.[*].subject.id.$in.[*]", 1},
				{"json", "body", "$.$and.[*].$or.[*].type.$in.[*]", 1},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := buildFields(tt.given, tt.format, tt.namespace)
			if !reflect.DeepEqual(actual, tt.expected) {
				t.Fatalf("case [%s]: Fail. \n Expected: %v \n Actual:   %v", tt.name, tt.expected, actual)
			}
		})
	}
}
