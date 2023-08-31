// Copyright 2019 - 2023 Weald Technology Trading Limited.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// string2eth provides methods for converting between number of Wei and a string
// represetation of the same, with the latter allowing all commonly-used
// representations as inputs.  String outputs are provided in a sensible format
// given the value.
package string2eth

import (
	"errors"
	"fmt"
	"math/big"
	"regexp"
	"strings"
)

var (
	ErrEmptyValue    = errors.New("failed to parse empty value")
	ErrInvalidFormat = errors.New("invalid format")
	ErrNegative      = errors.New("value resulted in negative number of Wei")
	ErrFractional    = errors.New("value resulted in fractional number of Wei")
	ErrUnknownUnit   = errors.New("unknown unit")
	ErrParseFailure  = errors.New("failed to parse")
)

// StringToWei turns a string in to number of Wei.
// The string can be a simple number of Wei, e.g. "1000000000000000" or it can
// be a number followed by a unit, e.g. "10 ether".  Unit names are
// case-insensitive, and can be either given names (e.g. "finney") or metric
// names (e.g. "mlliether").
// Note that this function expects use of the period as the decimal separator.
func StringToWei(input string) (*big.Int, error) {
	if input == "" {
		return nil, ErrEmptyValue
	}

	// Remove unused runes that may be in an input string.
	input = strings.ReplaceAll(input, " ", "")
	input = strings.ReplaceAll(input, "_", "")

	var result big.Int
	// Separate the number from the unit (if any)
	re := regexp.MustCompile(`^(-?[0-9]*(?:\.[0-9]*)?)([A-Za-z]+)?$`)
	subMatches := re.FindAllStringSubmatch(input, -1)
	var units string
	if len(subMatches) != 1 {
		return nil, ErrInvalidFormat
	}
	units = subMatches[0][2]
	if strings.Contains(subMatches[0][1], ".") {
		err := decimalStringToWei(subMatches[0][1], units, &result)
		if err != nil {
			return nil, err
		}
	} else {
		err := integerStringToWei(subMatches[0][1], units, &result)
		if err != nil {
			return nil, err
		}
	}

	// Ensure we don't have a negative number.
	if result.Cmp(new(big.Int)) < 0 {
		return nil, ErrNegative
	}

	return &result, nil
}

// StringToGWei turns a string in to number of GWei.
// See StringToWei for details.
// Any part of the value below 1GWei in denomination is lost.
func StringToGWei(input string) (uint64, error) {
	wei, err := StringToWei(input)
	if err != nil {
		return 0, err
	}

	return wei.Div(wei, billion).Uint64(), nil
}

// Used in WeiToString.
var (
	zero     = big.NewInt(0)
	thousand = big.NewInt(1000)
)

// Used in GWeiToString.
var billion = big.NewInt(1000000000)

// GWeiToString turns a number of GWei in to a string.
// See WeiToString for details.
func GWeiToString(input uint64, standard bool) string {
	return WeiToString(new(big.Int).Mul(new(big.Int).SetUint64(input), billion), standard)
}

// WeiToGWeiString turns a number of wei in to a Gwei string.
func WeiToGWeiString(input *big.Int) string {
	if input == nil {
		return "0"
	}

	intValue := new(big.Int).Div(input, billion)
	decValue := new(big.Int).Sub(input, new(big.Int).Mul(intValue, billion))

	// Return our value.
	if decValue.Cmp(zero) == 0 {
		return fmt.Sprintf("%s GWei", intValue)
	}
	decStr := strings.TrimRight(fmt.Sprintf("%09d", decValue.Int64()), "0")

	return fmt.Sprintf("%s.%s GWei", intValue, decStr)
}

// WeiToString turns a number of Wei in to a string.
// If the 'standard' argument is true then this will display the value
// in either (KMG)Wei or Ether only.
func WeiToString(input *big.Int, standard bool) string {
	if input == nil {
		return "0"
	}

	// Take a copy of the input so that we can mutate it.
	value := new(big.Int).Set(input)

	// Short circuit on 0.
	if value.Cmp(zero) == 0 {
		return "0"
	}

	// Step 1: work out simple units, keeping value as a whole number.
	value, unitPos := weiToStringStep1(value)

	// Step 2: move value to a fraction if sensible.
	outputValue, unitPos, desiredUnitPos, decimalPlace := weiToStringStep2(value, unitPos, standard)

	// Step 3: generate output.
	outputValue, unitPos = weiToStringStep3(outputValue, unitPos, desiredUnitPos, decimalPlace)

	if unitPos >= len(metricUnits) {
		return "overflow"
	}

	// Return our value.
	return fmt.Sprintf("%s %s", outputValue, metricUnits[unitPos])
}

// weiToStringStep1 steps the value down by thousands to obtain a smaller value
// with unit reference.
func weiToStringStep1(value *big.Int) (*big.Int, int) {
	unitPos := 0
	modInt := new(big.Int).Set(value)
	for value.Cmp(thousand) >= 0 && modInt.Mod(value, thousand).Cmp(zero) == 0 {
		unitPos++
		value = value.Div(value, thousand)
		modInt = modInt.Set(value)
	}

	return value, unitPos
}

