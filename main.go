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
		// if there are no files to generate, return nil
		if len(gen.Files) == 0 {
			return nil
		}

		// check if the first file's GoImportPath is included in the includeTypesGoPackages list
		// if not, return nil
		if !checkGoImportPath(gen, gen.Files[0]) {
			return nil
		}

		// generate the bridge
		genGovmomiBridge(gen)

		// generate the govmomi types
		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}

			// generate enum bridge
			if len(f.Enums) > 0 {
				genEnumBridge(gen, f)
			}

			// generate message bridge
			if len(f.Messages) > 0 {
				genMessageBridge(gen, f)
			}
		}

		gen.SupportedFeatures = gengo.SupportedFeatures
		return nil
	})
}

// checkGoImportPath checks if the given file's GoImportPath is included in the includeTypesGoPackages list
func checkGoImportPath(gen *protogen.Plugin, file *protogen.File) bool {
	for _, p := range includeTypesGoPackages {
		if p == string(file.GoImportPath) {
			return true
		}
	}
	return false
}
