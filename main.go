package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/a10y/classy/classfile"
	"github.com/fatih/color"
)

var (
	HeaderColorizer  *color.Color = color.New(color.FgYellow, color.Bold)
	SuccessColorizer              = color.New(color.BgWhite, color.FgBlack)
	ErrorColorizer                = color.New(color.FgWhite, color.BgRed)
	AuxColorizer                  = color.New(color.FgBlue, color.Bold)
	FieldTypeColor                = color.New(color.FgHiMagenta)
	FieldNameColor                = color.New(color.FgCyan, color.Bold)
	ParamTypeColor                = color.New(color.FgRed)
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
		// Skip over empty continuation slots for 8-byte constants
		if cpEntry == nil {
			continue
		}
		branch := "├──"
		if i == int(cf.ConstantPoolCount)-2 {
			branch = "└──"
		}
		fmt.Printf("  %v %02d: %v\n", branch, i+1, cpEntry.StringTag())
		fmt.Printf("  │\t\t%v\n", cpEntry.Repr(cf.ConstantPool))
	}
}

func printMethods(cf *classfile.ClassFile) {
	for i, meth := range cf.Methods {
		branch := "├──"
		if i == int(cf.MethodsCount)-1 {
			branch = "└──"
		}
		name := meth.Name(cf.ConstantPool)
		paramSlice, ret := classfile.ParseMethodDescriptor(meth.Descriptor(cf.ConstantPool))
		ret = FieldTypeColor.Sprint(ret)
		name = FieldNameColor.Sprint(name)
		params := ParamTypeColor.Sprint(strings.Join(paramSlice, ", "))
		repr := fmt.Sprintf("%v %v(%v)\n", ret, name, params)
		fmt.Printf("  %v %v", branch, repr)
	}
}

func printFields(cf *classfile.ClassFile) {
	for i, field := range cf.Fields {
		branch := "├──"
		if i == int(cf.FieldsCount)-1 {
			branch = "└──"
		}
		desc := FieldTypeColor.Sprint(classfile.ParseFieldDescriptor(field.Descriptor(cf.ConstantPool)))
		name := FieldNameColor.Sprint(field.Name(cf.ConstantPool))
		fmt.Printf("  %v %v %v\n", branch, desc, name)
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
