package bwjson

import "fmt"

func ExamplePrettyJson_1() {
  fmt.Println(PrettyJson(true))
  // Output: true
}

func ExamplePrettyJson_2() {
  fmt.Println(PrettyJson(100))
  // Output: 100
}

func ExamplePrettyJson_3() {
  fmt.Println(PrettyJson(`string`))
  // Output: "string"
}

func ExamplePrettyJson_4() {
  fmt.Println(PrettyJson([]interface{}{ false, 273, `something`}))
  // Output: [
  //   false,
  //   273,
  //   "something"
  // ]
}

func ExamplePrettyJson_5() {
  fmt.Println(PrettyJson(
    map[string]interface{}{
      "bool": true,
      "number": 273,
      "string": `something`,
      "array": []interface{}{ "one", true, 3 },
      "map": map[string]interface{}{
        "one": 1,
        "two": true,
        "three": "three",
      },
    },
  ))
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