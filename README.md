# go-string2eth

[![Tag](https://img.shields.io/github/tag/wealdtech/go-string2eth.svg)](https://github.com/wealdtech/go-string2eth/releases/)
[![License](https://img.shields.io/github/license/wealdtech/go-string2eth.svg)](LICENSE)
[![GoDoc](https://godoc.org/github.com/wealdtech/go-string2eth?status.svg)](https://godoc.org/github.com/wealdtech/go-string2eth)
[![Travis CI](https://img.shields.io/travis/wealdtech/go-string2eth.svg)](https://travis-ci.org/wealdtech/go-string2eth)
[![codecov.io](https://img.shields.io/codecov/c/github/wealdtech/go-string2eth.svg)](https://codecov.io/github/wealdtech/go-string2eth)

Go utility library to convert strings to Ether values and vice versa.

When converting strings to numeric values the process is case-insensitive and doesn't care about whitespace, so input values such as "0.1 Ether", "0.1Ether" and "0.1ether" would all result in the same result.  The standard unit denominations (Wei, Ether) are supported with or without SI prefixes (micro, milli, kilo, mega etc.), as are common names (Shannon, Babbage).

When converting numeric values to strings the user can select standard mode, in which case all values will be in units of Wei or Ether, or non-standard mode, in which case the full range of values will be used.

## Table of Contents

- [Install](#install)
- [Usage](#usage)
- [Maintainers](#maintainers)
- [Contribute](#contribute)
- [License](#license)

## Install

`go-string2eth` is a standard Go module which can be installed with:

```sh
go get github.com/wealdtech/go-string2eth
```

## Usage

`go-string2eth` converts from strings to Ether values and back again.

### Example

```go
package main

import (
	string2eth "github.com/wealdtech/go-string2eth"
)

func main() {

    // Convert a string value to a number of Wei
    value, err := string2eth.StringToWei("0.05 Ether")
    if err != nil {
        panic(err)
    }

    // Convert a number of Wei to a string value
    str := string2eth.WeiToString(value, true)

    fmt.Printf("0.05 Ether is %v Wei, is %s\n", value, str)
}
```

## Maintainers

Jim McDonald: [@mcdee](https://github.com/mcdee).

## Contribute

Contributions welcome. Please check out [the issues](https://github.com/wealdtech/go-string2eth/issues).

## License

[Apache-2.0](LICENSE) © 2019 Weald Technology Trading Ltd
