package headers

import (
	"fmt"
)

func ParseHeader(headers [][2]string) map[string]int {
	var keys []string
	items := make(map[string]int)
	for _, header := range headers {
		if header[0] == "" {
			continue
		}
		path := extractPath(header[0])
		if _, ok := items[path]; ok {
			items[path]++
		} else {
			items[path] = 1
			keys = append(keys, path)
		}
	}
	return items
}

func extractPath(expression string) string {
	return fmt.Sprintf("$.%s", expression)
}
