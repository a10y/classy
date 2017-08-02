package classy

import (
	"fmt"
	"strings"
)

// Parse descriptor
func ParseFieldDescriptor(descriptor string) string {
	name, _ := parseFieldDescriptor(descriptor)
	return name
}

func parseFieldDescriptor(descriptor string) (string, string) {
	// Read the field type
	switch descriptor[0] {
	case 'B':
		return "byte", descriptor[1:]
	case 'C':
		return "char", descriptor[1:]
	case 'D':
		return "double", descriptor[1:]
	case 'F':
		return "float", descriptor[1:]
	case 'I':
		return "int", descriptor[1:]
	case 'J':
		return "long", descriptor[1:]
	case 'S':
		return "short", descriptor[1:]
	case 'Z':
		return "boolean", descriptor[1:]
	case 'L':
		// Continue reading until we encounter a ';'
		pos := strings.IndexRune(descriptor, ';')
		if pos < 0 {
			panic(fmt.Errorf("Cannot find termination to descriptor"))
		}
		// Replace '/' with '.' for qualified class names
		clsName := strings.Replace(descriptor[1:pos], "/", ".", -1)
		return clsName, descriptor[pos+1:]
	case '[':
		name, rest := parseFieldDescriptor(descriptor[1:])
		return name + "[]", rest
	case 'V':
		return "void", descriptor[1:]
	default:
		panic(fmt.Errorf("Invalid basetype '%v' for descriptor '%v'", descriptor[0], descriptor))
	}
}

func ParseMethodDescriptor(descriptor string) ([]string, string) {
	var paramTypes []string
	if descriptor[0] != '(' {
		panic(fmt.Errorf("Invalid method descriptor '%v'", descriptor))
	}

	descriptor = descriptor[1:]
	finished := false
	for !finished {
		if len(descriptor) == 0 {
			panic(fmt.Errorf("this shouldn't happen!!!"))
		}

		switch descriptor[0] {
		case ')':
			descriptor = descriptor[1:]
			finished = true
			break
		default:
			name, desc := parseFieldDescriptor(descriptor)
			paramTypes = append(paramTypes, name)
			descriptor = desc
		}
	}
	retType, _ := parseFieldDescriptor(descriptor)
	return paramTypes, retType
}
