package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/a10y/classy/classfile"
	"github.com/fatih/color"
)

var (
	HeaderPrinter    *color.Color = color.New(color.FgYellow)
	SuccessColorizer              = color.New(color.BgGreen, color.FgWhite)
	ErrorColorizer                = color.New(color.FgWhite, color.BgRed)
)

func main() {
	data, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}
	classFile := classfile.ReadClassFile(data)

	validMsg := SuccessColorizer.Sprint("valid")
	if classFile.Magic != 0xCAFEBABE {
		validMsg = ErrorColorizer.Sprintf("INVALID")
	}

	fmt.Printf("Magic: 0x%X (%v)\n", classFile.Magic, validMsg)
	fmt.Printf("Major: %v\n", classFile.MajorVersion)
	fmt.Printf("Minor: %v\n", classFile.MinorVersion)

	HeaderPrinter.Printf("\nConstantPool:")
	fmt.Printf(" (%v entries)\n", classFile.ConstantPoolCount)
	printCP(classFile)

	HeaderPrinter.Printf("\nMethods:")
	fmt.Printf(" (%v entries)\n", classFile.MethodsCount)
	printMethods(classFile)

	HeaderPrinter.Printf("\nFields:")
	fmt.Printf(" (%v entries)\n", classFile.FieldsCount)
	printFields(classFile)

	fmt.Printf("\nAttrs: (%v entries)\n", classFile.AttrsCount)
	printAttrs(classFile)
}

func printCP(cf *classfile.ClassFile) {
	for i, cpEntry := range cf.ConstantPool {
		branch := "├──"
		if i == int(cf.ConstantPoolCount)-2 {
			branch = "└──"
		}
		fmt.Printf("  %v %02d: %v\n", branch, i+1, cpEntry.Display())
	}
}

func printMethods(cf *classfile.ClassFile) {
	for i, meth := range cf.Methods {
		branch := "├──"
		if i == int(cf.MethodsCount)-1 {
			branch = "└──"
		}
		fmt.Printf("  %v %v\n", branch, meth.Name(cf.ConstantPool))
	}
}

func printFields(cf *classfile.ClassFile) {
	for i, field := range cf.Fields {
		branch := "├──"
		if i == int(cf.FieldsCount)-1 {
			branch = "└──"
		}
		fmt.Printf("  %v %v\n", branch, field.Name(cf.ConstantPool))
	}
}

func printAttrs(cf *classfile.ClassFile) {
	for i, attr := range cf.Attrs {
		branch := "├──"
		if i == int(cf.AttrsCount)-1 {
			branch = "└──"
		}
		fmt.Printf("  %v %v\n", branch, attr.Name(cf.ConstantPool))
	}
}
