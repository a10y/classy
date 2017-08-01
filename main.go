package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/a10y/classy/classfile"
)

func main() {
	// Read all bytes in
	data, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}
	classFile := classfile.ReadClassFile(data)
	fmt.Printf("Magic: 0x%X\n", classFile.Magic)
	fmt.Printf("Major: %v\n", classFile.MajorVersion)
	fmt.Printf("Minor: %v\n", classFile.MinorVersion)
}
