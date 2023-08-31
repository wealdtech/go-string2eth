// Copyright 2019, 2022 Weald Technology Trading Ltd
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

package string2eth_test

import (
	"errors"
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	string2eth "github.com/wealdtech/go-string2eth"
)

func TestWeiToStringWithSmallEtherDecimalValue(t *testing.T) {
	expected := "1.000000000000000001 Ether"
	wei, _ := new(big.Int).SetString("1000000000000000001", 10)
	result := string2eth.WeiToString(wei, true)
	assert.Equal(t, expected, result, "Did not receive expected result")
}

func TestRoundTripWithSmallValue(t *testing.T) {
	first := "1 Wei"
	second, err := string2eth.StringToWei(first)
	assert.Nil(t, err, "Failed to convert Ether to Wei")
	assert.Equal(t, second, big.NewInt(1), "Unexpected result converting Ether to Wei")
	third := string2eth.WeiToString(second, false)
	assert.Equal(t, first, third, "Did not receive expected result")
	fourth := string2eth.WeiToString(second, true)
	assert.Equal(t, first, fourth, "Did not receive expected result")
}

func TestRoundTripWithNormalValue(t *testing.T) {
	first := "1 Ether"
	second, err := string2eth.StringToWei(first)
	assert.Nil(t, err, "Failed to convert Ether to Wei")
	assert.Equal(t, second, big.NewInt(1000000000000000000), "Unexpected result converting Ether to Wei")
	third := string2eth.WeiToString(second, true)
	assert.Equal(t, first, third, "Did not receive expected result")
}

func ExampleUnitToMultiplier() {
	multiplier, err := string2eth.UnitToMultiplier("ether")
	if err != nil {
		return
	}
	fmt.Println(multiplier.Text(10))
	// Output: 1000000000000000000
}

func _bigInt(input string) *big.Int {
	res, _ := new(big.Int).SetString(input, 10)
	return res
}

func TestStringToWei(t *testing.T) {
	tests := []struct {
		input  string
		result *big.Int
		err    error
	}{
		{ // 0
			input: "",
			err:   errors.New("failed to parse empty value"),
		},
		{ // 1
			input:  "1",
			result: _bigInt("1"),
		},
		{ // 2
			input:  "123456789",
			result: _bigInt("123456789"),
		},
		{ // 3
			input:  "123456789 Wei",
			result: _bigInt("123456789"),
		},
		{ // 4
			input:  "1000000000000000000000",
			result: _bigInt("1000000000000000000000"),
		},
		{ // 5
			input:  "0.024ether",
			result: _bigInt("24000000000000000"),
		},
		{ // 6
			input:  "85748574 microether",
			result: _bigInt("85748574000000000000"),
		},
		{ // 7
			input:  "85748574 milliether",
			result: _bigInt("85748574000000000000000"),
		},
		{ // 8
			input:  "1 ether",
			result: _bigInt("1000000000000000000"),
		},
		{ // 9
			input:  "1 kiloether",
			result: _bigInt("1000000000000000000000"),
		},
		{ // 10
			input:  "1 megaether",
			result: _bigInt("1000000000000000000000000"),
		},
		{ // 11
			input:  "1 gigaether",
			result: _bigInt("1000000000000000000000000000"),
		},
		{ // 12
			input:  "5000 Teraether",
			result: _bigInt("5000000000000000000000000000000000"),
		},
		{ // 13
			input:  "0.123 kwei",
			result: _bigInt("123"),
		},
		{ // 14
			input:  "0.0001 kiloether",
			result: _bigInt("100000000000000000"),
		},
		{ // 15
			input:  ".0000001 megaether",
			result: _bigInt("100000000000000000"),
		},
		{ // 16
			input:  "1. Mwei",
			result: _bigInt("1000000"),
		},
		{ // 17
			input:  "21 Gwei",
			result: _bigInt("21000000000"),
		},
		{ // 18
			input:  "1000 ",
			result: _bigInt("1000"),
		},
		{ // 19
			input:  "1000000000000000000000 Wei",
			result: _bigInt("1000000000000000000000"),
		},
		{ // 20
			input:  "2megawei",
			result: _bigInt("2000000"),
		},
		{ // 21
			input:  "2.876543megawei",
			result: _bigInt("2876543"),
		},
		{ // 22
			input: "2.8765432megawei",
			err:   errors.New("value resulted in fractional number of Wei"),
		},
		{ // 23
			input:  "2 mega wei",
			result: _bigInt("2000000"),
		},
		{ // 24
			input:  "    2    mega   wei    ",
			result: _bigInt("2000000"),
		},
		{ // 25
			input: "1000 foo",
			err:   errors.New("failed to parse 1000 foo"),
		},
		{ // 26
			input:  "2megawei",
			result: _bigInt("2000000"),
		},
		{ // 27
			input: "1000.5 foo",
			err:   errors.New("failed to parse 1000.5 foo"),
		},
		{ // 28
			input: "onehundred ether",
			err:   errors.New("failed to parse  onehundredether"),
		},
		{ // 29
			input: "onehundred.5 ether",
			err:   errors.New("invalid format"),
		},
		{ // 30
			input:  "0",
			result: _bigInt("0"),
		},
		{ // 31
			input:  "0 Ether",
			result: _bigInt("0"),
		},
		{ // 32
			input: "10 wei wei wei",
			err:   errors.New("failed to parse 10 weiweiwei"),
		},
		{ // 33
			input: "0.1wei",
			err:   errors.New("value resulted in fractional number of Wei"),
		},
		{ // 34
			input: "-2 wei",
			err:   errors.New("value resulted in negative number of Wei"),
		},
		{ // 35
			input: "@",
			err:   errors.New("invalid format"),
		},
		{ // 36
			input:  "5 Shannon",
			result: _bigInt("5000000000"),
		},
		{ // 37
			input:  "1_000_000 Ether",
			result: _bigInt("1000000000000000000000000"),
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			result, err := string2eth.StringToWei(test.input)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error(), fmt.Sprintf("Incorrect error at test %d", i))
			} else {
				require.Nil(t, test.err, fmt.Sprintf("Unexpected error at test %d", i))
				assert.Equal(t, test.result, result, fmt.Sprintf("Incorrect value at test %d", i))
			}
		})
	}
}

