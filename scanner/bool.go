package scanner

var (
	t = []byte("true")
	f = []byte("false")
)

// Boolean matches a boolean at the specified position
func Boolean(in []byte, pos int) (int, error) {
	switch in[pos] {
	case 't':
		return expect(in, pos, t...)
	case 'f':
		return expect(in, pos, f...)
	default:
		return 0, errUnexpectedValue
	}
}
