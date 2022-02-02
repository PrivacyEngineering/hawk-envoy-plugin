package format

import (
	"reflect"
	"testing"
)

func Test(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		expected float64
	}{
		{"0", []byte{byte(0), byte(0)}, 0.0},
		{"200", []byte{byte(200), byte(0)}, 200.0},
		{"555", []byte{byte(43), byte(2)}, 555.0},
		{"400", []byte{byte(144), byte(1)}, 400.0},
		{"404", []byte{byte(148), byte(1)}, 404.0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := FmtNumber(tt.data)
			if !reflect.DeepEqual(actual, tt.expected) {
				t.Fatalf("Actual [%v] different to expected [%v]", actual, tt.expected)
			}
		})
	}
}
