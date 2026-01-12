# jq

A high-performance Go implementation of the incredibly useful jq command line tool.

Rather than marshalling JSON elements into Go structs, jq manipulates JSON as raw `[]byte`. This is especially useful for applications that need to handle dynamic JSON data without predefined schemas.

## Features

-  **High Performance** - Zero-copy operations, nanosecond-level performance
-  **Pipe Operator** - Chain operations together like the original jq
-  **Base64 Decoding** - Built-in support for base64-encoded JSON
-  **Zero Allocations** - Most operations allocate no memory
-  **Fully Tested** - Comprehensive test suite with testify

## Installation

```bash
go get github.com/bubunyo/go-jq
```

## Quick Start

```go
package main

import (
	"fmt"
	"github.com/bubunyo/go-jq"
)

func main() {
	// Simple key extraction
	op, _ := jq.Parse(".hello")
	data := []byte(`{"hello":"world"}`)
	value, _ := op.Apply(data)
	fmt.Println(string(value)) // Output: "world"

	// Using pipes
	op, _ = jq.Parse(".user|.name")
	data = []byte(`{"user":{"name":"Alice","age":30}}`)
	value, _ = op.Apply(data)
	fmt.Println(string(value)) // Output: "Alice"
}
```

## Syntax Reference

### Basic Selectors

| Syntax | Description | Example Input | Example Output |
|--------|-------------|---------------|----------------|
| `.` | Unchanged input | `{"a":1}` | `{"a":1}` |
| `.foo` | Value at key | `{"foo":"bar"}` | `"bar"` |
| `.foo.bar` | Nested key access | `{"foo":{"bar":"baz"}}` | `"baz"` |
| `.[0]` | Array element at index | `["a","b","c"]` | `"a"` |
| `.[1:3]` | Array slice (inclusive) | `["a","b","c","d"]` | `["b","c","d"]` |
| `.[1:]` | Array from index onward | `["a","b","c","d"]` | `["b","c","d"]` |
| `.[:2]` | Array up to index | `["a","b","c","d"]` | `["a","b","c"]` |
| `.[]` | All array elements | `["a","b","c"]` | `["a","b","c"]` |

### Advanced Features

| Syntax | Description | Example |
|--------|-------------|---------|
| `\|` | Pipe operator - chain operations | `.foo\|.bar` |
| `b64_decode` | Decode base64-encoded JSON string | `.data\|b64_decode\|.field` |

## Examples

### Working with Objects

```go
data := []byte(`{
  "user": {
    "name": "Alice",
    "email": "alice@example.com"
  }
}`)

// Extract nested value
op, _ := jq.Parse(".user.name")
result, _ := op.Apply(data)
// result: "Alice"
```

### Working with Arrays

```go
data := []byte(`{
  "users": [
    {"name": "Alice", "age": 30},
    {"name": "Bob", "age": 25},
    {"name": "Charlie", "age": 35}
  ]
}`)

// Get specific array element
op, _ := jq.Parse(".users.[1].name")
result, _ := op.Apply(data)
// result: "Bob"

// Get array slice
op, _ = jq.Parse(".users.[0:1]")
result, _ = op.Apply(data)
// result: [{"name":"Alice","age":30},{"name":"Bob","age":25}]
```

### Using the Pipe Operator

The pipe operator `|` allows you to chain multiple operations together:

```go
data := []byte(`{
  "response": {
    "users": [
      {"name": "Alice", "active": true},
      {"name": "Bob", "active": false}
    ]
  }
}`)

// Chain multiple operations
op, _ := jq.Parse(".response|.users|.[0]|.name")
result, _ := op.Apply(data)
// result: "Alice"
```

### Base64 Decoding

Decode base64-encoded JSON strings and continue processing:

```go
data := []byte(`{
  "encoded": "eyJuYW1lIjoiQWxpY2UiLCJhZ2UiOjMwfQ=="
}`)

// The base64 string decodes to: {"name":"Alice","age":30}
op, _ := jq.Parse(".encoded|b64_decode|.name")
result, _ := op.Apply(data)
// result: "Alice"
```

### Building Operations Programmatically

You can also construct operations without parsing:

```go
// Create operations directly
op := jq.Chain(
	jq.Dot("user"),
	jq.Dot("address"),
	jq.Dot("city"),
)

data := []byte(`{"user":{"address":{"city":"New York"}}}`)
result, _ := op.Apply(data)
// result: "New York"

// Array operations
op = jq.Chain(
	jq.Dot("items"),
	jq.Index(2),
)

data = []byte(`{"items":["a","b","c","d"]}`)
result, _ = op.Apply(data)
// result: "c"
```

### Error Handling

```go
data := []byte(`{"foo":"bar"}`)

op, err := jq.Parse(".baz")
if err != nil {
	// Handle parse error
	panic(err)
}

result, err := op.Apply(data)
if err != nil {
	// Handle apply error (e.g., key not found)
	fmt.Println("Error:", err)
	// Error: key not found
}
```

## Performance

This implementation is designed for high performance with minimal allocations:

```
BenchmarkAny-8         	26988705	        44.47  ns/op	       0 B/op	       0 allocs/op
BenchmarkArray-8       	18326660	        65.33  ns/op	       0 B/op	       0 allocs/op
BenchmarkAsArray-8     	 2124834	       566.5   ns/op	    6528 B/op	       1 allocs/op
BenchmarkFindIndex-8   	17996996	        66.30  ns/op	       0 B/op	       0 allocs/op
BenchmarkFindKey-8     	18916466	        63.46  ns/op	       0 B/op	       0 allocs/op
BenchmarkFindRange-8   	11332197	       106.7   ns/op	      16 B/op	       1 allocs/op
BenchmarkNull-8        201436381	         5.861 ns/op	       0 B/op	       0 allocs/op
BenchmarkNumber-8      	49194944	        22.80  ns/op	       0 B/op	       0 allocs/op
BenchmarkObject-8      	20190363	        59.70  ns/op	       0 B/op	       0 allocs/op
BenchmarkString-8      	60229246	        19.87  ns/op	       0 B/op	       0 allocs/op
```

Most operations complete in under 350 nanoseconds with zero or minimal memory allocations.

## Testing

The project includes comprehensive tests using testify:

```bash
go test -v ./...
```

## Contributing

Contributions are welcome! Please feel free to submit issues or pull requests.

## License

See [LICENSE](LICENSE) file for details.

## Acknowledgments

This project was inspired by the original [jq](https://stedolan.github.io/jq/) command-line tool and builds upon the foundation laid by [savaki/jq](https://github.com/savaki/jq).
