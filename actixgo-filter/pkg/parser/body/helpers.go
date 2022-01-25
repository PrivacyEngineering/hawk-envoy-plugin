package body

import (
	"bytes"
)

type Element int

const (
	String Element = iota
	Number
	ObjectNode
	ArrayNode
	Boolean
	NullNode
	Comma
	Colon
	Invalid
)

var elementNames = map[Element]string{
	String:     "String",
	Number:     "Number",
	ObjectNode: "ObjectNode",
	ArrayNode:  "ArrayNode",
	Boolean:    "Boolean",
	NullNode:   "NullNode",
	Comma:      "Comma",
	Colon:      "Colon",
	Invalid:    "Invalid",
}

func isSpace(token byte) bool {
	switch token {
	case ' ', '\n', '\t', '\r':
		return true
	default:
		return false
	}
}
func isString(token byte) bool {
	return token == '"'
}
func isNumber(token byte) bool {
	switch token {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '-':
		return true
	default:
		return false
	}
}
func isObject(token byte) bool {
	return token == '{'
}
func isArray(token byte) bool {
	return token == '['
}
func isBoolean(token byte) bool {
	switch token {
	case 't', 'f':
		return true
	default:
		return false
	}
}
func isNullNode(token byte) bool {
	return token == 'n'
}
func isColon(token byte) bool {
	return token == ':'
}
func isComma(token byte) bool {
	return token == ','
}
func isDigit(token byte) bool {
	switch token {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '-', '.', 'e', 'E':
		return true
	default:
		return false
	}
}

func contains(expected []Element, value Element) bool {
	for _, element := range expected {
		if element == value {
			return true
		}
	}
	return false
}

func findCurrElement(token byte) Element {
	switch {
	case isString(token):
		return String
	case isNumber(token):
		return Number
	case isObject(token):
		return ObjectNode
	case isArray(token):
		return ArrayNode
	case isBoolean(token):
		return Boolean
	case isNullNode(token):
		return NullNode
	case isComma(token):
		return Comma
	case isColon(token):
		return Colon
	default:
		return Invalid
	}
}

func buildPath(chain []string) string {
	buff := bytes.NewBufferString("")
	for i, cc := range chain {
		buff.WriteString(cc)
		if i != len(chain)-1 {
			buff.WriteString(".")
		}
	}
	return buff.String()
}