// weiToStringStep2 starts to turn a value into a string, handling the case where
// the resultant output may be a decial.
func weiToStringStep2(value *big.Int, unitPos int, standard bool) (string, int, int, int) {
	// Because of the inaccuracy of floating point we use string manipulation
	// to place the decimal in the correct position.
	outputValue := value.Text(10)

	desiredUnitPos := unitPos
	if len(outputValue) > 3 {
		desiredUnitPos += len(outputValue) / 3
		if len(outputValue)%3 == 0 {
			desiredUnitPos--
		}
	}
	decimalPlace := len(outputValue)
	if desiredUnitPos > 3 && standard {
		// Because Gwei covers a large range allow anything up to 0.001 ETH
		// to display as Gwei.
		if desiredUnitPos == 4 {
			desiredUnitPos = 3
		} else {
			desiredUnitPos = 6
		}
	}
	for unitPos < desiredUnitPos {
		decimalPlace -= 3
		unitPos++
	}

	return outputValue, unitPos, desiredUnitPos, decimalPlace
}

// weiToStringStep3 finishes generation of the output value, ensuring the appropriate
// number of 0s and tidying up to provide a presentable result.
func weiToStringStep3(outputValue string, unitPos int, desiredUnitPos int, decimalPlace int) (string, int) {
	for unitPos > desiredUnitPos {
		outputValue += strings.Repeat("0", 3)
		decimalPlace += 3
		unitPos--
	}
	if decimalPlace <= 0 {
		outputValue = "0." + strings.Repeat("0", 0-decimalPlace) + outputValue
	} else if decimalPlace < len(outputValue) {
		outputValue = outputValue[:decimalPlace] + "." + outputValue[decimalPlace:]
	}

	// Trim trailing zeros if this is a decimal.
	if strings.Contains(outputValue, ".") {
		outputValue = strings.TrimRight(outputValue, "0")
	}

	return outputValue, unitPos
}

func decimalStringToWei(amount string, unit string, result *big.Int) error {
	// Because floating point maths is not accurate we need to break potentially
	// large decimal fractions in to two separate pieces: the integer part and the
	// decimal part.
	parts := strings.Split(amount, ".")

	// The value for the integer part of the number is easy.
	if parts[0] != "" {
		err := integerStringToWei(parts[0], unit, result)
		if err != nil {
			return fmt.Errorf("%w %s %s", ErrParseFailure, amount, unit)
		}
	}

	// The value for the decimal part of the number is harder.  We left-shift it
	// so that we end up multiplying two integers rather than two floats, as the
	// latter is unreliable.

	// Obtain multiplier.
	// This will never fail because it is already called above in integerStringToWei().
	multiplier, _ := UnitToMultiplier(unit)

	// Trim trailing 0s.
	trimmedDecimal := strings.TrimRight(parts[1], "0")
	if len(trimmedDecimal) == 0 {
		// Nothing more to do.
		return nil
	}
	var decVal big.Int
	decVal.SetString(trimmedDecimal, 10)

	// Divide multiplier by 10^len(trimmed decimal) to obtain sane value.
	div := big.NewInt(10)
	for i := 0; i < len(trimmedDecimal); i++ {
		multiplier.Div(multiplier, div)
	}

	// Ensure we don't have a fractional number of Wei.
	if multiplier.Sign() == 0 {
		return ErrFractional
	}

	var decResult big.Int
	decResult.Mul(multiplier, &decVal)

	// Add it to the integer result.
	result.Add(result, &decResult)

	return nil
}

func integerStringToWei(amount string, unit string, result *big.Int) error {
	// Obtain number.
	number := new(big.Int)
	_, success := number.SetString(amount, 10)
	if !success {
		return fmt.Errorf("%w %s %s", ErrParseFailure, amount, unit)
	}

	// Obtain multiplier.
	multiplier, err := UnitToMultiplier(unit)
	if err != nil {
		return fmt.Errorf("%w %s %s", ErrParseFailure, amount, unit)
	}

	result.Mul(number, multiplier)

	return nil
}

// Metric units.
var metricUnits = [...]string{
	"Wei",
	"KWei",
	"MWei",
	"GWei",
	"Microether",
	"Milliether",
	"Ether",
	"Kiloether",
	"Megaether",
	"Gigaether",
	"Teraether",
}

// UnitToMultiplier takes the name of an Ethereum unit and returns a multiplier.
//
//nolint:cyclop
func UnitToMultiplier(unit string) (*big.Int, error) {
	result := big.NewInt(0)
	switch strings.ToLower(unit) {
	case "", "wei":
		result.SetString("1", 10)
	case "ada", "kwei", "kilowei":
		result.SetString("1000", 10)
	case "babbage", "mwei", "megawei":
		result.SetString("1000000", 10)
	case "shannon", "gwei", "gigawei":
		result.SetString("1000000000", 10)
	case "szazbo", "micro", "microether":
		result.SetString("1000000000000", 10)
	case "finney", "milli", "milliether":
		result.SetString("1000000000000000", 10)
	case "eth", "ether":
		result.SetString("1000000000000000000", 10)
	case "einstein", "kilo", "kiloether":
		result.SetString("1000000000000000000000", 10)
	case "mega", "megaether":
		result.SetString("1000000000000000000000000", 10)
	case "giga", "gigaether":
		result.SetString("1000000000000000000000000000", 10)
	case "tera", "teraether":
		result.SetString("1000000000000000000000000000000", 10)
	default:
		return nil, fmt.Errorf("%w %s", ErrUnknownUnit, unit)
	}

	return result, nil
}
