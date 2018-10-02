package bwos

import (
	"fmt"
	"os"
)

func ExampleShortenFileSpec() {
	fmt.Printf(`%q`, ShortenFileSpec(os.Getenv(`HOME`)+`/bw`))
	// Output: "~/bw"
}

func ExampleShortenFileSpec_2() {
	fmt.Printf(`%q`, ShortenFileSpec(`/lib/bw`))
	// Output: "/lib/bw"
}
