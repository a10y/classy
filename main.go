package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/a10y/classy/classfile"
	"github.com/fatih/color"
)

var (
	HeaderColorizer  *color.Color = color.New(color.FgYellow, color.Bold)
	SuccessColorizer              = color.New(color.BgWhite, color.FgBlack)
	ErrorColorizer                = color.New(color.FgWhite, color.BgRed)
	AuxColorizer                  = color.New(color.FgBlue, color.Bold)
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

	AuxColorizer.Printf("Magic:")
	fmt.Printf(" 0x%X (%v)\n", classFile.Magic, validMsg)
	AuxColorizer.Printf("Major:")
	fmt.Printf(" %v\n", classFile.MajorVersion)
	AuxColorizer.Printf("Minor:")
	fmt.Printf(" %v\n", classFile.MinorVersion)

	HeaderColorizer.Printf("\nConstantPool:")
	fmt.Printf(" (%v entries)\n", classFile.ConstantPoolCount)
	printCP(classFile)

	HeaderColorizer.Printf("\nMethods:")
	fmt.Printf(" (%v entries)\n", classFile.MethodsCount)
	printMethods(classFile)

	HeaderColorizer.Printf("\nFields:")
	fmt.Printf(" (%v entries)\n", classFile.FieldsCount)
	printFields(classFile)

	HeaderColorizer.Printf("\nAttrs:")
	fmt.Printf(" (%v entries)\n", classFile.AttrsCount)
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
