package classfile

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

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
		binary.Read(reader, binary.BigEndian, &info.StringIndex)
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
