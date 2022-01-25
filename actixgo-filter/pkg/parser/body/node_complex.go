package body

type checkFn func(int, []Element) error
type processFn func(*xIterator, []Element)

func (iter *xIterator) ParseComplex(phases map[int][]Element, end byte, check checkFn, process processFn) ([]Element, error) {
	var parts []Element
	current := 0

	for token, hasMore := iter.next(); token != end && hasMore; token, hasMore = iter.next() {
		expected := phases[current]
		if isSpace(token) {
			continue
		}
		iter.element = findCurrElement(token)
		parts = append(parts, iter.element)

		if !contains(expected, iter.element) {
			return nil, UnexpectedElementError(iter.index, expected, iter.element)
		}
		if err := iter.processCurrent(); err != nil {
			return nil, err
		}

		process(iter, parts)
		current = nextPhase(current, phases)
	}
	return parts, check(iter.index, parts)
}

func nextPhase(current int, phases map[int][]Element) int {
	return (current + 1) % len(phases)
}
