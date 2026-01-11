package scanner_test

import (
	"testing"

	"github.com/bubunyo/jq/scanner"
)

func TestFindFrom(t *testing.T) {
	testCases := map[string]struct {
		In       string
		From     int
		Expected string
		HasErr   bool
	}{
		"all": {
			In:       `["a","b","c","d","e"]`,
			From:     0,
			Expected: `["a","b","c","d","e"]`,
		},
		"last": {
			In:       `["a","b","c","d","e"]`,
			From:     4,
			Expected: `["e"]`,
		},
		"middle": {
			In:       `["a","b","c","d","e"]`,
			From:     2,
			Expected: `["c","d","e"]`,
		},
		"mixed": {
			In:       `["a",{"hello":"world"},"c","d","e"]`,
			From:     0,
			Expected: `["a",{"hello":"world"},"c","d","e"]`,
		},
		"out of bounds": {
			In:     `["a",{"hello":"world"},"c","d","e"]`,
			From:   20,
			HasErr: true,
		},
	}

	for label, tc := range testCases {
		t.Run(label, func(t *testing.T) {
			data, err := scanner.FindFrom([]byte(tc.In), 0, tc.From)
			if tc.HasErr {
				if err == nil {
					t.FailNow()
				}
			} else {
				if string(data) != tc.Expected {
					t.FailNow()
				}
				if err != nil {
					t.FailNow()
				}
			}
		})
	}
}
