package format

import "math"

func FmtNumber(data []byte) float64 {
	var sum float64

	for i, val := range data {
		sum += float64(val) * math.Pow(256, float64(i))
	}
	return sum
}
