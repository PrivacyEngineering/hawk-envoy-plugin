package body

import (
	"bytes"
	"math/big"
)

// stop in the current (until ")
func (iter *xIterator) consumeString() error {
	buff := bytes.NewBufferString("")
	for token, hasMore := iter.next(); hasMore && token != '"'; token, hasMore = iter.next() {
		if len(iter.data)-1 == iter.index {
			return MalformedStringError(iter.index)
		}
		buff.WriteByte(token)
		if token == '\\' {
			token, hasMore = iter.next()
			if len(iter.data)-1 == iter.index {
				return MalformedStringError(iter.index)
			}
			buff.WriteByte(token)
		}
	}
	iter.content = buff.String()
	return nil
}

// stop in the current (until e)
func (iter *xIterator) consumeBoolean() error {
	token, content := iter.data[iter.index], ""

	if 't' == token {
		content = iter.readCurrentPlus(3)
	} else if 'f' == token {
		content = iter.readCurrentPlus(4)
	}

	if content != "true" && content != "false" {
		return MalformedBooleanError(iter.index, content)
	}
	return nil
}

// stop in the current (until l)
func (iter *xIterator) consumeNullNode() error {
	content := iter.readCurrentPlus(3)
	if content != "null" {
		return MalformedNullNodeError(iter.index, content)
	}
	return nil
}

// stop in the current token (until last digit)
func (iter *xIterator) consumeNumber() error {
	token, hasMore := iter.data[iter.index], true
	buff := bytes.NewBufferString("")
	for ; isDigit(token) && hasMore; token, hasMore = iter.next() {
		buff.WriteByte(token)
	}
	content := buff.String()
	_, _, err := big.ParseFloat(content, 10, 0, big.ToNearestEven)
	if err != nil {
		return MalformedNumberError(iter.index, content)
	}
	//return to last digit position
	if hasMore {
		iter.previous()
	}
	return nil
}

func (iter *xIterator) readCurrentPlus(times int) string {
	token, hasMore := iter.data[iter.index], true
	buff := bytes.NewBufferString("")
	buff.WriteByte(token)
	for i := 0; hasMore && i < times; i++ {
		token, hasMore = iter.next()
		buff.WriteByte(token)
	}

	return buff.String()
}
