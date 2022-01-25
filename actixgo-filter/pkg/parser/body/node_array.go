package body

const endArray = byte(']')

var phasesArray = map[int][]Element{
	0: {String, Number, ObjectNode, ArrayNode, Boolean, NullNode},
	1: {Comma},
}

// ParseArray start with [ and will end with ]
// structure [ value (, value)* ]
func (iter *xIterator) ParseArray() error {
	iter.addChain("[*]")
	parts, err := iter.ParseComplex(phasesArray, endArray, checkArray, func(*xIterator, []Element) {})
	if contains(parts, String) || contains(parts, Number) || contains(parts, Boolean) {
		iter.addPath()
	}
	iter.removeLastOfChain()
	return err
}

// value => 1
// value , value => 3 = 1 + 2
// value , value , value => 5 = 1 + 2 + 2
func checkArray(index int, parts []Element) error {
	size := len(parts)
	if size == 0 || (size-1)%2 == 0 {
		return nil
	}
	return MalformedArrayNode(index, parts)
}
