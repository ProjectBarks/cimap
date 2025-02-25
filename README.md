# CaseInsensitiveMap

## Overview

`CaseInsensitiveMap` `(cimap)` is a Go package that provides a map-like data structure with case-insensitive keys. It is designed to handle string keys in a way that ignores case differences, making it useful for scenarios where key case sensitivity is not desired.

## Features

- **Case-Insensitive Keys**: Keys are treated in a case-insensitive manner, allowing for more flexible key management.
- **Generic Support**: The map supports generic types, allowing you to store any type of value.
- **Custom Hashing**: You can set a custom hash function for the map.
- **JSON Serialization**: The map can be easily serialized and deserialized to and from JSON.
- **Iterators**: Provides iterators for keys and key-value pairs.

## Installation

To install the package, use the following command:

```bash
go get github.com/projectbarks/cimap
```

## Usage

Here's a basic example of how to use the `CaseInsensitiveMap`:

```go
package main

import (
	"fmt"
	"cimap"
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

-   **Time per operation**: Over 50% speed improvement compared to native case insensitive map.
-   **No additional allocations**: `CIMap` uses **0 B/op** and **0 allocs/op** for Add, Get, and Delete.
-   **Inline case handling**: By converting characters to lowercase inline without extra string allocations, `CIMap` avoids overhead from creating new strings.

Below is an example of benchmark data comparing a **base** case-insensitive map approach (`Base`) vs. `CaseInsensitiveMap` (`CIMap`). The results show substantial performance gains, particularly in reduced time per operation (ns/op), **no allocations** (`B/op` and `allocs/op`), and overall higher efficiency for Add, Get, and Delete operations.

```lang=bash
          │    sec/op     │   sec/op     vs base                │
Add/16       45.04n ±  9%   20.85n ± 4%  -53.69% (p=0.000 n=10)
Get/16      131.35n ±  6%   59.49n ± 9%  -54.71% (p=0.000 n=10)
Delete/16    66.89n ± 10%   22.39n ± 6%  -66.53% (p=0.000 n=10)
geomean      73.41n         30.28n       -58.75%
```

```lang=bash
          │     B/op      │   B/op     vs base                     │
Add/16        95.50 ± 39%   0.00 ± 0%  -100.00% (p=0.000 n=10)
Get/16        20.00 ±  0%   0.00 ± 0%  -100.00% (p=0.000 n=10)
Delete/16     16.00 ±  0%   0.00 ± 0%  -100.00% (p=0.000 n=10)
geomean       31.26                    ?                       ¹ ²
¹ summaries must be >0 to compute geomean
² ratios must be >0 to compute geomean
```

```lang=bash
          │   allocs/op   │ allocs/op   vs base                     │
Add/16       1.000 ± 0%     0.000 ± 0%  -100.00% (p=0.000 n=10)
Get/16       0.000 ± 0%     0.000 ± 0%         ~ (p=1.000 n=10) ¹
Delete/16    0.000 ± 0%     0.000 ± 0%         ~ (p=1.000 n=10) ¹
geomean                 ²               ?                       ² ³
```

## Contributing

Contributions are welcome! Please feel free to submit a pull request or open an issue.

### Development 

Comparing the benchmark performance stats:
```lang=bash
go install golang.org/x/perf/cmd/benchstat@latest
go test -benchmem -run=^$ -bench '^(Benchmark)' cimap -count=10 > bench-all.txt
grep 'Base-' bench-all.txt | sed 's|Base-||g' > bench-old.txt
grep 'CIMap-' bench-all.txt | sed 's|CIMap-||g' > bench-new.txt
benchstat bench-old.txt bench-new.txt

```

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.