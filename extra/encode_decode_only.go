package extra

import (
	"strings"

	jsoniter "github.com/json-iterator/go"
)

// SupportEncodeOrDecodeOnly extend jsoniter to support read-only or write-only
func SupportEncodeOrDecodeOnly() {
	jsoniter.RegisterExtension(&encodeOrDecodeOnlyExtension{})
}

// encodeOrDecodeOnlyExtension extend jsoniter to support decode-only or encode-only by tag
type encodeOrDecodeOnlyExtension struct {
	jsoniter.DummyExtension
}

// UpdateStructDescriptor update descriptor if any field is read-only or write-only
// NOTE: for now ,the options must set by json tag
func (extension *encodeOrDecodeOnlyExtension) UpdateStructDescriptor(descriptor *jsoniter.StructDescriptor) {
	for _, field := range descriptor.Fields {
		encodeOnly := extension.isEncodeOnly(field)
		decodeOnly := extension.isDecodeOnly(field)

		if encodeOnly && decodeOnly {
			continue
		} else if encodeOnly {
			extension.setEncodeOnly(field)
		} else if decodeOnly {
			extension.setDecodeOnly(field)
		}
	}
}

func (extension *encodeOrDecodeOnlyExtension) setEncodeOnly(field *jsoniter.Binding) {
	field.FromNames = []string{}
}

func (extension *encodeOrDecodeOnlyExtension) setDecodeOnly(field *jsoniter.Binding) {
	field.ToNames = []string{}
}

func (extension *encodeOrDecodeOnlyExtension) isDecodeOnly(field *jsoniter.Binding) bool {
	return extension.hasJSONOption(field, "<-")
}

func (extension *encodeOrDecodeOnlyExtension) isEncodeOnly(field *jsoniter.Binding) bool {
	return extension.hasJSONOption(field, "->")
}

func (extension *encodeOrDecodeOnlyExtension) hasJSONOption(field *jsoniter.Binding, opt string) bool {
	tag, ok := field.Field.Tag().Lookup("json")
	if !ok {
		return false
	}
	parts := strings.Split(tag, ",")

	for _, part := range parts {
		equal := strings.EqualFold(strings.TrimSpace(part), opt)
		if equal {
			return true
		}
	}
	return false
}
