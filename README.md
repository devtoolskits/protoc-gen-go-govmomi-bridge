# protoc-gen-go-govmomi-bridge
protoc-gen-go-govmomi-bridge is a protoc go plugin to generate a compatibility layer between vSphere vim25 proto definitions and govmomi type definitions.

An example proto package is included in the `proto` directory. The generated files are included in the `gen` directory.

## Usage

1. Install the plugin

```bash
go install github.com/jiayinzhang-mint/protoc-gen-go-govmomi-bridge@latest
```

2. Include this plugin in your buf.gen.yaml
   
   Parameters:
    - `include_enum_proto_files`: A list of proto files with enum definitions to include.
    - `include_types_go_packages`: A list of go packages with govmomi-related types definitions to include.

```yaml
plugins:
  - name: go-govmomi-bridge
    out: gen
    opt:
      - paths=source_relative
      - include_enum_proto_files=proto/v1/enum.proto
      - include_types_go_packages=github.com/jiayinzhang-mint/protoc-gen-go-govmomi-bridge/fixture/gen/proto/v1
```

# Author
**Jiayin Zhang**

* <https://github.com/jiayinzhang-mint>


## License

Released under [Apache-2.0 License](https://github.com/jiayinzhang-mint/protoc-gen-go-govmomi-bridge/blob/main/LICENSE)