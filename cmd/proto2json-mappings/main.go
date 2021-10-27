package main

import (
	"fmt"

	"github.com/edfoh/proto-to-json-examples/cmd/proto2json-mappings/mapping"
)

func main() {
	mappings, err := mapping.Load()
	if err != nil {
		fmt.Printf("error loading mappings: %v", err)
		panic(err)
	}
	fmt.Printf("%v\n", mappings)
}
