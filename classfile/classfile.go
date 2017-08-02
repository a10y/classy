// Types and functions for parsing JVM Classfiles.
package classfile

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

// Shortcut for a uint16 that represents an access flag for the class, field, or method
type Access uint16

const (
	ACC_PUBLIC     Access = 0x0001
	ACC_FINAL             = 0x0010
	ACC_SUPER             = 0x0020
	ACC_INTERFACE         = 0x0200
	ACC_ABSTRACT          = 0x0400
	ACC_SYNTHETIC         = 0x1000
	ACC_ANNOTATION        = 0x2000
	ACC_ENUM              = 0x4000
)

// 1-byte header before every entry in the constant pool conveying type information.
// All tags are defined below and begin with "CONSTANT_"
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

// An in-memory representation of a .class file that is loadable by a JVM.
// This closely mirrors the actual serialized layout of a classfile and all its nested
// components, but with some syntactic flourishes to make it cleaner to work with. We are
// storing a parsed representation of the classfile, meaning runs of bytes in the file
// that correspond to arrays of variable-length items correspond to slices below, as with
// ConstantPool and Fields.
type ClassFile struct {
	Magic             uint32
	MinorVersion      uint16
	MajorVersion      uint16
	ConstantPoolCount uint16
	ConstantPool      []CpEntry
	AccessFlags       uint16
	ThisClass         uint16
	SuperClass        uint16
	InterfacesCount   uint16
	Interfaces        []uint16
	FieldsCount       uint16
	Fields            []FieldInfo
	MethodsCount      uint16
	Methods           []MethodInfo
	AttrsCount        uint16
	Attrs             []AttrInfo
}

// Deserialize a classfile from its raw byte representation into a ClassFile.
func ReadClassFile(raw []byte) *ClassFile {
	var classFile ClassFile

	reader := bytes.NewReader(raw)
	binary.Read(reader, binary.BigEndian, &classFile.Magic)
	binary.Read(reader, binary.BigEndian, &classFile.MinorVersion)
	binary.Read(reader, binary.BigEndian, &classFile.MajorVersion)
	binary.Read(reader, binary.BigEndian, &classFile.ConstantPoolCount)

	// Read constant pool entries
	for i := uint16(1); i <= classFile.ConstantPoolCount-1; i++ {
		ent := readCpEntry(reader)
		fmt.Printf("%v: %v\n", i, ent.Display())
		classFile.ConstantPool = append(classFile.ConstantPool, ent)
		// Check to see if the entry is one of the 8-byte varieties
		// If so, we skip a slot
		if _, ok := ent.(*CONSTANT_Double_info); ok {
			classFile.ConstantPool = append(classFile.ConstantPool, nil)
			i++
		}
		if _, ok := ent.(*CONSTANT_Long_info); ok {
			classFile.ConstantPool = append(classFile.ConstantPool, nil)
			i++
		}
	}

	binary.Read(reader, binary.BigEndian, &classFile.AccessFlags)
	binary.Read(reader, binary.BigEndian, &classFile.ThisClass)
	binary.Read(reader, binary.BigEndian, &classFile.SuperClass)
	binary.Read(reader, binary.BigEndian, &classFile.InterfacesCount)

	classFile.Interfaces = make([]uint16, classFile.InterfacesCount)
	for i := uint16(0); i < classFile.InterfacesCount; i++ {
		binary.Read(reader, binary.BigEndian, &classFile.Interfaces[i])
	}

	binary.Read(reader, binary.BigEndian, &classFile.FieldsCount)
	for i := uint16(0); i < classFile.FieldsCount; i++ {
		classFile.Fields = append(classFile.Fields, readField(reader))
	}

	binary.Read(reader, binary.BigEndian, &classFile.MethodsCount)
	for i := uint16(0); i < classFile.MethodsCount; i++ {
		classFile.Methods = append(classFile.Methods, readMethod(reader))
	}

	binary.Read(reader, binary.BigEndian, &classFile.AttrsCount)
	for i := uint16(0); i < classFile.AttrsCount; i++ {
		classFile.Attrs = append(classFile.Attrs, readAttr(reader))
	}

	return &classFile
}

