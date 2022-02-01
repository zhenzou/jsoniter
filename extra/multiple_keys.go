package extra

import (
	"strings"

	jsoniter "github.com/json-iterator/go"
)

// SupportMultipleKeys extend jsoniter to support encode to multiple keys and decode from multiple keys
// NOTE: if multiple keys exist in the same time, witch decode to the field is undefined
func SupportMultipleKeys() {
	jsoniter.RegisterExtension(&multipleKeysExtension{})
}

// multipleKeysExtension
type multipleKeysExtension struct {
	jsoniter.DummyExtension
}

// UpdateStructDescriptor update descriptor if any field set another keys
// NOTE: for now ,the options must set by json tag
func (extension *multipleKeysExtension) UpdateStructDescriptor(descriptor *jsoniter.StructDescriptor) {
	for _, field := range descriptor.Fields {
		decodeNames := extension.getDecodeNames(field)
		if len(decodeNames) > 0 {
			field.FromNames = append(field.FromNames, decodeNames...)
		}

		encodeNames := extension.getEncodeNames(field)
		if len(encodeNames) > 0 {
			field.ToNames = append(field.ToNames, encodeNames...)
		}
	}
}

func (extension *multipleKeysExtension) getDecodeNames(field *jsoniter.Binding) []string {
	options, ok := extension.getJSONOptionValue(field, "<")
	if !ok || options == "" {
		return nil
	}
	return strings.Split(options, " ")
}

func (extension *multipleKeysExtension) getEncodeNames(field *jsoniter.Binding) []string {
	options, ok := extension.getJSONOptionValue(field, ">")
	if !ok || options == "" {
		return nil
	}
	return strings.Split(options, " ")
}

func (extension *multipleKeysExtension) getJSONOptionValue(field *jsoniter.Binding, opt string) (string, bool) {
	tag, ok := field.Field.Tag().Lookup("json")
	if !ok {
		return "", false
	}
	return extension.extractOptionValue(tag, opt)
}

func (extension *multipleKeysExtension) extractOptionValue(tag string, opt string) (string, bool) {
	parts := strings.Split(tag, ",")
	for _, part := range parts {
		part := strings.TrimSpace(part)
		if strings.HasPrefix(part, opt) {
			strs := strings.Split(part, ":")
			if len(strs) > 1 {
				return strs[1], true
			}
		}
	}
	return "", false
}
