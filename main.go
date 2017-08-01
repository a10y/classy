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
	fmt.Printf("\nConstantPool: %v\n", classFile.ConstantPoolCount)
	printCP(classFile)
	fmt.Printf("\nMethods: %v\n", classFile.MethodsCount)
	printMethods(classFile)
	fmt.Printf("\nFields: %v\n", classFile.FieldsCount)
	printFields(classFile)
}

func printCP(cf *classfile.ClassFile) {
	for i, cpEntry := range cf.ConstantPool {
		fmt.Printf("  |%02d: %v\n", i+1, cpEntry.Display())
	}
}

func printMethods(cf *classfile.ClassFile) {
	for _, meth := range cf.Methods {
		fmt.Printf("  | %v\n", meth.Name(cf.ConstantPool))
	}
}

func printFields(cf *classfile.ClassFile) {
	for _, field := range cf.Fields {
		fmt.Printf("  | %v\n", field.Name(cf.ConstantPool))
	}
}
