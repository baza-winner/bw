package ansi_test

import (
	"fmt"
	"testing"

	"github.com/baza-winner/bwcore/ansi"
	_ "github.com/baza-winner/bwcore/ansi/tags"
)

func TestString(t *testing.T) {
	tests := []struct {
		in  ansi.A
		out string
	}{
		{
			in:  ansi.A{ansi.MustTag("ansiErr"), "ERR: <ansiPath>some<ansi> expects arg"},
			out: "\x1b[91;1mERR: \x1b[38;5;252;1msome\x1b[91;1m expects arg\x1b[0m",
		},
		{
			in:  ansi.A{ansi.MustTag("ansiErr"), ""},
			out: "",
		},
		{
			in:  ansi.A{S: ""},
			out: "",
		},
		{
			in:  ansi.A{S: "some"},
			out: "some",
		},
		{
			in:  ansi.A{S: "some <ansiVar>thing<ansi> good"},
			out: "some \x1b[38;5;201;1mthing\x1b[0m good\x1b[0m",
		},
	}
	for _, test := range tests {
		got := ansi.StringA(test.in)
		if got != test.out {
			t.Errorf(ansi.String(
				fmt.Sprintf("From(%#v)\n    => <ansiErr>%q<ansi>\n, want <ansiOK>%q", test.in, got, test.out),
			))
		}
	}
}

// func TestConcat(t *testing.T) {
// 	tests := []struct {
// 		in  []string
// 		out string
// 	}{
// 		{
// 			in: []string{
// 				"",
// 				ansi.String("<ansiPath>some<ansi> expects arg"),
// 			},
// 			out: "\x1b[38;5;252;1msome\x1b[0m expects arg\x1b[0m",
// 		},
// 		{
// 			in: []string{
// 				ansi.StringA(ansi.A{Default: ansi.MustTag("ansiErr"), S: "ERR: "}),
// 				ansi.String("<ansiPath>some<ansi> expects arg"),
// 			},
// 			out: "\x1b[91;1mERR: \x1b[38;5;252;1msome\x1b[0m expects arg\x1b[0m",
// 		},
// 	}
// 	for _, test := range tests {
// 		got := ansi.Concat(test.in...)
// 		if got != test.out {
// 			t.Errorf(ansi.String(
// 				fmt.Sprintf("Concat(%#v)\n    => <ansiErr>%q<ansi>\n, want <ansiOK>%q", test.in, got, test.out),
// 			))
// 		}
// 	}
// }

func ExampleAnsi() {
	fmt.Printf("%q",
		ansi.String("some <ansiVar>thing<ansi> good"),
	)
	// Output:
	// "some \x1b[38;5;201;1mthing\x1b[0m good\x1b[0m"
}

func ExampleAnsi2() {
	fmt.Printf("%q",
		ansi.String(fmt.Sprintf("some <ansiVar>%s<ansi> good", "thing")),
	)
	// Output:
	// "some \x1b[38;5;201;1mthing\x1b[0m good\x1b[0m"
}

func ExampleAnsi3() {
	fmt.Printf("%q",
		ansi.String(fmt.Sprintf("some <ansiVar>%s<ansi> good", "thing")),
	)
	// Output:
	// "some \x1b[38;5;201;1mthing\x1b[0m good\x1b[0m"
}

func ExampleAnsi4() {
	fmt.Printf("%q",
		ansi.StringA(ansi.A{ansi.MustTag("ansiErr"), fmt.Sprintf("ERR: <ansiPath>%s<ansi> expects arg\n\n", "some")}),
	)
	// Output:
	// "\x1b[91;1mERR: \x1b[38;5;252;1msome\x1b[91;1m expects arg\n\n\x1b[0m"
}
