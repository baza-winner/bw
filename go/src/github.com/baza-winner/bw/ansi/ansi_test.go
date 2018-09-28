package ansi

import (
	"fmt"
	"testing"
)

func TestAnsi(t *testing.T) {
	cases := []struct {
		defaultAnsiName string
		source          string
		result          string
	}{
		{
			defaultAnsiName: "Err",
			source:          "ERR: <ansiCmd>some<ansi> expects arg",
			result:          "\x1b[31m\x1b[1mERR: \x1b[97m\x1b[1msome\x1b[31m\x1b[1m expects arg\x1b[0m",
		},
		{
			defaultAnsiName: "Err",
			source:          "ERR: <ansiCmd>some<ansi> expects arg\n\n",
			result:          "\x1b[31m\x1b[1mERR: \x1b[97m\x1b[1msome\x1b[31m\x1b[1m expects arg\x1b[0m\n\n",
		},
		{
			defaultAnsiName: "",
			source:          "some <ansiOutline>thing<ansi> good",
			result:          "\x1b[0msome \x1b[38;5;201m\x1b[1mthing\x1b[0m good\x1b[0m",
		},
	}
	for _, c := range cases {
		got := Ansi(c.defaultAnsiName, c.source)
		if got != c.result {
			t.Errorf(Ansi("", "Ansi(%q, %q)\n    => <ansiErr>%q<ansi>\n, want <ansiOK>%q"), c.defaultAnsiName, c.source, got, c.result)
		}
	}
}

func ExampleAnsi() {
	fmt.Printf(`%q`, Ansi("Err", "ERR: <ansiCmd>some<ansi> expects arg\n\n"))
	// Output: "\x1b[31m\x1b[1mERR: \x1b[97m\x1b[1msome\x1b[31m\x1b[1m expects arg\x1b[0m\n\n"
}
