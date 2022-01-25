package body

func ParseBody(data []byte) (map[string]int, error) {
	if data == nil || len(data) == 0 {
		return nil, nil
	}
	iter := xIterator{
		data:  data,
		size:  len(data),
		index: -1,
		chain: []string{"$"},
		paths: map[string]int{},
	}
	err := iter.ParseValue()
	if err != nil {
		return nil, err
	}

	return iter.paths, nil
}

type xIterator struct {
	data    []byte
	size    int
	index   int
	element Element
	content string
	chain   []string
	paths   map[string]int
}

func (iter *xIterator) ParseValue() error {
	expected := []Element{String, Number, ObjectNode, ArrayNode, Boolean, NullNode}
	for token, hasMore := iter.next(); hasMore; token, hasMore = iter.next() {
		if isSpace(token) {
			continue
		}
		iter.element = findCurrElement(token)
		if !contains(expected, iter.element) {
			return UnexpectedElementError(iter.index, expected, iter.element)
		}
		if err := iter.processCurrent(); err != nil {
			return err
		}
		expected = []Element{}
	}
	return nil
}

func (iter *xIterator) processCurrent() error {
	switch iter.element {
	case String:
		return iter.consumeString()
	case Number:
		return iter.consumeNumber()
	case ObjectNode:
		return iter.ParseObject()
	case ArrayNode:
		return iter.ParseArray()
	case Boolean:
		return iter.consumeBoolean()
	case NullNode:
		return iter.consumeNullNode()
	case Colon, Comma:
		return nil
	}

	return InvalidStateError(iter.index, iter.element)
}

func (iter *xIterator) next() (byte, bool) {
	if iter.endOfData() {
		return 0, false
	}
	iter.index++
	token := iter.data[iter.index]
	return token, true
}
func (iter *xIterator) previous() byte {
	iter.index--
	token := iter.data[iter.index]
	return token
}
func (iter *xIterator) endOfData() bool {
	return iter.size-1 == iter.index
}
func (iter *xIterator) remaining() int {
	return iter.size - iter.index
}

func (iter *xIterator) addChain(text string) {
	iter.chain = append(iter.chain, text)
}
func (iter *xIterator) removeLastOfChain() {
	iter.chain = iter.chain[:len(iter.chain)-1]
}
func (iter *xIterator) addPath() {
	path := buildPath(iter.chain)
	if val, ok := iter.paths[path]; ok {
		iter.paths[path] = val + 1
	} else {
		iter.paths[path] = 1
	}
}
