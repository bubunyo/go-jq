package jq

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	reArray = regexp.MustCompile(`^\s*\[\s*(?:(\d+))?\s*(?:(:))?\s*(?:(\d+))?\s*\]\s*$`)
)

// Must is a convenience method similar to template.Must
func Must(op Op, err error) Op {
	if err != nil {
		panic(fmt.Errorf("unable to parse selector; %v", err.Error()))
	}

	return op
}

// Parse takes a string representation of a selector and returns the corresponding Op definition
func Parse(selector string) (Op, error) {
	// Split selector by pipe operator to get pipeline segments
	pipeSegments := strings.Split(selector, "|")

	ops := make([]Op, 0, len(pipeSegments))
	for _, segment := range pipeSegments {
		op, err := parseSegment(segment)
		if err != nil {
			return nil, err
		}
		ops = append(ops, op)
	}

	return Chain(ops...), nil
}

func parseArray(key string) (Op, bool) {
	match := FindIndices(key)

	if len(match) == 0 {
		return nil, false
	}

	matches := match[0]

	if matches[1]+matches[2]+matches[3] == "" {
		return From(0), true
	}

	if matches[2] == "" {
		idx, err := strconv.Atoi(matches[1])
		if err != nil {
			return nil, false
		}
		return Index(idx), true
	}

	if matches[1] == "" {
		to, err := strconv.Atoi(matches[3])
		if err != nil {
			return nil, false
		}
		return To(to), true
	}

	if matches[3] == "" {
		from, err := strconv.Atoi(matches[1])
		if err != nil {
			return nil, false
		}
		return From(from), true
	}

	from, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil, false
	}

	to, err := strconv.Atoi(matches[3])
	if err != nil {
		return nil, false
	}

	return Range(from, to), true
}

func FindIndices(key string) [][]string {
	return reArray.FindAllStringSubmatch(key, -1)
}

// namedOperations is a registry of operations that can be referenced by name
var namedOperations = map[string]func() OpFunc{
	"b64_decode": B64Decode,
}

// isNamedOperation checks if a segment name corresponds to a registered operation
func isNamedOperation(name string) bool {
	_, exists := namedOperations[name]
	return exists
}

// getNamedOperation retrieves a named operation from the registry
func getNamedOperation(name string) (Op, error) {
	fn, exists := namedOperations[name]
	if !exists {
		// This should not happen if isNamedOperation is checked first
		return nil, fmt.Errorf("unknown operation: %s", name)
	}
	return fn(), nil
}

// parseDotNotation parses a dot-notation selector (e.g., "a.b.c" or "a.[0]")
func parseDotNotation(selector string) (Op, error) {
	segments := strings.Split(selector, ".")

	ops := make([]Op, 0, len(segments))
	for _, segment := range segments {
		key := strings.TrimSpace(segment)
		if key == "" {
			continue
		}

		if op, ok := parseArray(key); ok {
			ops = append(ops, op)
			continue
		}

		ops = append(ops, Dot(key))
	}

	return Chain(ops...), nil
}

// parseSegment parses a single pipe segment, which can be either a named operation or a dot-notation selector
func parseSegment(segment string) (Op, error) {
	segment = strings.TrimSpace(segment)

	// Check for empty segments
	if segment == "" {
		return nil, fmt.Errorf("empty pipe segment")
	}

	// Check if it's a named operation (no dots or brackets)
	if !strings.Contains(segment, ".") && !strings.Contains(segment, "[") {
		// Try to get as named operation
		if isNamedOperation(segment) {
			return getNamedOperation(segment)
		}
	}

	// Otherwise, parse as dot notation selector
	return parseDotNotation(segment)
}
