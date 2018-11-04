package bwerr

import (
	"fmt"

	_ "github.com/baza-winner/bwcore/ansi/tags"
)

func ExampleError_1() {
	fmt.Printf(`%q`,
		From("error message"),
	)
	// Output: "\x1b[91;1mERR: \x1b[0merror message at \x1b[32;1mgithub.com/baza-winner/bwcore/bwerr.ExampleError_1\x1b[38;5;243m@\x1b[97;1mbwerr_test.go:11\x1b[0m"
}

// func ExampleDebug() {
// 	// Output: \x1b[32;1mgithub.com/baza-winner/bwcore/bwerror_test.ExampleDebug\x1b[38;5;243m@\x1b[97;1mbwerror_test.go:17: \x1b[38;5;201;1mvarA\x1b[0m: \x1b[0m\x1b[96;1m(struct { s string }){s:(string)string value}\x1b[0m\x1b[0m, \x1b[38;5;201;1mvarB\x1b[0m: \x1b[0m\x1b[96;1m(struct { i int }){i:(int)273}\x1b[0m\x1b[0m\x1b[0m
// 	varA := struct{ s string }{s: "string value"}
// 	varB := struct{ i int }{i: 273}
// 	Debug("varA", varA, "varB", varB) // s, _ := AnsiDebug("varA", varA, "varB", varB) // fmt.Printf("%q\n", s) //
// 	// Output: \x1b[32;1mgithub.com/baza-winner/bwcore/bwerror_test.ExampleDebug@bwerror_test.go:17: varA: (struct { s string }){s:(string)string value}, varB: (struct { i int }){i:(int)273}
// }

// func ExampleDebug2() {
// 	varA := struct{ s string }{s: "string value"}
// 	varB := struct{ i int }{i: 273}
// 	Debug("!HERE", "varA", varA, "varB", varB) // s, _ := AnsiDebug("!HERE", "varA", varA, "varB", varB) // fmt.Printf("%q", s) //
// 	// Output: "\x1b[93;1m!HERE\x1b[0m, \x1b[32;1mgithub.com/baza-winner/bwcore/bwerror_test.ExampleDebug2\x1b[38;5;243m@\x1b[97;1mbwerror_test.go:24: \x1b[38;5;201;1mvarA\x1b[0m: \x1b[0m\x1b[96;1m(struct { s string }){s:(string)string value}\x1b[0m\x1b[0m, \x1b[38;5;201;1mvarB\x1b[0m: \x1b[0m\x1b[96;1m(struct { i int }){i:(int)273}\x1b[0m\x1b[0m\x1b[0m"
// }
