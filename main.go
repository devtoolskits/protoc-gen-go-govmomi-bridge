package main

import (
	"flag"
	"strings"

	gengo "google.golang.org/protobuf/cmd/protoc-gen-go/internal_gengo"
	"google.golang.org/protobuf/compiler/protogen"
)

type arrayFlags []string

func (i *arrayFlags) String() string {
	return strings.Join(*i, ",")
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var (
	includeEnumFiles       arrayFlags
	includeTypesGoPackages arrayFlags
)

func main() {
	var flags flag.FlagSet

	flags.Var(&includeEnumFiles, "include_enum_proto_files", "list of proto files with enum definitions to include")
	flags.Var(&includeTypesGoPackages, "include_types_go_packages", "list of go packages with govmomi-related types definitions to include")

	protogen.Options{
		ParamFunc: flags.Set,
	}.Run(func(gen *protogen.Plugin) error {
		if len(gen.Files) == 0 {
			return nil
		}

		// types bridge is a general purpose helper, we only need to generate it once for required package
		f := gen.Files[0]

		for _, p := range includeTypesGoPackages {
			if p == string(f.GoImportPath) {
				genGovmomiBridge(gen, f)
			}
		}

		// handle enums
		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}

			if len(f.Enums) > 0 {
				genEnumBridge(gen, f)
			}

			if len(f.Messages) > 0 {
				genMessageBridge(gen, f)
			}
		}

		gen.SupportedFeatures = gengo.SupportedFeatures
		return nil
	})
}