func TestWeiToString(t *testing.T) {
	tests := []struct {
		input     *big.Int
		canonical bool
		result    string
		err       error
	}{
		{ // 0
			input:     _bigInt("1"),
			canonical: true,
			result:    "1 Wei",
		},
		{ // 1
			input:     _bigInt("2034"),
			canonical: true,
			result:    "2.034 KWei",
		},
		{ // 2
			input:     _bigInt("1234567890"),
			canonical: true,
			result:    "1.23456789 GWei",
		},
		{ // 3
			input:     _bigInt(""),
			canonical: true,
			result:    "0",
		},
		{ // 4
			input:     _bigInt("1000000000000000000"),
			canonical: true,
			result:    "1 Ether",
		},
		{ // 5
			input:     _bigInt("1000000000000000001"),
			canonical: true,
			result:    "1.000000000000000001 Ether",
		},
		{ // 6
			input:     _bigInt("1"),
			canonical: true,
			result:    "1 Wei",
		},
		{ // 7
			input:     _bigInt("999"),
			canonical: true,
			result:    "999 Wei",
		},
		{ // 8
			input:     _bigInt("1000"),
			canonical: true,
			result:    "1 KWei",
		},
		{ // 9
			input:     _bigInt("1001"),
			canonical: true,
			result:    "1.001 KWei",
		},
		{ // 10
			input:     _bigInt("999999"),
			canonical: true,
			result:    "999.999 KWei",
		},
		{ // 11
			input:     _bigInt("1000000"),
			canonical: true,
			result:    "1 MWei",
		},
		{ // 12
			input:     nil,
			canonical: true,
			result:    "0",
		},
		{ // 13
			input:     _bigInt("1000001"),
			canonical: true,
			result:    "1.000001 MWei",
		},
		{ // 14
			input:     _bigInt("999999999"),
			canonical: true,
			result:    "999.999999 MWei",
		},
		{ // 15
			input:     _bigInt("1000000000"),
			canonical: true,
			result:    "1 GWei",
		},
		{ // 16
			input:     _bigInt("1000000001"),
			canonical: true,
			result:    "1.000000001 GWei",
		},
		{ // 17
			input:     _bigInt("999999999999"),
			canonical: true,
			result:    "999.999999999 GWei",
		},
		{ // 18
			input:     _bigInt("1000000000000"),
			canonical: true,
			result:    "1000 GWei",
		},
		{ // 19
			input:     _bigInt("1000000000000"),
			canonical: false,
			result:    "1 Microether",
		},
		{ // 20
			input:     _bigInt("1000000000001"),
			canonical: true,
			result:    "1000.000000001 GWei",
		},
		{ // 21
			input:     _bigInt("1000000000001"),
			canonical: false,
			result:    "1.000000000001 Microether",
		},
		{ // 22
			input:     _bigInt("999999999999999"),
			canonical: true,
			result:    "999999.999999999 GWei",
		},
		{ // 23
			input:     _bigInt("999999999999999"),
			canonical: false,
			result:    "999.999999999999 Microether",
		},
		{ // 24
			input:     _bigInt("1000000000000000"),
			canonical: true,
			result:    "0.001 Ether",
		},
		{ // 25
			input:     _bigInt("1000000000000000"),
			canonical: false,
			result:    "1 Milliether",
		},
		{ // 26
			input:     _bigInt("1000000000000001"),
			canonical: true,
			result:    "0.001000000000000001 Ether",
		},
		{ // 27
			input:     _bigInt("1000000000000001"),
			canonical: false,
			result:    "1.000000000000001 Milliether",
		},
		{ // 28
			input:     _bigInt("999999999999999999"),
			canonical: true,
			result:    "0.999999999999999999 Ether",
		},
		{ // 29
			input:     _bigInt("999999999999999999"),
			canonical: false,
			result:    "999.999999999999999 Milliether",
		},
		{ // 30
			input:     _bigInt("1000000000000000000"),
			canonical: true,
			result:    "1 Ether",
		},
		{ // 31
			input:     _bigInt("1000000000000000000"),
			canonical: false,
			result:    "1 Ether",
		},
		{ // 32
			input:     _bigInt("1000000000000000001"),
			canonical: true,
			result:    "1.000000000000000001 Ether",
		},
		{ // 33
			input:     _bigInt("1000000000000000001"),
			canonical: false,
			result:    "1.000000000000000001 Ether",
		},
		{ // 34
			input:     _bigInt("999999999999999999999"),
			canonical: true,
			result:    "999.999999999999999999 Ether",
		},
		{ // 35
			input:     _bigInt("999999999999999999999"),
			canonical: false,
			result:    "999.999999999999999999 Ether",
		},
		{ // 36
			input:     _bigInt("1000000000000000000000"),
			canonical: true,
			result:    "1000 Ether",
		},
		{ // 37
			input:     _bigInt("1000000000000000000000"),
			canonical: false,
			result:    "1 Kiloether",
		},
		{ // 38
			input:     _bigInt("1000000000000000000001"),
			canonical: true,
			result:    "1000.000000000000000001 Ether",
		},
		{ // 39
			input:     _bigInt("1000000000000000000001"),
			canonical: false,
			result:    "1.000000000000000000001 Kiloether",
		},
		{ // 40
			input:     _bigInt("999999999999999999999999"),
			canonical: true,
			result:    "999999.999999999999999999 Ether",
		},
		{ // 41
			input:     _bigInt("999999999999999999999999"),
			canonical: false,
			result:    "999.999999999999999999999 Kiloether",
		},
		{ // 42
			input:     _bigInt("1000000000000000000000000"),
			canonical: true,
			result:    "1000000 Ether",
		},
		{ // 43
			input:     _bigInt("1000000000000000000000000"),
			canonical: false,
			result:    "1 Megaether",
		},
		{ // 44
			input:     _bigInt("1000000000000000000000001"),
			canonical: true,
			result:    "1000000.000000000000000001 Ether",
		},
		{ // 45
			input:     _bigInt("1000000000000000000000001"),
			canonical: false,
			result:    "1.000000000000000000000001 Megaether",
		},
		{ // 46
			input:     _bigInt("999999999999999999999999999"),
			canonical: true,
			result:    "999999999.999999999999999999 Ether",
		},
		{ // 47
			input:     _bigInt("999999999999999999999999999"),
			canonical: false,
			result:    "999.999999999999999999999999 Megaether",
		},
		{ // 48
			input:     _bigInt("1000000000000000000000000000"),
			canonical: true,
			result:    "1000000000 Ether",
		},
		{ // 49
			input:     _bigInt("1000000000000000000000000000"),
			canonical: false,
			result:    "1 Gigaether",
		},
		{ // 50
			input:     _bigInt("1000000000000000000000000001"),
			canonical: true,
			result:    "1000000000.000000000000000001 Ether",
		},
		{ // 51
			input:     _bigInt("1000000000000000000000000001"),
			canonical: false,
			result:    "1.000000000000000000000000001 Gigaether",
		},
		{ // 52
			input:     _bigInt("999999999999999999999999999999"),
			canonical: true,
			result:    "999999999999.999999999999999999 Ether",
		},
		{ // 53
			input:     _bigInt("999999999999999999999999999999"),
			canonical: false,
			result:    "999.999999999999999999999999999 Gigaether",
		},
		{ // 54
			input:     _bigInt("1000000000000000000000000000000"),
			canonical: true,
			result:    "1000000000000 Ether",
		},
		{ // 55
			input:     _bigInt("1000000000000000000000000000000"),
			canonical: false,
			result:    "1 Teraether",
		},
		{ // 56
			input:     _bigInt("1000000000000000000000000000001"),
			canonical: true,
			result:    "1000000000000.000000000000000001 Ether",
		},
		{ // 57
			input:     _bigInt("1000000000000000000000000000001"),
			canonical: false,
			result:    "1.000000000000000000000000000001 Teraether",
		},
		{ // 58
			input:     _bigInt("999999999999999999999999999999999"),
			canonical: true,
			result:    "999999999999999.999999999999999999 Ether",
		},
		{ // 59
			input:     _bigInt("999999999999999999999999999999999"),
			canonical: false,
			result:    "999.999999999999999999999999999999 Teraether",
		},
		{ // 60
			input:     _bigInt("1000000000000000000000000000000000"),
			canonical: true,
			result:    "1000000000000000 Ether",
		},
		{ // 61
			input:     _bigInt("1000000000000000000000000000000000"),
			canonical: false,
			result:    "overflow",
		},
		{ // 62
			input:     _bigInt(""),
			canonical: true,
			result:    "0",
		},
		{ // 63
			input:     _bigInt("0"),
			canonical: true,
			result:    "0",
		},
		{ // 64
			input:     _bigInt("999999999999"),
			canonical: true,
			result:    "999.999999999 GWei",
		},
		{ // 65
			input:     _bigInt("999999999999999"),
			canonical: true,
			result:    "999999.999999999 GWei",
		},
		{ // 66
			input:     _bigInt("1000000000000000"),
			canonical: true,
			result:    "0.001 Ether",
		},
	}

	for i, test := range tests {
		result := string2eth.WeiToString(test.input, test.canonical)
		assert.Equal(t, test.result, result, fmt.Sprintf("Incorrect value at test %d", i))
	}
}

