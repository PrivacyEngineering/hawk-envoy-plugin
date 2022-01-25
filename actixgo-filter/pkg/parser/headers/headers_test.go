package headers

import (
	"reflect"
	"testing"
)

func TestParseHeader(t *testing.T) {
	tests := []struct {
		name     string
		given    [][2]string
		expected map[string]int
	}{
		{
			"simple header",
			[][2]string{{"Content-Length", "23"}},
			map[string]int{"$.Content-Length": 1},
		},
		{
			"same header twice",
			[][2]string{{"Content-Length", "23"}, {"Content-Length", "6544"}},
			map[string]int{"$.Content-Length": 2},
		},
		{
			"empty",
			[][2]string{},
			map[string]int{},
		},
		{
			"nil",
			nil,
			map[string]int{},
		},
		{
			"empty header",
			[][2]string{{}},
			map[string]int{},
		},
		{
			"real example",
			[][2]string{
				{"Accept", "*/*"},
				{"Accept-Encoding", "gzip, deflate, br"},
				{"Cache-Control", "no-cache"},
				{"Content-Length", "23"},
				{"Content-Type", "application/json"},
				{"Host", "10.111.48.236"},
				{"Postman-Token", "13b6bb46-c82c-413a-bdcf-5b565dc87382"},
				{"User-Agent", "PostmanRuntime/7.29.0"},
				{"X-B3-Parentspanid", "15fb3b9c18415168"},
				{"X-B3-Sampled", "0"},
				{"X-B3-Spanid", "0fa56070b66f8014"},
				{"X-B3-Traceid", "746f0bce27226f6d15fb3b9c18415168"},
				{"X-Envoy-Attempt-Count", "1"},
				{"X-Envoy-Internal", "true"},
				{"X-Forwarded-Client-Cert", "By=spiffe://cluster.local/ns/httpbin-gateway/sa/httpbin;Hash=65cc9fa283682d39a8698afca4439317a68d6e8946cf9bce1199ae56b78288d9;Subject=\"\";URI=spiffe://cluster.local/ns/istio-system/sa/istio-ingressgateway"},
			},
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
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := ParseHeader(tt.given)
			if !reflect.DeepEqual(actual, tt.expected) {
				t.Fatalf("case [%s]: Fail. \n Expected: %v \n Actual:   %v", tt.name, tt.expected, actual)
			}
		})
	}
}
