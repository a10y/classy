package classy

import (
	"strings"
)

// Shortcut for a uint16 that represents an access flag for the class, field, or method
type Access uint16

const (
	AccPublic    Access = 0x0001
	AccPrivate          = 0x0002
	AccProtected        = 0x0004
	AccStatic           = 0x0008
	AccFinal            = 0x0010
	AccSuper            = 0x0020
	AccVolatile         = 0x0040

	AccTransient = 0x0080
	AccVarargs   = 0x0080

	AccNative     = 0x0100
	AccInterface  = 0x0200
	AccAbstract   = 0x0400
	AccStrict     = 0x0800
	AccSynthetic  = 0x1000
	AccAnnotation = 0x2000
	AccEnum       = 0x4000
)

// MethodFlagsRepr returns the string representation of flags for a method, in the order
// one would expect them to appear if written in a Java source file.
func MethodFlagsRepr(acc Access) string {
	var text []string

	if (acc & AccPublic) > 0 {
		text = append(text, "public")
	}

	if (acc & AccPrivate) > 0 {
		text = append(text, "private")
	}

	if (acc & AccProtected) > 0 {
		text = append(text, "protected")
	}

	if (acc & AccAbstract) > 0 {
		text = append(text, "abstract")
	}

	if (acc & AccStatic) > 0 {
		text = append(text, "static")
	}

	if (acc & AccFinal) > 0 {
		text = append(text, "final")
	}

	if (acc & AccVolatile) > 0 {
		text = append(text, "volatile")
	}

	if (acc & AccNative) > 0 {
		text = append(text, "native")
	}

	if (acc & AccEnum) > 0 {
		text = append(text, "enum")
	}

	return strings.Join(text, " ")
}

// FieldFlagsRepr returns the string representation of modifiers for a field, in the
// order one would expect them to appear in a Java source file.
func FieldFlagsRepr(acc Access) string {
	var text []string

	if (acc & AccPublic) > 0 {
		text = append(text, "public")
	}

	if (acc & AccPrivate) > 0 {
		text = append(text, "private")
	}

	if (acc & AccProtected) > 0 {
		text = append(text, "protected")
	}

	if (acc & AccStatic) > 0 {
		text = append(text, "static")
	}

	if (acc & AccFinal) > 0 {
		text = append(text, "final")
	}

	if (acc & AccVolatile) > 0 {
		text = append(text, "volatile")
	}

	if (acc & AccTransient) > 0 {
		text = append(text, "transient")
	}

	return strings.Join(text, " ")
}
