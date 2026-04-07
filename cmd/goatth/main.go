// Command goatth extracts embedded GoATTH assets to disk.
//
// Usage:
//
//	go tool goatth -out=css/goatth-base.css
//	go run github.com/guilycst/GoATTH-penguinui/cmd/goatth@latest -out=goatth.css
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/guilycst/GoATTH-penguinui/assets"
)

func main() {
	out := flag.String("out", "goatth-base.css", "output path for extracted CSS")
	flag.Parse()

	data, err := assets.StylesCSS()
	if err != nil {
		fmt.Fprintf(os.Stderr, "goatth: %v\n", err)
		os.Exit(1)
	}

	if dir := filepath.Dir(*out); dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "goatth: mkdir %s: %v\n", dir, err)
			os.Exit(1)
		}
	}

	if err := os.WriteFile(*out, data, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "goatth: write %s: %v\n", *out, err)
		os.Exit(1)
	}

	fmt.Printf("goatth: wrote %s (%d bytes)\n", *out, len(data))
}
