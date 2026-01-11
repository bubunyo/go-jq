package scanner

// FindTo finds the elements of an array between the specified indexes; inclusive
func FindTo(in []byte, pos, to int) ([]byte, error) {
	pos, err := skipSpace(in, pos)
	if err != nil {
		return nil, err
	}

	if v := in[pos]; v != '[' {
		return nil, newError(pos, v)
	}
	pos++

	idx := 0
	itemStart := pos

	for {
		pos, err = skipSpace(in, pos)
		if err != nil {
			return nil, err
		}

		// data
		pos, err = Any(in, pos)
		if err != nil {
			return nil, err
		}

		if idx == to {
			data := in[itemStart:pos]
			result := make([]byte, 0, len(data)+2)
			result = append(result, '[')
			result = append(result, data...)
			result = append(result, ']')
			return result, nil
		}

		pos, err = skipSpace(in, pos)
		if err != nil {
			return nil, err
		}

		switch in[pos] {
		case ',':
			pos++
		case ']':
			return nil, errIndexOutOfBounds
		}

		idx++
	}
}