// Read a constant pool entry from the classfile
func readCpEntry(reader *bytes.Reader) CpEntry {
	var tag ConstantTag
	binary.Read(reader, binary.BigEndian, &tag)
	//fmt.Printf("tag %v\n", tag)
	switch tag {
	case CONSTANT_Class:
		var info CONSTANT_Class_info
		info.Tag = tag
		binary.Read(reader, binary.BigEndian, &info.NameIndex)
		return &info
	case CONSTANT_Fieldref:
		var info CONSTANT_Fieldref_info
		info.Tag = tag
		binary.Read(reader, binary.BigEndian, &info.ClassIndex)
		binary.Read(reader, binary.BigEndian, &info.NameAndTypeIndex)
		return &info
	case CONSTANT_Methodref:
		var info CONSTANT_Methodref_info
		info.Tag = tag
		binary.Read(reader, binary.BigEndian, &info.ClassIndex)
		binary.Read(reader, binary.BigEndian, &info.NameAndTypeIndex)
		return &info
	case CONSTANT_InterfaceMethodref:
		var info CONSTANT_InterfaceMethodref_info
		info.Tag = tag
		binary.Read(reader, binary.BigEndian, &info.ClassIndex)
		binary.Read(reader, binary.BigEndian, &info.NameAndTypeIndex)
		return &info
	case CONSTANT_String:
		var info CONSTANT_String_info
		info.Tag = tag
		binary.Read(reader, binary.BigEndian, &info.stringIndex)
		return &info
	case CONSTANT_Integer:
		var info CONSTANT_Integer_info
		info.Tag = tag
		binary.Read(reader, binary.BigEndian, &info.Value)
		return &info
	case CONSTANT_Float:
		var info CONSTANT_Float_info
		info.Tag = tag
		binary.Read(reader, binary.BigEndian, &info.Value)
		return &info
	case CONSTANT_Long:
		var info CONSTANT_Long_info
		info.Tag = tag
		binary.Read(reader, binary.BigEndian, &info.HighBytes)
		binary.Read(reader, binary.BigEndian, &info.LowBytes)
		// Push a BS copy as well
		return &info
	case CONSTANT_Double:
		var info CONSTANT_Double_info
		info.Tag = tag
		binary.Read(reader, binary.BigEndian, &info.HighBytes)
		binary.Read(reader, binary.BigEndian, &info.LowBytes)
		return &info
	case CONSTANT_NameAndType:
		var info CONSTANT_NameAndType_info
		info.Tag = tag
		binary.Read(reader, binary.BigEndian, &info.NameIndex)
		binary.Read(reader, binary.BigEndian, &info.DescriptorIndex)
		return &info
	case CONSTANT_Utf8:
		var info CONSTANT_Utf8_info
		info.Tag = tag
		binary.Read(reader, binary.BigEndian, &info.Length)
		info.Bytes = make([]byte, info.Length)
		io.ReadFull(reader, info.Bytes)
		return &info
	case CONSTANT_MethodHandle:
		var info CONSTANT_MethodHandle_info
		info.Tag = tag
		binary.Read(reader, binary.BigEndian, &info.ReferenceKind)
		binary.Read(reader, binary.BigEndian, &info.ReferenceIndex)
		return &info
	case CONSTANT_MethodType:
		var info CONSTANT_MethodType_info
		info.Tag = tag
		binary.Read(reader, binary.BigEndian, &info.DescriptorIndex)
		return &info
	case CONSTANT_InvokeDynamic:
		var info CONSTANT_InvokeDynamic_info
		info.Tag = tag
		binary.Read(reader, binary.BigEndian, &info.BootstrapMethodAttrIndex)
		binary.Read(reader, binary.BigEndian, &info.NameAndTypeIndex)
		return &info
	default:
		panic(fmt.Errorf("Invalid Tag '%v' for constant pool entry", tag))
	}
}

// Read a field_info struct
func readField(reader *bytes.Reader) FieldInfo {
	var fieldInfo FieldInfo
	binary.Read(reader, binary.BigEndian, &fieldInfo.AccessFlags)
	binary.Read(reader, binary.BigEndian, &fieldInfo.NameIndex)
	binary.Read(reader, binary.BigEndian, &fieldInfo.DescriptorIndex)
	binary.Read(reader, binary.BigEndian, &fieldInfo.AttrsCount)
	for i := uint16(0); i < fieldInfo.AttrsCount; i++ {
		fieldInfo.Attrs = append(fieldInfo.Attrs, readAttr(reader))
	}
	return fieldInfo
}

// Read an attribute_info struct
func readAttr(reader *bytes.Reader) AttrInfo {
	var attrInfo AttrInfo
	binary.Read(reader, binary.BigEndian, &attrInfo.NameIndex)
	binary.Read(reader, binary.BigEndian, &attrInfo.AttrLength)
	attrInfo.AttrData = make([]byte, attrInfo.AttrLength)
	io.ReadFull(reader, attrInfo.AttrData)
	return attrInfo
}

