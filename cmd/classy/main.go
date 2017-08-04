package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/a10y/classy"
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

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %v FILENAME\n", os.Args[0])
	os.Exit(-1)
}

func main() {
	if len(os.Args) == 1 {
		usage()
	}

	data, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}
	classFile, err := classy.ReadClassFile(data)
	if classFile == nil {
		fmt.Fprintf(os.Stderr, "Error parsing %v: %v\n", os.Args[1], err.Error())
		os.Exit(-1)
	}

	AuxColorizer.Printf("Major:")
	fmt.Printf(" %v\n", classFile.MajorVersion)
	AuxColorizer.Printf("Minor:")
	fmt.Printf(" %v\n", classFile.MinorVersion)

	HeaderColorizer.Printf("\nConstantPool:")
	fmt.Printf(" (%v entries)\n", classFile.ConstantPoolCount-1)
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

func printCP(cf *classy.ClassFile) {
	for i, cpEntry := range cf.ConstantPool {
		// Skip over empty continuation slots for 8-byte constants
		if cpEntry == nil {
			continue
		}
		branch := "├──"
		branch2 := "│"
		if i == int(cf.ConstantPoolCount)-2 {
			branch = "└──"
			branch2 = ""
		}
		fmt.Printf("  %v %02d: %v\n", branch, i+1, cpEntry.StringTag())
		fmt.Printf("  %v\t\t%v\n", branch2, cpEntry.Repr(cf.ConstantPool))
	}
}

func printMethods(cf *classy.ClassFile) {
	for i, meth := range cf.Methods {
		branch := "├──"
		if i == int(cf.MethodsCount)-1 {
			branch = "└──"
		}
		name := meth.Name(cf.ConstantPool)
		paramSlice, ret := classy.ParseMethodDescriptor(meth.Descriptor(cf.ConstantPool))
		flags := FieldTypeColor.Sprint(classy.MethodFlagsRepr(meth.AccessFlags))
		ret = FieldTypeColor.Sprint(ret)
		name = FieldNameColor.Sprint(name)
		params := ParamTypeColor.Sprint(strings.Join(paramSlice, ", "))
		repr := fmt.Sprintf("%v %v %v(%v)\n", flags, ret, name, params)
		fmt.Printf("  %v %v", branch, repr)
	}
}

func printFields(cf *classy.ClassFile) {
	for i, field := range cf.Fields {
		branch := "├──"
		if i == int(cf.FieldsCount)-1 {
			branch = "└──"
		}
		flags := FieldTypeColor.Sprint(classy.FieldFlagsRepr(field.AccessFlags))
		desc := FieldTypeColor.Sprint(classy.ParseFieldDescriptor(field.Descriptor(cf.ConstantPool)))
		name := FieldNameColor.Sprint(field.Name(cf.ConstantPool))
		fmt.Printf("  %v %v %v %v\n", branch, flags, desc, name)
	}
}

func printAttrs(cf *classy.ClassFile) {
	for i, attr := range cf.Attrs {
		branch := "├──"
		if i == int(cf.AttrsCount)-1 {
			branch = "└──"
		}
		fmt.Printf("  %v %v\n", branch, attr.Name(cf.ConstantPool))
	}
}

// TODO: disassembly?
// TODO: colorize attributes and constant pool entries
// TODO: add constant pool info for methodref, instancemethodref, etc.
// TODO: show access flags for fields/methods
