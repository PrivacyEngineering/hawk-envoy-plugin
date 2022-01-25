package body

import (
	"reflect"
	"testing"
)

func TestParseBody(t *testing.T) {
	tests := []struct {
		name     string
		body     []byte
		expected map[string]int
	}{
		{
			"nil",
			nil,
			nil,
		},
		{
			"empty string",
			[]byte(` `),
			map[string]int{},
		},
		{
			"number",
			[]byte(`1234567980`),
			map[string]int{},
		},
		{
			"decimal",
			[]byte(`12345.67980`),
			map[string]int{},
		},
		{
			"boolean true",
			[]byte(`true`),
			map[string]int{},
		},
		{
			"boolean false",
			[]byte(`false`),
			map[string]int{},
		},
		{
			"single string",
			[]byte(`"single value"`),
			map[string]int{},
		},
		{
			"empty object",
			[]byte(`{}`),
			map[string]int{},
		},
		{
			"empty array",
			[]byte(`[]`),
			map[string]int{},
		},
		{
			"array of strings",
			[]byte(`["aa","bb","cc"]`),
			map[string]int{"$.[*]": 1},
		},
		{
			"array of object",
			[]byte(`[{"a":1},{"a":2}]`),
			map[string]int{"$.[*].a": 2},
		},
		{
			"array of object and bool",
			[]byte(`[{"a":1},true]`),
			map[string]int{
				"$.[*].a": 1,
				"$.[*]":   1,
			},
		},
		{
			"single property",
			[]byte(`{"single": "property"}`),
			map[string]int{"$.single": 1},
		},
		{
			"two properties",
			[]byte(`{"one": "11111", "two": "2222"}`),
			map[string]int{
				"$.one": 1,
				"$.two": 1,
			},
		},
		{
			"three properties with numbers",
			[]byte(`{"one": 1111, "two": 2222, "three": 3333}`),
			map[string]int{
				"$.one":   1,
				"$.two":   1,
				"$.three": 1,
			},
		},
		{
			"four properties",
			[]byte(`{"one": "1111", "two": 99.88, "three": "all-lll", "four": true}`),
			map[string]int{
				"$.one":   1,
				"$.two":   1,
				"$.three": 1,
				"$.four":  1,
			},
		},
		{
			"nested 2 levels",
			[]byte(`{"one": { "two": 2 } }`),
			map[string]int{
				"$.one.two": 1,
			},
		},
		{
			"nested 3 levels",
			[]byte(`{"one": { "two": { "three": true } } }`),
			map[string]int{
				"$.one.two.three": 1,
			},
		},
		{
			"2 nested 2 levels",
			[]byte(`{ "p_one": { "c_one": 1, "c_two": 2 }, "p_two": { "c_one": 1, "c_two": 2 } }`),
			map[string]int{
				"$.p_one.c_one": 1,
				"$.p_one.c_two": 1,
				"$.p_two.c_one": 1,
				"$.p_two.c_two": 1,
			},
		},
		{
			"array of objects",
			[]byte(`[ { "user": { "email":"a", "lastName":"aa" }, "p":false }, { "user": { "email":"b", "lastName":"bb" }, "p":true } ]`),
			map[string]int{
				"$.[*].user.email":    2,
				"$.[*].user.lastName": 2,
				"$.[*].p":             2,
			},
		},
		{
			"array of objects",
			[]byte(`
{
  "$and": [
    {"subject_type": {"$in": ["xx"]}},
    {"created_at": {"$gte": "xx", "$lte": "xx" } },
    { "status": { "$in": [ "created", { "crazy": "yes", "$blot": [] } ] } },
    {
      "$or": [
        { "subject.identity": { "$in": [ "xx" ] } },
        { "subject.id": { "$in": [ "xx" ] } }
      ]
    },
    {
      "$or": [
        {"type": {"$in": ["xx", "yy"]}},
        {"payload.peculiarity_type": { "$in": [] }}
      ]
    }
  ]
}
`),
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
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := ParseBody(tt.body)
			if err != nil {
				t.Errorf("fail with error: %v", err)
			}

			//sort.SliceStable(actual, func(i, j int) bool {
			//	return actual[i].Path < actual[j].Path
			//})
			//sort.SliceStable(tt.expected, func(i, j int) bool {
			//	return tt.expected[i].Path < tt.expected[j].Path
			//})
			if !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("case [%s]: Fail. \n Expected: %v \n Actual:   %v", tt.name, tt.expected, actual)
			}
		})
	}
}
