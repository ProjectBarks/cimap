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

> :warning: **Note**: Code performs best when there are string allocations due to
>			`strings.ToLower` / `strings.ToUpper`. This would occur if the strings in the 
> 			map either contain Unicode or mismatching case-characters. If you can gaurantee 
> 			the inputs provided will only be one character case consider not using
> 			case insensitivity at all. 

```bash
                         │         sec/op          │    sec/op     vs base                │
Add/ASCII_Mismatch/16                  70.66n ± 3%    54.02n ± 2%  -23.55% (p=0.000 n=25)
Add/ASCII_Match/16                     19.35n ± 3%    35.19n ± 6%  +81.86% (p=0.000 n=25)
Add/Unicode/16                         148.8n ± 5%    127.2n ± 3%  -14.52% (p=0.000 n=25)
Get/ASCII_Mismatch/16                  115.0n ± 5%    101.2n ± 6%  -12.00% (p=0.000 n=25)
Get/ASCII_Match/16                     53.59n ± 3%   104.00n ± 4%  +94.07% (p=0.000 n=25)
Get/Unicode/16                         578.0n ± 3%    652.9n ± 7%  +12.96% (p=0.000 n=25)
Delete/ASCII_Mismatch/16               80.17n ± 4%    50.76n ± 4%  -36.68% (p=0.000 n=25)
Delete/ASCII_Match/16                  30.87n ± 3%    50.80n ± 4%  +64.56% (p=0.000 n=25)
Delete/Unicode/16                      551.9n ± 3%    475.2n ± 3%  -13.90% (p=0.001 n=25)
```

```bash
                         │          B/op           │    B/op      vs base                     │
Add/ASCII_Mismatch/16                 92.00 ± 1%      0.00 ±  0%  -100.00% (p=0.000 n=25)
Add/ASCII_Match/16                    80.00 ± 3%      0.00 ±  0%  -100.00% (p=0.000 n=25)
Add/Unicode/16                        117.0 ± 3%       0.0 ±  0%  -100.00% (p=0.000 n=25)
Get/ASCII_Mismatch/16                 33.00 ± 0%      0.00 ±  0%  -100.00% (p=0.000 n=25)
Get/ASCII_Match/16                    0.000 ± 0%     0.000 ±  0%         ~ (p=1.000 n=25) ¹
Get/Unicode/16                        112.0 ± 0%       0.0 ±  0%  -100.00% (p=0.000 n=25)
Delete/ASCII_Mismatch/16              33.00 ± 0%      0.00 ±  0%  -100.00% (p=0.000 n=25)
Delete/ASCII_Match/16                 0.000 ± 0%     0.000 ±  0%         ~ (p=1.000 n=25) ¹
Delete/Unicode/16                   120.000 ± 1%     3.000 ± 33%   -97.50% (p=0.000 n=25)
```

```bash
                         │        allocs/op        │ allocs/op   vs base                     │
Add/ASCII_Mismatch/16                 1.000 ± 0%     0.000 ± 0%  -100.00% (p=0.000 n=25)
Add/ASCII_Match/16                    0.000 ± 0%     0.000 ± 0%         ~ (p=1.000 n=25) ¹
Add/Unicode/16                        1.000 ± 0%     0.000 ± 0%  -100.00% (p=0.000 n=25)
Get/ASCII_Mismatch/16                 1.000 ± 0%     0.000 ± 0%  -100.00% (p=0.000 n=25)
Get/ASCII_Match/16                    0.000 ± 0%     0.000 ± 0%         ~ (p=1.000 n=25) ¹
Get/Unicode/16                        1.000 ± 0%     0.000 ± 0%  -100.00% (p=0.000 n=25)
Delete/ASCII_Mismatch/16              1.000 ± 0%     0.000 ± 0%  -100.00% (p=0.000 n=25)
Delete/ASCII_Match/16                 0.000 ± 0%     0.000 ± 0%         ~ (p=1.000 n=25) ¹
Delete/Unicode/16                     2.000 ± 0%     0.000 ± 0%  -100.00% (p=0.000 n=25)
```

## Contributing

Contributions are welcome! Please feel free to submit a pull request or open an issue.

### Development

Comparing the benchmark performance stats.

#### Setup 

```bash
go install golang.org/x/perf/cmd/benchstat@latest
```

#### Running benchmarks

```bash
go test -benchmem -run=^$ -bench '^(Benchmark)' github.com/projectbarks/cimap -count=10 > benchmark/bench-all.txt
grep 'Base-' benchmark/bench-all.txt | sed 's|Base-||g' > benchmark/bench-old.txt
grep 'CIMap-' benchmark/bench-all.txt | sed 's|CIMap-||g' > benchmark/bench-new.txt
benchstat benchmark/bench-old.txt benchmark/bench-new.txt

```

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
