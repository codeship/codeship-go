package main

import (
	"fmt"
	"os"

	"golang.org/x/build/version"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("USAGE: %v <version> [commands as normal]\n",
			os.Args[0])
		os.Exit(1)
	}

	v := os.Args[1]
	os.Args = append(os.Args[0:1], os.Args[2:]...)

	version.Run("go" + v)
}
