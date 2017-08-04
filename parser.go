package classy

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

const ClassFileMagic uint32 = 0xCAFEBABE

// Deserialize a classfile from its raw byte representation into a ClassFile.
func ReadClassFile(raw []byte) (classFile *ClassFile, err error) {
	defer func() {
		if e := recover(); e != nil {
			classFile = nil
			err = e.(error)
		}
	}()
	classFile = new(ClassFile)
	err = nil

	reader := bytes.NewReader(raw)
	safeReadBinary(reader, binary.BigEndian, &classFile.Magic)

	if classFile.Magic != ClassFileMagic {
		panic(fmt.Errorf("Invalid magic: %v", classFile.Magic))
	}

	safeReadBinary(reader, binary.BigEndian, &classFile.MinorVersion)
	safeReadBinary(reader, binary.BigEndian, &classFile.MajorVersion)
	safeReadBinary(reader, binary.BigEndian, &classFile.ConstantPoolCount)

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

	safeReadBinary(reader, binary.BigEndian, &classFile.AccessFlags)
	safeReadBinary(reader, binary.BigEndian, &classFile.ThisClass)
	safeReadBinary(reader, binary.BigEndian, &classFile.SuperClass)
	safeReadBinary(reader, binary.BigEndian, &classFile.InterfacesCount)

	classFile.Interfaces = make([]uint16, classFile.InterfacesCount)
	for i := uint16(0); i < classFile.InterfacesCount; i++ {
		safeReadBinary(reader, binary.BigEndian, &classFile.Interfaces[i])
	}

	safeReadBinary(reader, binary.BigEndian, &classFile.FieldsCount)
	for i := uint16(0); i < classFile.FieldsCount; i++ {
		classFile.Fields = append(classFile.Fields, readField(reader))
	}

	safeReadBinary(reader, binary.BigEndian, &classFile.MethodsCount)
	for i := uint16(0); i < classFile.MethodsCount; i++ {
		classFile.Methods = append(classFile.Methods, readMethod(reader))
	}

	safeReadBinary(reader, binary.BigEndian, &classFile.AttrsCount)
	for i := uint16(0); i < classFile.AttrsCount; i++ {
		classFile.Attrs = append(classFile.Attrs, readAttr(reader))
	}

	return
}

