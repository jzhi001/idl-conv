package code_gen

import (
	"errors"
	"fmt"
	"strings"
)

const (
	Star  = "*"
	Slice = "[]"
)

var builtinTypeMapping = map[Token]string{
	"float32": "float",
	"float64": "float",
	"int":     "int32",
	"int16":   "int32",
	"int32":   "int32",
	"int64":   "int64",
	"int8":    "int32",
	"rune":    "rune",
	"string":  "string",
	"uint":    "uint32",
	"uint16":  "uint32",
	"uint32":  "uint32",
	"uint64":  "uint64",
	"uint8":   "uint32",
	//"uintptr": "uintptr",
}

type FieldDesc struct {
	OrigFName       Token
	FName           Token
	FType           Token
	IsPointer       bool
	IsSlice         bool
	IsPrimitive     bool
	IsMap           bool   // TODO support map<A, B>
	BackTickComment string // TODO support `protobuf:"dog_name"`
	SlashedComment  string
}

func (f *FieldDesc) String() string {
	s := fmt.Sprintf("%s %s `json:\"%s\"`", f.FName.String(), f.FType.String(), Camel2Snake(f.FName.String()))
	if len(f.SlashedComment) > 0 {
		s += " // " + f.SlashedComment
	}
	return s
}

func (f *FieldDesc) ProtobufVer(idx int) string {

	if f.IsPrimitive {
		f.FType = Token(builtinTypeMapping[f.FType])
	}

	s := fmt.Sprintf("%s %s = %d;", f.FType, Camel2Snake(f.FName.String()), idx)
	if f.IsSlice {
		s = "repeated " + s
	}
	return s
}

type StructDesc struct {
	TName        Token
	Fields       []*FieldDesc
	dependencies []FieldDesc
}

func (t *StructDesc) String() string {
	s := fmt.Sprintf("type %s struct {\n", t.TName)

	for _, f := range t.Fields {
		s += f.String() + "\n"
	}

	s += "}"
	return s
}

func (t *StructDesc) ProtobufVer() string {
	s := fmt.Sprintf("message %s{\n", t.TName)

	for i, f := range t.Fields {
		s += strings.Repeat(Space, 4) + f.ProtobufVer(i+1) + NewLine
	}

	s += "}"
	return s
}

func (t *StructDesc) GetField(fieldName string) (*FieldDesc, error) {

	for _, f := range t.Fields {
		if f.FName.String() == fieldName {
			return f, nil
		}
	}

	return nil, errors.New("no such field")
}

func parseField(fName, fType Token) *FieldDesc {

	f := &FieldDesc{
		FName:     fName,
		FType:     fType,
		OrigFName: fName,
	}

	if f.FType.StartsWith(Slice) {
		f.IsSlice = true
		f.FType = f.FType[2:]
	}

	if f.FType.StartsWith(Star) {
		f.IsPointer = true
		f.FType = f.FType[1:]
	}

	if _, found := builtinTypeMapping[f.FType]; found {
		f.IsPrimitive = true
	}

	return f
}

func Parse(tokens []Token) (typeDescList []*StructDesc, err error) {

	it := NewTokenIterator(tokens)

	var tk Token

	for tk, err = it.NextNonNewLine(); it.HasNext(); tk, err = it.NextNonNewLine() {

		if err != nil {
			return nil, err
		}

		if tk != Type {
			return nil, errors.New("expect 'type' keyword, but got " + tk.Quote())
		}

		var typeName, fieldName, fieldType Token

		typeName, err = it.NextNonNewLine()
		_, err = it.NextNonNewLine() // skip 'struct'
		_, err = it.NextNonNewLine() // skip '{'

		var fieldDescList []*FieldDesc

		for fieldName, err = it.NextNonNewLine(); it.HasNext() && fieldName != RtCurlBrace; fieldName, err = it.NextNonNewLine() {

			if err != nil {
				return nil, err
			}

			fieldType, err = it.NextNonNewLine()
			if err != nil {
				return nil, err
			}

			fieldDescList = append(fieldDescList, parseField(fieldName, fieldType))
		}

		typeDescList = append(typeDescList, &StructDesc{
			TName:  typeName,
			Fields: fieldDescList,
		})
	}

	return typeDescList, nil
}
