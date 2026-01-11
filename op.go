package jq

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/bubunyo/jq/scanner"
)

// Op defines a single transformation to be applied to a []byte
type Op interface {
	Apply([]byte) ([]byte, error)
}

// OpFunc provides a convenient func type wrapper on Op
type OpFunc func([]byte) ([]byte, error)

// Apply executes the transformation defined by OpFunc
func (fn OpFunc) Apply(in []byte) ([]byte, error) {
	return fn(in)
}

// Dot extract the specific key from the map provided; to extract a nested value, use the Dot Op in conjunction with the
// Chain Op
func Dot(key string) OpFunc {
	key = strings.TrimSpace(key)
	if key == "" {
		return func(in []byte) ([]byte, error) { return in, nil }
	}

	k := []byte(key)

	return func(in []byte) ([]byte, error) {
		return scanner.FindKey(in, 0, k)
	}
}

// Chain executes a series of operations in the order provided
func Chain(filters ...Op) OpFunc {
	return func(in []byte) ([]byte, error) {
		if filters == nil {
			return in, nil
		}

		var err error
		data := in
		for _, filter := range filters {
			data, err = filter.Apply(data)
			if err != nil {
				return nil, err
			}
		}

		return data, nil
	}
}

// Index extracts a specific element from the array provided
func Index(index int) OpFunc {
	return func(in []byte) ([]byte, error) {
		return scanner.FindIndex(in, 0, index)
	}
}

// Range extracts a selection of elements from the array provided, inclusive
func Range(from, to int) OpFunc {
	return func(in []byte) ([]byte, error) {
		return scanner.FindRange(in, 0, from, to)
	}
}

// From extracts all elements from the array provided from the given index onward, inclusive
func From(from int) OpFunc {
	return func(in []byte) ([]byte, error) {
		return scanner.FindFrom(in, 0, from)
	}
}

// To extracts all elements from the array provided up to the given index, inclusive
func To(to int) OpFunc {
	return func(in []byte) ([]byte, error) {
		return scanner.FindTo(in, 0, to)
	}
}

// B64Decode decodes a base64-encoded JSON string value
func B64Decode() OpFunc {
	return func(in []byte) ([]byte, error) {
		// Remove leading/trailing whitespace
		in = bytes.TrimSpace(in)

		// Check if input is a quoted JSON string
		if len(in) < 2 || in[0] != '"' || in[len(in)-1] != '"' {
			return nil, fmt.Errorf("b64_decode expects a JSON string")
		}

		// Extract the string content (without quotes)
		encoded := in[1 : len(in)-1]

		// Decode base64
		decoded, err := base64.StdEncoding.DecodeString(string(encoded))
		if err != nil {
			return nil, fmt.Errorf("b64_decode failed: %v", err)
		}

		// Return decoded bytes (should be valid JSON)
		return decoded, nil
	}
}