// Read a constant pool entry from the classfile
func readCpEntry(reader *bytes.Reader) CpEntry {
	var tag ConstantTag
	safeReadBinary(reader, binary.BigEndian, &tag)
	//fmt.Printf("tag %v\n", tag)
	switch tag {
	case CONSTANT_Class:
		var info CONSTANT_Class_info
		info.Tag = tag
		safeReadBinary(reader, binary.BigEndian, &info.NameIndex)
		return &info
	case CONSTANT_Fieldref:
		var info CONSTANT_Fieldref_info
		info.Tag = tag
		safeReadBinary(reader, binary.BigEndian, &info.ClassIndex)
		safeReadBinary(reader, binary.BigEndian, &info.NameAndTypeIndex)
		return &info
	case CONSTANT_Methodref:
		var info CONSTANT_Methodref_info
		info.Tag = tag
		safeReadBinary(reader, binary.BigEndian, &info.ClassIndex)
		safeReadBinary(reader, binary.BigEndian, &info.NameAndTypeIndex)
		return &info
	case CONSTANT_InterfaceMethodref:
		var info CONSTANT_InterfaceMethodref_info
		info.Tag = tag
		safeReadBinary(reader, binary.BigEndian, &info.ClassIndex)
		safeReadBinary(reader, binary.BigEndian, &info.NameAndTypeIndex)
		return &info
	case CONSTANT_String:
		var info CONSTANT_String_info
		info.Tag = tag
		safeReadBinary(reader, binary.BigEndian, &info.StringIndex)
		return &info
	case CONSTANT_Integer:
		var info CONSTANT_Integer_info
		info.Tag = tag
		safeReadBinary(reader, binary.BigEndian, &info.Value)
		return &info
	case CONSTANT_Float:
		var info CONSTANT_Float_info
		info.Tag = tag
		safeReadBinary(reader, binary.BigEndian, &info.Value)
		return &info
	case CONSTANT_Long:
		var info CONSTANT_Long_info
		info.Tag = tag
		safeReadBinary(reader, binary.BigEndian, &info.HighBytes)
		safeReadBinary(reader, binary.BigEndian, &info.LowBytes)
		// Push a BS copy as well
		return &info
	case CONSTANT_Double:
		var info CONSTANT_Double_info
		info.Tag = tag
		safeReadBinary(reader, binary.BigEndian, &info.HighBytes)
		safeReadBinary(reader, binary.BigEndian, &info.LowBytes)
		return &info
	case CONSTANT_NameAndType:
		var info CONSTANT_NameAndType_info
		info.Tag = tag
		safeReadBinary(reader, binary.BigEndian, &info.NameIndex)
		safeReadBinary(reader, binary.BigEndian, &info.DescriptorIndex)
		return &info
	case CONSTANT_Utf8:
		var info CONSTANT_Utf8_info
		info.Tag = tag
		safeReadBinary(reader, binary.BigEndian, &info.Length)
		info.Bytes = make([]byte, info.Length)
		io.ReadFull(reader, info.Bytes)
		return &info
	case CONSTANT_MethodHandle:
		var info CONSTANT_MethodHandle_info
		info.Tag = tag
		safeReadBinary(reader, binary.BigEndian, &info.ReferenceKind)
		safeReadBinary(reader, binary.BigEndian, &info.ReferenceIndex)
		return &info
	case CONSTANT_MethodType:
		var info CONSTANT_MethodType_info
		info.Tag = tag
		safeReadBinary(reader, binary.BigEndian, &info.DescriptorIndex)
		return &info
	case CONSTANT_InvokeDynamic:
		var info CONSTANT_InvokeDynamic_info
		info.Tag = tag
		safeReadBinary(reader, binary.BigEndian, &info.BootstrapMethodAttrIndex)
		safeReadBinary(reader, binary.BigEndian, &info.NameAndTypeIndex)
		return &info
	default:
		panic(fmt.Errorf("Invalid Tag '%v' for constant pool entry", tag))
	}
}

// Read a field_info struct
func readField(reader *bytes.Reader) FieldInfo {
	var fieldInfo FieldInfo
	safeReadBinary(reader, binary.BigEndian, &fieldInfo.AccessFlags)
	safeReadBinary(reader, binary.BigEndian, &fieldInfo.NameIndex)
	safeReadBinary(reader, binary.BigEndian, &fieldInfo.DescriptorIndex)
	safeReadBinary(reader, binary.BigEndian, &fieldInfo.AttrsCount)
	for i := uint16(0); i < fieldInfo.AttrsCount; i++ {
		fieldInfo.Attrs = append(fieldInfo.Attrs, readAttr(reader))
	}
	return fieldInfo
}

// Read an attribute_info struct
func readAttr(reader *bytes.Reader) AttrInfo {
	var attrInfo AttrInfo
	safeReadBinary(reader, binary.BigEndian, &attrInfo.NameIndex)
	safeReadBinary(reader, binary.BigEndian, &attrInfo.AttrLength)
	attrInfo.AttrData = make([]byte, attrInfo.AttrLength)
	io.ReadFull(reader, attrInfo.AttrData)
	return attrInfo
}

// Reads a method_info struct
func readMethod(reader *bytes.Reader) MethodInfo {
	var methodInfo MethodInfo
	safeReadBinary(reader, binary.BigEndian, &methodInfo.AccessFlags)
	safeReadBinary(reader, binary.BigEndian, &methodInfo.NameIndex)
	safeReadBinary(reader, binary.BigEndian, &methodInfo.DescriptorIndex)
	safeReadBinary(reader, binary.BigEndian, &methodInfo.AttrsCount)
	for i := uint16(0); i < methodInfo.AttrsCount; i++ {
		methodInfo.Attrs = append(methodInfo.Attrs, readAttr(reader))
	}
	return methodInfo
}

func safeReadBinary(reader io.Reader, order binary.ByteOrder, location interface{}) {
	if err := binary.Read(reader, order, location); err != nil {
		panic(err)
	}
}
