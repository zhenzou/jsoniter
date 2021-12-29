package extra

import (
	"strings"

	jsoniter "github.com/json-iterator/go"
)

// SupportReadWriteOptions extend jsoniter to support read-only or write-only
func SupportReadWriteOptions() {
	jsoniter.RegisterExtension(&readWriteOptions{})
}

// readWriteOptions extend jsoniter to support read-only or write-only
// <- read-only means can only decode
// -> write-only means can only encode
type readWriteOptions struct {
	jsoniter.DummyExtension
}

// UpdateStructDescriptor update descriptor if any field is read-only or write-only
// NOTE: for now ,the options must set by json tag
func (extension *readWriteOptions) UpdateStructDescriptor(descriptor *jsoniter.StructDescriptor) {
	for _, field := range descriptor.Fields {
		writeOnly := extension.isWriteOnly(field)
		readOnly := extension.isReadOnly(field)

		if writeOnly && readOnly {
			continue
		} else if writeOnly {
			extension.setWriteOnly(field)
		} else if readOnly {
			extension.setReadOnly(field)
		}
	}
}

func (extension *readWriteOptions) setWriteOnly(field *jsoniter.Binding) {
	field.FromNames = []string{}
}

func (extension *readWriteOptions) setReadOnly(field *jsoniter.Binding) {
	field.ToNames = []string{}
}

func (extension *readWriteOptions) isReadOnly(field *jsoniter.Binding) bool {
	return extension.hasJSONOption(field, "<-")
}

func (extension *readWriteOptions) isWriteOnly(field *jsoniter.Binding) bool {
	return extension.hasJSONOption(field, "->")
}

func (extension *readWriteOptions) hasJSONOption(field *jsoniter.Binding, opt string) bool {
	tag, ok := field.Field.Tag().Lookup("json")
	if !ok {
		return false
	}
	parts := strings.Split(tag, ",")

	for _, part := range parts {
		equal := strings.EqualFold(part, strings.TrimSpace(opt))
		if equal {
			return true
		}
	}
	return false
}
