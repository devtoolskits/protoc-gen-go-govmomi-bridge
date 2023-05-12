package main

import (
	"bytes"
	"path/filepath"
	"strings"
	"text/template"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/descriptorpb"
)

/*
genEnumBridge generates the bridge between the enum and the govmomi type, including the following functions:

- ToGovmomi: convert the enum to the govmomi type

- MustToGovmomi: convert the enum to the govmomi type, panic if the conversion fails

- FromGovmomi: convert the govmomi type to the enum
*/
func genEnumBridge(gen *protogen.Plugin, file *protogen.File) (*protogen.GeneratedFile, []*protogen.Enum) {
	if len(includeEnumFiles) > 0 {
		included := false
		for _, f := range includeEnumFiles {
			if f == *file.Proto.Name {
				included = true
				break
			}
		}
		if !included {
			return nil, nil
		}
	}

	filename := file.GeneratedFilenamePrefix + "_govmomi.pb.go"
	g := gen.NewGeneratedFile(filename, file.GoImportPath)

	g.P("// Code generated by protoc-gen-go-govmomi-bridge. DO NOT EDIT.")
	g.P()
	g.P("package ", file.GoPackageName)
	g.P()

	// handle import
	g.P("import (")
	g.P(`"errors"`)
	g.P(`vmomiTypes "github.com/vmware/govmomi/vim25/types"`)
	g.P(")")

	g.P()

	// add general enum undefined error type
	g.P("var ErrEnumUndefined = errors.New(\"enum undefined or unspecified\")")

	g.P()

	// provide a helper function to get the pointer of an enum
	g.P("// NewPointer returns a pointer to the given value.")
	g.P("func NewPointer[T any](v T) *T {")
	g.P("return &v")
	g.P("}")

	g.P()

	// provide a helper interface for all enums
	g.P("// Enum is the interface that all enums must implement.")
	g.P("type Enum interface {")
	g.P("FromGovmomi(string)")
	g.P("}")

	for _, enum := range file.Enums {
		name := enum.GoIdent.GoName

		if len(enum.Values) == 0 {
			continue
		}

		enumNameGovmomi := enum.Values[0].GoIdent.GoName[len(name)+1 : len(name)*2+1]

		g.P("// ToGovmomi converts the enum to the govmomi type, return ErrEnumUndefined if the conversion fails")
		g.P("func (x *", name, ") ToGovmomi() (*vmomiTypes.", enumNameGovmomi, ", error) {")
		g.P("switch *x {")
		for _, value := range enum.Values {

			if strings.Contains(value.GoIdent.GoName, "UNSPECIFIED") {
				continue
			}

			enumValue := value.GoIdent.GoName[len(name)+1:]

			g.P("case ", value.GoIdent.GoName, ":")
			g.P("return NewPointer(vmomiTypes.", enumValue, "), nil")
		}
		g.P("default:")
		g.P("return nil, ErrEnumUndefined")
		g.P("}")
		g.P("}")
		g.P()

		g.P("// MustToGovmomi converts the enum to the govmomi type, panic if the conversion fails")
		g.P("func (x *", name, ") MustToGovmomi() *vmomiTypes.", enumNameGovmomi, " {")
		g.P("v, err := x.ToGovmomi()")
		g.P("if err != nil {")
		g.P("return nil")
		g.P("}")
		g.P("return v")
		g.P("}")
		g.P()

		g.P("// FromGovmomi converts the govmomi type to the enum,")
		g.P("func (x *", name, ") FromGovmomi(v string) {")
		g.P("switch vmomiTypes.", enumNameGovmomi, "(v) {")
		for _, value := range enum.Values {

			if strings.Contains(value.GoIdent.GoName, "UNSPECIFIED") {
				continue
			}

			enumValue := value.GoIdent.GoName[len(name)+1:]

			g.P("case vmomiTypes.", enumValue, ":")
			g.P("*x = ", value.GoIdent.GoName)
		}
		g.P("default:")
		g.P("}")
		g.P("}")
		g.P()
	}

	return g, file.Enums
}

/*
genGovmomiBridge generates the bridge between messages and govmomi types, including the following functions:

- FromGovmomi: convert a govmomi struct to a message
*/
func genGovmomiBridge(gen *protogen.Plugin) *protogen.GeneratedFile {
	f := gen.Files[0]
	// handle messages
	t, err := template.New("types_govmomi_tamplate").Parse(TypesGovmomiTemplate)
	if err != nil {
		panic("fail to parse template:" + err.Error())
	}

	var buf bytes.Buffer

	type data struct {
		GoPackageName string
	}

	err = t.Execute(&buf, &data{
		GoPackageName: string(f.GoPackageName),
	})
	if err != nil {
		panic("fail to execute template:" + err.Error())
	}

	filename := filepath.Join(filepath.Dir(f.GeneratedFilenamePrefix), "govmomi_bridge.pb.go")
	g := gen.NewGeneratedFile(filename, f.GoImportPath)
	if _, err := g.Write(buf.Bytes()); err != nil {
		panic("fail to write file:" + err.Error())
	}

	return g
}

type Tag string

const (
	TagMessageGovmomiAlias   Tag = "govmomi_alias" // govmomi type alias
	TagFieldGovmomiFieldName Tag = "name"          // govmomi field name
	TagFieldGovmomiFieldType Tag = "type"          // govmomi field type
)

// genMessageBridge generates the bridge between messages and govmomi types, like alias
func genMessageBridge(gen *protogen.Plugin, file *protogen.File) {
	filename := file.GeneratedFilenamePrefix + "_govmomi.pb.go"
	g := gen.NewGeneratedFile(filename, file.GoImportPath)

	g.P("// Code generated by protoc-gen-go-govmomi-bridge. DO NOT EDIT.")
	g.P()
	g.P("package ", file.GoPackageName)
	g.P()

	for _, message := range file.Messages {
		parseMessageOptions(message, g)
	}
}

func parseMessageOptions(m *protogen.Message, g *protogen.GeneratedFile) map[Tag]string {
	res := make(map[Tag]string)

	opts := m.Desc.Options().(*descriptorpb.MessageOptions)
	if opts == nil {
		return res
	}

	// get string inside double quote
	// only one option is supported here
	aliasName := strings.Split(opts.String(), "\"")[1]
	g.P("type ", aliasName, " = ", m.GoIdent.GoName)
	g.P(" ")

	return res
}
