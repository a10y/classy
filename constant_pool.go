package classy

import (
	"fmt"
	"math"
)

// ConstantTag is a 1-byte header before every entry in the constant pool conveying type
// information.  All tags are defined below and begin with "CONSTANT_"
type ConstantTag byte

const (
	CONSTANT_Class              ConstantTag = 7
	CONSTANT_Fieldref                       = 9
	CONSTANT_Methodref                      = 10
	CONSTANT_InterfaceMethodref             = 11
	CONSTANT_String                         = 8
	CONSTANT_Integer                        = 3
	CONSTANT_Float                          = 4
	CONSTANT_Long                           = 5
	CONSTANT_Double                         = 6
	CONSTANT_NameAndType                    = 12
	CONSTANT_Utf8                           = 1
	CONSTANT_MethodHandle                   = 15
	CONSTANT_MethodType                     = 16
	CONSTANT_InvokeDynamic                  = 18
)

// CONSTANT_Class_info represents constant pool entries for classes.
// Corresponds to eponymous struct in the spec.
type CONSTANT_Class_info struct {
	Tag       ConstantTag
	NameIndex uint16
}

// Constant pool entry for field references.
// Corresponds to eponymous struct in the spec.
type CONSTANT_Fieldref_info struct {
	Tag              ConstantTag
	ClassIndex       uint16
	NameAndTypeIndex uint16
}

// Constant pool entry referencing a method.
// Corresponds to eponymous struct in the spec.
type CONSTANT_Methodref_info struct {
	Tag              ConstantTag
	ClassIndex       uint16
	NameAndTypeIndex uint16
}

// Corresponds to eponymous struct in the spec.
type CONSTANT_InterfaceMethodref_info struct {
	Tag              ConstantTag
	ClassIndex       uint16
	NameAndTypeIndex uint16
}

// Corresponds to eponymous struct in the spec.
type CONSTANT_String_info struct {
	Tag         ConstantTag
	StringIndex uint16
}

// Corresponds to eponymous struct in the spec.
type CONSTANT_Integer_info struct {
	Tag   ConstantTag
	Value uint32
}

// Corresponds to eponymous struct in the spec.
type CONSTANT_Float_info struct {
	Tag   ConstantTag
	Value float32
}

// Corresponds to eponymous struct in the spec.
type CONSTANT_Long_info struct {
	Tag       ConstantTag
	HighBytes uint32
	LowBytes  uint32
}

// Corresponds to eponymous struct in the spec.
type CONSTANT_Double_info struct {
	Tag       ConstantTag
	HighBytes uint32
	LowBytes  uint32
}

// Corresponds to eponymous struct in the spec.
type CONSTANT_NameAndType_info struct {
	Tag             ConstantTag
	NameIndex       uint16
	DescriptorIndex uint16
}

// Corresponds to eponymous struct in the spec.
type CONSTANT_Utf8_info struct {
	Tag    ConstantTag
	Length uint16
	Bytes  []byte
}

// Corresponds to eponymous struct in the spec.
type CONSTANT_MethodHandle_info struct {
	Tag            ConstantTag
	ReferenceKind  byte
	ReferenceIndex uint16
}

// Corresponds to eponymous struct in the spec.
type CONSTANT_MethodType_info struct {
	Tag             ConstantTag
	DescriptorIndex uint16
}

// Corresponds to eponymous struct in the spec.
type CONSTANT_InvokeDynamic_info struct {
	Tag                      ConstantTag
	BootstrapMethodAttrIndex uint16
	NameAndTypeIndex         uint16
}

func (i *CONSTANT_Class_info) StringTag() string {
	return "CONSTANT_Class"
}

func (i *CONSTANT_Class_info) RawTag() ConstantTag {
	return i.Tag
}

func (i *CONSTANT_Class_info) Name(cp []CpEntry) string {
	idx := i.NameIndex - 1
	ent := cp[idx].(*CONSTANT_Utf8_info)
	return string(ent.Bytes[:ent.Length])
}

func (i *CONSTANT_Class_info) Repr(cp []CpEntry) string {
	return i.Name(cp)
}

func (i *CONSTANT_Fieldref_info) StringTag() string {
	return "CONSTANT_Fieldref"
}

func (i *CONSTANT_Fieldref_info) RawTag() ConstantTag {
	return i.Tag
}

func (i *CONSTANT_Fieldref_info) Repr(ignored []CpEntry) string {
	return ""
}

func (i *CONSTANT_Methodref_info) StringTag() string {
	return "CONSTANT_Methodref"
}

func (i *CONSTANT_Methodref_info) RawTag() ConstantTag {
	return i.Tag
}

