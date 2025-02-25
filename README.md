# CaseInsensitiveMap

[![Go Reference](https://pkg.go.dev/badge/github.com/projectbarks/cimap.svg)](https://pkg.go.dev/github.com/projectbarks/cimap)
[![Main Actions Status](https://github.com/projectbarks/cimap/workflows/Go/badge.svg)](https://github.com/projectbarks/cimap/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/projectbarks/cimap)](https://goreportcard.com/report/github.com/projectbarks/cimap)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](./LICENSE.md)

`CaseInsensitiveMap` `(cimap)` is a high performance map for case-insensitive keys.

## Features

- **Case-Insensitive Keys**: Keys are treated in a case-insensitive manner, allowing for more flexible key management.
- **Generic Support**: The map supports generic types, allowing you to store any type of value.
- **Custom Hashing**: You can set a custom hash function for the map.
- **JSON Serialization**: The map can be easily serialized and deserialized to and from JSON.
- **Iterators**: Provides iterators for keys and key-value pairs.

## Installation

You need Golang [1.23.x](https://go.dev/dl/) or above

```bash
go get github.com/projectbarks/cimap
```

## Usage & Documentation

Function documentation can be found at [GoDocs](https://pkg.go.dev/github.com/projectbarks/cimap). Here's a basic example of how to use the `CaseInsensitiveMap`:

```go
package main

import (
	"fmt"
	"github.com/projectbarks/cimap"
)

func main() {
	// Create a new case-insensitive map
	m := cimap.New[string]()

	// Add some key-value pairs
	m.Add("KeyOne", "Value1")
	m.Add("keytwo", "Value2")

	// Retrieve values
	val, found := m.Get("KEYONE")
	if found {
		fmt.Println("Found:", val)
	} else {
		fmt.Println("Key not found")
	}

	// Iterate over keys
	m.Keys()(func(key string) bool {
		fmt.Println("Key:", key)
		return true
	})

	// Serialize to JSON
	jsonData, _ := m.MarshalJSON()
	fmt.Println("JSON:", string(jsonData))
}
```

## Performance

- **Time per operation**: Over 50% speed improvement compared to native case insensitive map.
- **No additional allocations**: `CIMap` uses **0 B/op** and **0 allocs/op** for Add, Get, Delete, and more. By converting characters to lowercase inline without extra string allocations, `CIMap` avoids overhead from creating new strings.

> :warning: This code performs best when there are lots of string allocations due to `strings.ToLower` 
>			for strings either containing unicode or uppercase characters. If you can gaurantee
>           the inputs provided will only be one casing use a native map instead.

```bash
          │    sec/op     │   sec/op     vs base                │
                  │         sec/op          │    sec/op     vs base                 │
Add/Upper/16                   137.3n ±  2%   129.5n ± 11%         ~ (p=0.105 n=10)
Add/Lower/16                   22.68n ±  4%   49.48n ± 16%  +118.21% (p=0.000 n=10)
Add/Unicode/16                 363.4n ±  4%   331.9n ±  7%    -8.69% (p=0.000 n=10)
Get/Upper/16                  118.30n ±  7%   93.06n ±  7%   -21.34% (p=0.000 n=10)
Get/Lower/16                   50.27n ± 10%   96.12n ± 12%   +91.22% (p=0.000 n=10)
Get/Unicode/16                 578.0n ±  6%   618.3n ±  5%         ~ (p=0.052 n=10)
Delete/Upper/16                79.37n ±  8%   48.65n ± 10%   -38.71% (p=0.000 n=10)
Delete/Lower/16                30.38n ±  5%   49.97n ± 11%   +64.51% (p=0.000 n=10)
Delete/Unicode/16              547.1n ±  4%   468.0n ±  5%   -14.46% (p=0.000 n=10)
geomean                        119.9n         133.4n         +11.23%
```

```bash
          │     B/op      │   B/op     vs base                     │
Add/Upper/16                  122.0 ±  5%       0.0 ±  0%  -100.00% (p=0.000 n=10)
Add/Lower/16                  46.50 ±  5%      0.00 ±  0%  -100.00% (p=0.000 n=10)
Add/Unicode/16                142.5 ± 28%       0.0 ±  0%  -100.00% (p=0.000 n=10)
Get/Upper/16                  33.00 ±  0%      0.00 ±  0%  -100.00% (p=0.000 n=10)
Get/Lower/16                  0.000 ±  0%     0.000 ±  0%         ~ (p=1.000 n=10) ¹
Get/Unicode/16                112.0 ±  0%       0.0 ±  0%  -100.00% (p=0.000 n=10)
Delete/Upper/16               33.00 ±  0%      0.00 ±  0%  -100.00% (p=0.000 n=10)
Delete/Lower/16               0.000 ±  0%     0.000 ±  0%         ~ (p=1.000 n=10) ¹
Delete/Unicode/16           119.000 ±  1%     3.000 ± 33%   -97.48% (p=0.000 n=10)
```

```bash
                  │        allocs/op        │ allocs/op   vs base                     │
Add/Upper/16                   1.000 ± 0%     0.000 ± 0%  -100.00% (p=0.000 n=10)
Add/Lower/16                   0.000 ± 0%     0.000 ± 0%         ~ (p=1.000 n=10) ¹
Add/Unicode/16                 2.000 ± 0%     0.000 ± 0%  -100.00% (p=0.000 n=10)
Get/Upper/16                   1.000 ± 0%     0.000 ± 0%  -100.00% (p=0.000 n=10)
Get/Lower/16                   0.000 ± 0%     0.000 ± 0%         ~ (p=1.000 n=10) ¹
Get/Unicode/16                 1.000 ± 0%     0.000 ± 0%  -100.00% (p=0.000 n=10)
Delete/Upper/16                1.000 ± 0%     0.000 ± 0%  -100.00% (p=0.000 n=10)
Delete/Lower/16                0.000 ± 0%     0.000 ± 0%         ~ (p=1.000 n=10) ¹
Delete/Unicode/16              2.000 ± 0%     0.000 ± 0%  -100.00% (p=0.000 n=10)
```

## Contributing

Contributions are welcome! Please feel free to submit a pull request or open an issue.

### Development

Comparing the benchmark performance stats:

Setup 
```bash
go install golang.org/x/perf/cmd/benchstat@latest
```

```bash
go test -benchmem -run=^$ -bench '^(Benchmark)' github.com/projectbarks/cimap -count=10 > benchmark/bench-all.txt
grep 'Base-' benchmark/bench-all.txt | sed 's|Base-||g' > benchmark/bench-old.txt
grep 'CIMap-' benchmark/bench-all.txt | sed 's|CIMap-||g' > benchmark/bench-new.txt
benchstat benchmark/bench-old.txt benchmark/bench-new.txt

```

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
