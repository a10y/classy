// Types and functions for parsing JVM Classfiles.
package classy

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
	AccessFlags       Access
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

// An entry that exists in the classfile's constant pool
type CpEntry interface {
	// Get the string representation of the tag
	StringTag() string
	// Get the raw integer tag
	RawTag() ConstantTag
	// Get a string representation of the entry
	Repr([]CpEntry) string
}

// Classfile field information
type FieldInfo struct {
	AccessFlags     Access
	NameIndex       uint16
	DescriptorIndex uint16
	AttrsCount      uint16
	Attrs           []AttrInfo
}

// Information corresponding to attributes
type AttrInfo struct {
	NameIndex  uint16
	AttrLength uint32
	AttrData   []byte
}

// Corresponds to the method_info type in the spec.
type MethodInfo struct {
	AccessFlags     Access
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

// Get the name of the attribute
func (i *AttrInfo) Name(cp []CpEntry) string {
	idx := i.NameIndex - 1
	ent := cp[idx].(*CONSTANT_Utf8_info)
	return string(ent.Bytes[:ent.Length])
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
