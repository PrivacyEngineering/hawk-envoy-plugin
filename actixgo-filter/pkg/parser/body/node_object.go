package body

const endObject = byte('}')

var phasesObject = map[int][]Element{
	0: {String},
	1: {Colon},
	2: {String, Number, ObjectNode, ArrayNode, Boolean, NullNode},
	3: {Comma},
}

// ParseObject start with { and will end with }
// structure { string : value (, string : value )* }
func (iter *xIterator) ParseObject() error {
	parts, err := iter.ParseComplex(phasesObject, endObject, checkObject, chainProcessObject)
	if err != nil {
		return err
	}
	if len(parts) > 0 {
		prev := parts[len(parts)-1]
		if prev == String || prev == Number || prev == Boolean {
			iter.addPath()
		}
	}
	iter.removeLastOfChain()
	return err
}

// string : value => 3
// string : value , string : value => 7 = 3 + 4
// string : value , string : value , string : value => 11 = 3 + 4 + 4
func checkObject(index int, parts []Element) error {
	size := len(parts)
	if size == 0 || (size-3)%4 == 0 {
		return nil
	}
	return MalformedObjectNode(index, parts)
}

func chainProcessObject(iter *xIterator, parts []Element) {
	if (len(parts)-1)%4 == 0 {
		iter.addChain(iter.content)
	}
	if parts[len(parts)-1] == Comma {
		prev := parts[len(parts)-2]
		if prev == String || prev == Number || prev == Boolean {
			iter.addPath()
		}
		iter.removeLastOfChain()
	}
}
