package ansi_test

import (
	"fmt"
	"testing"

	"github.com/baza-winner/bwcore/ansi"
)

func TestString(t *testing.T) {
	ansi.MustAddTag("ansiOutline",
		ansi.SGRCodeOfColor256(ansi.Color256{Code: 201}),
	)
	ansi.MustAddTag("ansiErr",
		ansi.MustSGRCodeOfColor8(ansi.Color8{Color: ansi.SGRColorRed, Bright: true}),
	)
	ansi.MustAddTag("ansiCmd",
		ansi.MustSGRCodeOfColor8(ansi.Color8{Color: ansi.SGRColorWhite, Bright: true}),
		ansi.MustSGRCodeOfCmd(ansi.SGRCmdBold),
	)
	ansi.MustAddTag("ansiOK",
		ansi.MustSGRCodeOfColor8(ansi.Color8{Color: ansi.SGRColorGreen, Bright: true}),
	)
	tests := []struct {
		In  ansi.A
		Out string
	}{
		{
			In:  ansi.A{ansi.MustTag("ansiErr"), "ERR: <ansiCmd>some<ansi> expects arg"},
			Out: "\x1b[91mERR: \x1b[97;1msome\x1b[91m expects arg\x1b[0m",
		},
		{
			In:  ansi.A{S: "some <ansiOutline>thing<ansi> good"},
			Out: "some \x1b[38;5;201mthing\x1b[0m good\x1b[0m",
		},
	}
	for _, test := range tests {
		got := ansi.String(test.In)
		if got != test.Out {
			t.Errorf(
				// "From(%#v)\n    => %q\n, want %q", test.In, got.S, test.Out.S,
				ansi.String(
					ansi.A{S: fmt.Sprintf("From(%#v)\n    => <ansiErr>%q<ansi>\n, want <ansiOK>%q", test.In, got, test.Out)},
				),
			)
		}
	}
}

func TestConcat(t *testing.T) {
	tests := []struct {
		In  []string
		Out string
	}{
		{
			In: []string{
				"",
				ansi.String(ansi.A{S: "<ansiCmd>some<ansi> expects arg"}),
			},
			Out: "\x1b[97;1msome\x1b[0m expects arg\x1b[0m",
		},
		{
			In: []string{
				ansi.String(ansi.A{Default: ansi.MustTag("ansiErr"), S: "ERR: "}),
				ansi.String(ansi.A{S: "<ansiCmd>some<ansi> expects arg"}),
			},
			Out: "\x1b[91mERR: \x1b[97;1msome\x1b[0m expects arg\x1b[0m",
		},
	}
	for _, test := range tests {
		got := ansi.Concat(test.In...)
		if got != test.Out {
			t.Errorf(
				ansi.String(
					ansi.A{S: fmt.Sprintf("Concat(%#v)\n    => <ansiErr>%q<ansi>\n, want <ansiOK>%q", test.In, got, test.Out)},
				),
			)
		}
	}
}

func ExampleAnsi() {
	fmt.Printf("%q",
		ansi.String(ansi.A{S: "some <ansiOutline>thing<ansi> good"}),
	)
	// Output:
	// "some \x1b[38;5;201mthing\x1b[0m good\x1b[0m"
}

func ExampleAnsi2() {
	fmt.Printf("%q",
		ansi.String(ansi.A{S: fmt.Sprintf("some <ansiOutline>%s<ansi> good", "thing")}),
	)
	// Output:
	// "some \x1b[38;5;201mthing\x1b[0m good\x1b[0m"
}

func ExampleAnsi3() {
	fmt.Printf("%q",
		ansi.String(ansi.A{S: fmt.Sprintf("some <ansiOutline>%s<ansi> good", "thing")}),
	)
	// Output:
	// "some \x1b[38;5;201mthing\x1b[0m good\x1b[0m"
}

func ExampleAnsi4() {
	fmt.Printf("%q",
		ansi.String(ansi.A{ansi.MustTag("ansiErr"), fmt.Sprintf("ERR: <ansiCmd>%s<ansi> expects arg\n\n", "some")}),
	)
	// Output:
	// "\x1b[91mERR: \x1b[97;1msome\x1b[91m expects arg\n\n\x1b[0m"
}
