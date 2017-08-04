[![GoDoc](https://godoc.org/github.com/a10y/classy?status.svg)](https://godoc.org/github.com/a10y/classy)

# classy: Classfile parsing libray and tool

## Installing

```
go get github.com/a10y/classy
```


## In Action

The following Java file:


```java
public class more {

    private final static String HIGHLANDER = "connor macleod";
    protected double THERE_CAN_BE_ONLY = 1.0;

    protected final String name(String... words) {
        return words[0];
    }

    protected native String fakeNative();
}
```

When compiled to `more.class`, when run through classy produces the following output (color not preserved)

```
Magic: 0xCAFEBABE (valid)
Major: 52
Minor: 0

ConstantPool: (26 entries)
  ├── 01: CONSTANT_Methodref
  │
  ├── 02: CONSTANT_Fieldref
  │
  ├── 03: CONSTANT_Class
  │		more
  ├── 04: CONSTANT_Class
  │		java/lang/Object
  ├── 05: CONSTANT_Utf8
  │		"HIGHLANDER"
  ├── 06: CONSTANT_Utf8
  │		"Ljava/lang/String;"
  ├── 07: CONSTANT_Utf8
  │		"ConstantValue"
  ├── 08: CONSTANT_String
  │		"connor macleod"
  ├── 09: CONSTANT_Utf8
  │		"THERE_CAN_BE_ONLY"
  ├── 10: CONSTANT_Utf8
  │		"D"
  ├── 11: CONSTANT_Utf8
  │		"<init>"
  ├── 12: CONSTANT_Utf8
  │		"()V"
  ├── 13: CONSTANT_Utf8
  │		"Code"
  ├── 14: CONSTANT_Utf8
  │		"LineNumberTable"
  ├── 15: CONSTANT_Utf8
  │		"name"
  ├── 16: CONSTANT_Utf8
  │		"([Ljava/lang/String;)Ljava/lang/String;"
  ├── 17: CONSTANT_Utf8
  │		"fakeNative"
  ├── 18: CONSTANT_Utf8
  │		"()Ljava/lang/String;"
  ├── 19: CONSTANT_Utf8
  │		"SourceFile"
  ├── 20: CONSTANT_Utf8
  │		"more.java"
  ├── 21: CONSTANT_NameAndType
  │
  ├── 22: CONSTANT_NameAndType
  │
  ├── 23: CONSTANT_Utf8
  │		"more"
  ├── 24: CONSTANT_Utf8
  │		"java/lang/Object"
  └── 25: CONSTANT_Utf8
  │		"connor macleod"

Methods: (3 entries)
  ├── public void <init>()
  ├── protected final java.lang.String name(java.lang.String[])
  └── protected native java.lang.String fakeNative()

Fields: (2 entries)
  ├── private static final java.lang.String HIGHLANDER
  └── protected double THERE_CAN_BE_ONLY

Attrs: (1 entries)
  └── SourceFile
```