func TestGWeiToString(t *testing.T) {
	tests := []struct {
		name      string
		input     uint64
		canonical bool
		result    string
	}{
		{
			name:      "Zero",
			input:     0,
			canonical: true,
			result:    "0",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := string2eth.GWeiToString(test.input, test.canonical)
			require.Equal(t, test.result, result)
		})
	}
}

func TestStringToGWei(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		result uint64
		err    error
	}{
		{
			name:   "Zero",
			input:  "0",
			result: 0,
		},
		{
			name:  "Invalid",
			input: "@",
			err:   errors.New("invalid format"),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := string2eth.StringToGWei(test.input)
			if test.err != nil {
				require.NotNil(t, err)
				require.Equal(t, test.err.Error(), err.Error())
			} else {
				require.Nil(t, err)
				require.Equal(t, test.result, result)
			}
		})
	}
}

func TestWeiToGWeiString(t *testing.T) {
	tests := []struct {
		name   string
		input  *big.Int
		result string
	}{
		{
			name:   "Nil",
			result: "0",
		},
		{
			name:   "1",
			input:  big.NewInt(1),
			result: "0.000000001 GWei",
		},
		{
			name:   "999",
			input:  big.NewInt(999),
			result: "0.000000999 GWei",
		},
		{
			name:   "1000",
			input:  big.NewInt(1000),
			result: "0.000001 GWei",
		},
		{
			name:   "999999",
			input:  big.NewInt(999999),
			result: "0.000999999 GWei",
		},
		{
			name:   "1000000",
			input:  big.NewInt(1000000),
			result: "0.001 GWei",
		},
		{
			name:   "999999999",
			input:  big.NewInt(999999999),
			result: "0.999999999 GWei",
		},
		{
			name:   "100000000",
			input:  big.NewInt(1000000000),
			result: "1 GWei",
		},
		{
			name:   "999000000000",
			input:  big.NewInt(999000000000),
			result: "999 GWei",
		},
		{
			name:   "999000050000",
			input:  big.NewInt(999000050000),
			result: "999.00005 GWei",
		},
		{
			name:   "1000000000000",
			input:  big.NewInt(10000000000000),
			result: "10000 GWei",
		},
		{
			name:   "1000010000000",
			input:  big.NewInt(10000100000000),
			result: "10000.1 GWei",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := string2eth.WeiToGWeiString(test.input)
			require.Equal(t, test.result, result)
		})
	}
}
