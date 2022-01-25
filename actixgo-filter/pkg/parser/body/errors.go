package body

import (
	"errors"
	"fmt"
)

func UnexpectedElementError(index int, expected []Element, actual Element) error {
	return errors.New(fmt.Sprintf("json-parser [UnexpectedElementError] (index %d): Unexpected element in json. Expected: %v, Current: %v", index, expected, actual))
}

func MalformedStringError(index int) error {
	return errors.New(fmt.Sprintf("json-parser [MalformedStringError] (index %d): Missing closing quotes for string value", index))
}

func MalformedBooleanError(index int, actual string) error {
	return errors.New(fmt.Sprintf(`json-parser [MalformedBooleanError] (index %d): Expected "true"" or "false" but got %s `, index, actual))
}

func MalformedNullNodeError(index int, actual string) error {
	return errors.New(fmt.Sprintf(`json-parser [MalformedNullNodeError] (index %d): Expected "null"" but got %s `, index, actual))
}

func MalformedNumberError(index int, actual string) error {
	return errors.New(fmt.Sprintf(`json-parser [MalformedNumberError] (index %d): Expected a valid number but got %s`, index, actual))
}

func InvalidStateError(index int, actual Element) error {
	name := elementToString(actual)
	return errors.New(fmt.Sprintf(`json-parser [InvalidStateError] (index %d): Internal error. Inconsistent state. Actual %v`, index, name))
}

func MalformedObjectNode(index int, parts []Element) error {
	names := elementsToString(parts)
	return errors.New(fmt.Sprintf(`json-parser [MalformedObjectNode] (index %d): Invalid object structure. Actual parts (%d): %v`, index, len(names), names))
}

func MalformedArrayNode(index int, parts []Element) error {
	names := elementsToString(parts)
	return errors.New(fmt.Sprintf(`json-parser [MalformedArrayNode] (index %d): Invalid array structure. Actual parts (%d): %v`, index, len(names), names))
}

func elementToString(e Element) string {
	return elementNames[e]
}
func elementsToString(elements []Element) []string {
	var names []string
	for _, e := range elements {
		names = append(names, elementToString(e))
	}
	return names
}
