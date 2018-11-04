package bwjson_test

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/baza-winner/bwcore/bwjson"
)

func ExamplePretty_1() {
	fmt.Println(
		bwjson.Pretty(true),
	)
	// Output: true
}

func ExamplePretty_2() {
	fmt.Println(
		bwjson.Pretty(100),
	)
	// Output: 100
}

func ExamplePretty_3() {
	fmt.Println(
		bwjson.Pretty(`string`),
	)
	// Output: "string"
}

func ExamplePretty_4() {
	fmt.Println(
		bwjson.Pretty([]interface{}{false, 273, `something`}),
	)
	// Output: [
	//   false,
	//   273,
	//   "something"
	// ]
}

func ExamplePretty_5() {
	fmt.Println(
		bwjson.Pretty(
			map[string]interface{}{
				"bool":   true,
				"number": 273,
				"string": `something`,
				"array":  []interface{}{"one", true, 3},
				"map": map[string]interface{}{
					"one":   1,
					"two":   true,
					"three": "three",
				},
			},
		),
	)
	// Output:
	// {
	//   "array": [
	//     "one",
	//     true,
	//     3
	//   ],
	//   "bool": true,
	//   "map": {
	//     "one": 1,
	//     "three": "three",
	//     "two": true
	//   },
	//   "number": 273,
	//   "string": "something"
	// }
}

type EnumType uint8

const (
	Enum0 EnumType = iota
	Enum1
	Enum2
)

const _EnumType_name = "Enum0Enum1Enum2"

var _EnumType_index = [...]uint8{0, 5, 10, 15}

func (i EnumType) String() string {
	if i >= EnumType(len(_EnumType_index)-1) {
		return "EnumType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _EnumType_name[_EnumType_index[i]:_EnumType_index[i+1]]
}

func (v EnumType) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.String())
}

type someStruct struct {
	boolField bool
	numField  int
	strField  string
	enumField EnumType
}

func (v someStruct) MarshalJSON() ([]byte, error) {
	result := map[string]interface{}{}
	result["boolField"] = v.boolField
	result["numField"] = v.numField
	result["strField"] = v.strField
	result["enumField"] = v.enumField
	return json.Marshal(result)
}

func ExamplePretty_6() {
	v := someStruct{
		boolField: true,
		numField:  273,
		strField:  "something",
		enumField: Enum2,
	}
	fmt.Println(
		bwjson.Pretty(v),
	)
	// Output:
	// {
	//   "boolField": true,
	//   "enumField": "Enum2",
	//   "numField": 273,
	//   "strField": "something"
	// }
}
