package scanner

// Any returns the position of the end of the current element that begins at pos; handles any valid json element
func Any(in []byte, pos int) (int, error) {
	pos, err := skipSpace(in, pos)
	if err != nil {
		return 0, err
	}

	switch in[pos] {
	case '"':
		return String(in, pos)
	case '{':
		return Object(in, pos)
	case '.', '-', '1', '2', '3', '4', '5', '6', '7', '8', '9', '0':
		return Number(in, pos)
	case '[':
		return Array(in, pos)
	case 't', 'f':
		return Boolean(in, pos)
	case 'n':
		return Null(in, pos)
	default:
		max := len(in) - pos
		if max > 20 {
			max = 20
		}

		return 0, opErr{
			pos:     pos,
			msg:     "invalid object",
			content: string(in[pos : pos+max]),
		}
	}
}
