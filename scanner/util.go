package scanner

import (
	"errors"
	"fmt"
	"unicode"
	"unicode/utf8"
)

var (
	errUnexpectedEOF    = errors.New("unexpected EOF")
	errKeyNotFound      = errors.New("key not found")
	errIndexOutOfBounds = errors.New("index out of bounds")
	errToLessThanFrom   = errors.New("to index less than from index")
	errFromOutOfBounds  = errors.New("from index out of bounds")
	errUnexpectedValue  = errors.New("unexpected value")
)

func skipSpace(in []byte, pos int) (int, error) {
	for {
		r, size := utf8.DecodeRune(in[pos:])
		if size == 0 {
			return 0, errUnexpectedEOF
		}
		if !unicode.IsSpace(r) {
			break
		}
		pos += size
	}

	return pos, nil
}

func expect(in []byte, pos int, content ...byte) (int, error) {
	if pos+len(content) > len(in) {
		return 0, errUnexpectedEOF
	}

	for _, b := range content {
		if v := in[pos]; v != b {
			return 0, errUnexpectedValue
		}
		pos++
	}

	return pos, nil
}

func newError(pos int, b byte) error {
	return fmt.Errorf("invalid character at position, %v; %v", pos, string([]byte{b}))
}
