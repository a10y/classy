package classy

import (
	"strings"
)

// Shortcut for a uint16 that represents an access flag for the class, field, or method
type Access uint16

const (
	ACC_PUBLIC    Access = 0x0001
	ACC_PRIVATE          = 0x0002
	ACC_PROTECTED        = 0x0004
	ACC_STATIC           = 0x0008
	ACC_FINAL            = 0x0010
	ACC_SUPER            = 0x0020
	ACC_VOLATILE         = 0x0040

	ACC_TRANSIENT = 0x0080
	ACC_VARARGS   = 0x0080

	ACC_NATIVE     = 0x0100
	ACC_INTERFACE  = 0x0200
	ACC_ABSTRACT   = 0x0400
	ACC_STRICT     = 0x0800
	ACC_SYNTHETIC  = 0x1000
	ACC_ANNOTATION = 0x2000
	ACC_ENUM       = 0x4000
)

func MethodFlagsRepr(acc Access) string {
	var text []string

	if (acc & ACC_PUBLIC) > 0 {
		text = append(text, "public")
	}

	if (acc & ACC_PRIVATE) > 0 {
		text = append(text, "private")
	}

	if (acc & ACC_PROTECTED) > 0 {
		text = append(text, "protected")
	}

	if (acc & ACC_ABSTRACT) > 0 {
		text = append(text, "abstract")
	}

	if (acc & ACC_STATIC) > 0 {
		text = append(text, "static")
	}

	if (acc & ACC_FINAL) > 0 {
		text = append(text, "final")
	}

	if (acc & ACC_VOLATILE) > 0 {
		text = append(text, "volatile")
	}

	if (acc & ACC_NATIVE) > 0 {
		text = append(text, "native")
	}

	if (acc & ACC_ENUM) > 0 {
		text = append(text, "enum")
	}

	return strings.Join(text, " ")
}

func FieldFlagsRepr(acc Access) string {
	var text []string

	if (acc & ACC_PUBLIC) > 0 {
		text = append(text, "public")
	}

	if (acc & ACC_PRIVATE) > 0 {
		text = append(text, "private")
	}

	if (acc & ACC_PROTECTED) > 0 {
		text = append(text, "protected")
	}

	if (acc & ACC_STATIC) > 0 {
		text = append(text, "static")
	}

	if (acc & ACC_FINAL) > 0 {
		text = append(text, "final")
	}

	if (acc & ACC_VOLATILE) > 0 {
		text = append(text, "volatile")
	}

	if (acc & ACC_TRANSIENT) > 0 {
		text = append(text, "transient")
	}

	return strings.Join(text, " ")
}