func (i *CONSTANT_Methodref_info) Repr(ignored []CpEntry) string {
	// Fill me in
	return ""
}

func (i *CONSTANT_InterfaceMethodref_info) StringTag() string {
	return "CONSTANT_InterfaceMethodref"
}

func (i *CONSTANT_InterfaceMethodref_info) RawTag() ConstantTag {
	return i.Tag
}

func (i *CONSTANT_InterfaceMethodref_info) Repr(ignored []CpEntry) string {
	// Fill me in
	return ""
}

func (i *CONSTANT_String_info) StringTag() string {
	return "CONSTANT_String"
}

func (i *CONSTANT_String_info) RawTag() ConstantTag {
	return i.Tag
}

func (i *CONSTANT_String_info) Repr(cp []CpEntry) string {
	idx := i.StringIndex - 1
	ent := cp[idx].(*CONSTANT_Utf8_info)
	val := string(ent.Bytes[:ent.Length])
	return fmt.Sprintf("\"%v\"", val)
}

func (i *CONSTANT_Integer_info) StringTag() string {
	return "CONSTANT_Integer"
}

func (i *CONSTANT_Integer_info) RawTag() ConstantTag {
	return i.Tag
}

func (i *CONSTANT_Integer_info) Repr(ignored []CpEntry) string {
	// Fill me in
	return fmt.Sprintf("%v", i.Value)
}

func (i *CONSTANT_Float_info) StringTag() string {
	return "CONSTANT_Float"
}

func (i *CONSTANT_Float_info) RawTag() ConstantTag {
	return i.Tag
}

func (i *CONSTANT_Float_info) Repr(ignored []CpEntry) string {
	// Fill me in
	return fmt.Sprintf("%v", i.Value)
}

func (li *CONSTANT_Long_info) Value() uint64 {
	res := uint64(li.HighBytes) << 32
	res &= uint64(li.LowBytes)
	return res
}

func (i *CONSTANT_Long_info) StringTag() string {
	return "CONSTANT_Long"
}

func (i *CONSTANT_Long_info) RawTag() ConstantTag {
	return i.Tag
}

func (i *CONSTANT_Long_info) Repr(ignored []CpEntry) string {
	// Fill me in
	return fmt.Sprintf("%v", i.Value())
}

func (i *CONSTANT_Double_info) Value() float64 {
	binary := uint64(i.HighBytes)
	binary &= uint64(i.LowBytes)
	return math.Float64frombits(binary)
}

func (i *CONSTANT_Double_info) StringTag() string {
	return "CONSTANT_Double"
}

func (i *CONSTANT_Double_info) RawTag() ConstantTag {
	return i.Tag
}

func (i *CONSTANT_Double_info) Repr(ignored []CpEntry) string {
	return fmt.Sprintf("%v", i.Value())
}

func (i *CONSTANT_NameAndType_info) StringTag() string {
	return "CONSTANT_NameAndType"
}

func (i *CONSTANT_NameAndType_info) Repr(ignored []CpEntry) string {
	// Fill me in
	return ""
}

func (i *CONSTANT_NameAndType_info) RawTag() ConstantTag {
	return i.Tag
}

func (i *CONSTANT_Utf8_info) StringTag() string {
	return "CONSTANT_Utf8"
}

func (i *CONSTANT_Utf8_info) RawTag() ConstantTag {
	return i.Tag
}

func (i *CONSTANT_Utf8_info) Value() string {
	return string(i.Bytes[:i.Length])
}

func (i *CONSTANT_Utf8_info) Repr(ignored []CpEntry) string {
	// Fill me in
	return fmt.Sprintf("\"%v\"", i.Value())
}

func (i *CONSTANT_MethodHandle_info) StringTag() string {
	return "CONSTANT_MethodHandle"
}

func (i *CONSTANT_MethodHandle_info) RawTag() ConstantTag {
	return i.Tag
}

func (i *CONSTANT_MethodHandle_info) Repr(ignored []CpEntry) string {
	// Fill me in
	return ""
}

func (i *CONSTANT_MethodType_info) StringTag() string {
	return "CONSTANT_MethodType"
}

func (i *CONSTANT_MethodType_info) RawTag() ConstantTag {
	return i.Tag
}

func (i *CONSTANT_MethodType_info) Repr(ignored []CpEntry) string {
	// Fill me in
	return ""
}

func (i *CONSTANT_InvokeDynamic_info) StringTag() string {
	return "CONSTANT_InvokeDynamic"
}

func (i *CONSTANT_InvokeDynamic_info) RawTag() ConstantTag {
	return i.Tag
}

func (i *CONSTANT_InvokeDynamic_info) Repr(ignored []CpEntry) string {
	// Fill me in
	return ""
}
