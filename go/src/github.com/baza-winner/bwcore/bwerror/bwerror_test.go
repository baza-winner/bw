package bwerror

import (
	"fmt"
)

func ExampleError_1() {
	fmt.Printf(`%+q`, Error("error message"))
	// Output: "\x1b[0m\x1b[31m\x1b[1mERR:\x1b[0m error message\x1b[0m"
}

func ExampleError_2() {
	fmt.Printf(`%+q`, Error("error message with <ansiOutline>ansi<ansi> formatting"))
	// Output: "\x1b[0m\x1b[31m\x1b[1mERR:\x1b[0m error message with \x1b[38;5;201m\x1b[1mansi\x1b[0m formatting\x1b[0m"
}

func ExampleError_3() {
	fmt.Printf(`%+q`, Error("error message: <ansiOutline>string <ansiPrimary>%s<ansi>, <ansiOutline>number <ansiPrimary>%d<ansi>, <ansiOutline>value <ansiSecondary>%+v", "string value", 273, map[string]interface{}{"some": "thing"}))
	// Output: "\x1b[0m\x1b[31m\x1b[1mERR:\x1b[0m error message: \x1b[38;5;201m\x1b[1mstring \x1b[36m\x1b[1mstring value\x1b[0m, \x1b[38;5;201m\x1b[1mnumber \x1b[36m\x1b[1m273\x1b[0m, \x1b[38;5;201m\x1b[1mvalue \x1b[36m\x1b[22mmap[some:thing]\x1b[0m"
}