// Reads a method_info struct
func readMethod(reader *bytes.Reader) MethodInfo {
	var methodInfo MethodInfo
	binary.Read(reader, binary.BigEndian, &methodInfo.AccessFlags)
	binary.Read(reader, binary.BigEndian, &methodInfo.NameIndex)
	binary.Read(reader, binary.BigEndian, &methodInfo.DescriptorIndex)
	binary.Read(reader, binary.BigEndian, &methodInfo.AttrsCount)
	for i := uint16(0); i < methodInfo.AttrsCount; i++ {
		methodInfo.Attrs = append(methodInfo.Attrs, readAttr(reader))
	}
	return methodInfo
}

// An entry that exists in the classfile's constant pool
type CpEntry interface {
	// Get the string representation of the tag
	StringTag() string
	// Get the raw integer tag
	RawTag() ConstantTag
	// Get a string representation of the entry
	Display() string
}

// Classfile field information
type FieldInfo struct {
	AccessFlags     uint16
	NameIndex       uint16
	DescriptorIndex uint16
	AttrsCount      uint16
	Attrs           []AttrInfo
}

// Get the String name of the field
// It requires looking up a CONSTANT_Utf8 entry in the constant pool.
func (i *FieldInfo) Name(cp []CpEntry) string {
	idx := i.NameIndex - 1
	ent := cp[idx].(*CONSTANT_Utf8_info)
	return string(ent.Bytes[:ent.Length])
}

// Get the string representation of field's type descriptor
func (i *FieldInfo) Descriptor(cp []CpEntry) string {
	idx := i.DescriptorIndex - 1
	ent := cp[idx].(*CONSTANT_Utf8_info)
	return string(ent.Bytes[:ent.Length])
}

// Information corresponding to attributes
type AttrInfo struct {
	NameIndex  uint16
	AttrLength uint32
	AttrData   []byte
}

// Get the name of the attribute
func (i *AttrInfo) Name(cp []CpEntry) string {
	idx := i.NameIndex - 1
	ent := cp[idx].(*CONSTANT_Utf8_info)
	return string(ent.Bytes[:ent.Length])
}

// Corresponds to the method_info type in the spec.
type MethodInfo struct {
	AccessFlags     uint16
	NameIndex       uint16
	DescriptorIndex uint16
	AttrsCount      uint16
	Attrs           []AttrInfo
}

// Get the name of the method
func (i *MethodInfo) Name(cp []CpEntry) string {
	idx := i.NameIndex - 1
	ent := cp[idx].(*CONSTANT_Utf8_info)
	return string(ent.Bytes[:ent.Length])
}

func (i *MethodInfo) Descriptor(cp []CpEntry) string {
	idx := i.DescriptorIndex - 1
	ent := cp[idx].(*CONSTANT_Utf8_info)
	return string(ent.Bytes[:ent.Length])
}

// Constant pool entry for classes.
// Corresponds to eponymous struct in the spec.
type CONSTANT_Class_info struct {
	Tag       ConstantTag
	NameIndex uint16
}

func (i *CONSTANT_Class_info) StringTag() string {
	return "CONSTANT_Class"
}

func (i *CONSTANT_Class_info) RawTag() ConstantTag {
	return i.Tag
}

func (i *CONSTANT_Class_info) Display() string {
	// Fill me in
	return i.StringTag()
}

// Constant pool entry for field references.
// Corresponds to eponymous struct in the spec.
type CONSTANT_Fieldref_info struct {
	Tag              ConstantTag
	ClassIndex       uint16
	NameAndTypeIndex uint16
}

func (i *CONSTANT_Fieldref_info) StringTag() string {
	return "CONSTANT_Fieldref"
}

func (i *CONSTANT_Fieldref_info) RawTag() ConstantTag {
	return i.Tag
}

func (i *CONSTANT_Fieldref_info) Display() string {
	// Fill me in
	return i.StringTag()
}

// Constant pool entry referencing a method.
// Corresponds to eponymous struct in the spec.
type CONSTANT_Methodref_info struct {
	Tag              ConstantTag
	ClassIndex       uint16
	NameAndTypeIndex uint16
}

func (i *CONSTANT_Methodref_info) StringTag() string {
	return "CONSTANT_Methodref"
}

func (i *CONSTANT_Methodref_info) RawTag() ConstantTag {
	return i.Tag
}

func (i *CONSTANT_Methodref_info) Display() string {
	// Fill me in
	return i.StringTag()
}

// Corresponds to eponymous struct in the spec.
type CONSTANT_InterfaceMethodref_info struct {
	Tag              ConstantTag
	ClassIndex       uint16
	NameAndTypeIndex uint16
}

func (i *CONSTANT_InterfaceMethodref_info) StringTag() string {
	return "CONSTANT_InterfaceMethodref"
}

