package classfile

import (
	"fmt"
)

// Parse descriptor
func ParseDescriptor(descriptor string) string {
	// Read the field type
	switch descriptor[0] {
	case 'B':
		return "byte"
	case 'C':
		return "char"
	case 'D':
		return "double"
	case 'F':
		return "float"
	case 'I':
		return "int"
	case 'J':
		return "long"
	case 'S':
		return "short"
	case 'Z':
		return "boolean"
	case 'L':
		return ParseDescriptor(descriptor[1:])
	case '[':
		return ParseDescriptor(descriptor[1:]) + "[]"
	default:
		panic(fmt.Errorf("Invalid basetype '%v' for descriptor '%v'", descriptor[0], descriptor))
	}
}