func (i *CONSTANT_InterfaceMethodref_info) RawTag() ConstantTag {
	return i.Tag
}

func (i *CONSTANT_InterfaceMethodref_info) Display() string {
	// Fill me in
	return i.StringTag()
}

// Corresponds to eponymous struct in the spec.
type CONSTANT_String_info struct {
	Tag         ConstantTag
	stringIndex uint16
}

func (i *CONSTANT_String_info) StringTag() string {
	return "CONSTANT_String"
}

func (i *CONSTANT_String_info) RawTag() ConstantTag {
	return i.Tag
}

func (i *CONSTANT_String_info) Display() string {
	// Fill me in
	return i.StringTag()
}

// Corresponds to eponymous struct in the spec.
type CONSTANT_Integer_info struct {
	Tag   ConstantTag
	Value uint32
}

func (i *CONSTANT_Integer_info) StringTag() string {
	return "CONSTANT_Integer"
}

func (i *CONSTANT_Integer_info) RawTag() ConstantTag {
	return i.Tag
}

func (i *CONSTANT_Integer_info) Display() string {
	// Fill me in
	return i.StringTag()
}

// Corresponds to eponymous struct in the spec.
type CONSTANT_Float_info struct {
	Tag   ConstantTag
	Value float32
}

func (i *CONSTANT_Float_info) StringTag() string {
	return "CONSTANT_Float"
}

func (i *CONSTANT_Float_info) RawTag() ConstantTag {
	return i.Tag
}

func (i *CONSTANT_Float_info) Display() string {
	// Fill me in
	return i.StringTag()
}

// Corresponds to eponymous struct in the spec.
type CONSTANT_Long_info struct {
	Tag       ConstantTag
	HighBytes uint32
	LowBytes  uint32
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

func (i *CONSTANT_Long_info) Display() string {
	// Fill me in
	return i.StringTag()
}

// Corresponds to eponymous struct in the spec.
type CONSTANT_Double_info struct {
	Tag       ConstantTag
	HighBytes uint32
	LowBytes  uint32
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

func (i *CONSTANT_Double_info) Display() string {
	return fmt.Sprintf("%v\t\t\t(%v)", i.StringTag(), i.Value())
}

// Corresponds to eponymous struct in the spec.
type CONSTANT_NameAndType_info struct {
	Tag             ConstantTag
	NameIndex       uint16
	DescriptorIndex uint16
}

func (i *CONSTANT_NameAndType_info) StringTag() string {
	return "CONSTANT_NameAndType"
}

func (i *CONSTANT_NameAndType_info) Display() string {
	// Fill me in
	return i.StringTag()
}

func (i *CONSTANT_NameAndType_info) RawTag() ConstantTag {
	return i.Tag
}

// Corresponds to eponymous struct in the spec.
type CONSTANT_Utf8_info struct {
	Tag    ConstantTag
	Length uint16
	Bytes  []byte
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

func (i *CONSTANT_Utf8_info) Display() string {
	// Fill me in
	return fmt.Sprintf("%v\t\t\t'%v'", i.StringTag(), i.Value())
}

// Corresponds to eponymous struct in the spec.
type CONSTANT_MethodHandle_info struct {
	Tag            ConstantTag
	ReferenceKind  byte
	ReferenceIndex uint16
}

func (i *CONSTANT_MethodHandle_info) StringTag() string {
	return "CONSTANT_MethodHandle"
}

func (i *CONSTANT_MethodHandle_info) RawTag() ConstantTag {
	return i.Tag
}

func (i *CONSTANT_MethodHandle_info) Display() string {
	// Fill me in
	return i.StringTag()
}

// Corresponds to eponymous struct in the spec.
type CONSTANT_MethodType_info struct {
	Tag             ConstantTag
	DescriptorIndex uint16
}

func (i *CONSTANT_MethodType_info) StringTag() string {
	return "CONSTANT_MethodType"
}

func (i *CONSTANT_MethodType_info) RawTag() ConstantTag {
	return i.Tag
}

func (i *CONSTANT_MethodType_info) Display() string {
	// Fill me in
	return i.StringTag()
}

// Corresponds to eponymous struct in the spec.
type CONSTANT_InvokeDynamic_info struct {
	Tag                      ConstantTag
	BootstrapMethodAttrIndex uint16
	NameAndTypeIndex         uint16
}

func (i *CONSTANT_InvokeDynamic_info) StringTag() string {
	return "CONSTANT_InvokeDynamic"
}

func (i *CONSTANT_InvokeDynamic_info) RawTag() ConstantTag {
	return i.Tag
}

func (i *CONSTANT_InvokeDynamic_info) Display() string {
	// Fill me in
	return i.StringTag()
}
