package extra

import (
	"testing"

	"github.com/stretchr/testify/require"

	jsoniter "github.com/json-iterator/go"
)

func Test_multipleKeysExtension_extractOptionValue(t *testing.T) {
	type args struct {
		tag string
		opt string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 bool
	}{
		{
			name: "without decode",
			args: args{
				tag: "origin",
				opt: "<",
			},
			want:  "",
			want1: false,
		},
		{
			name: "without decode",
			args: args{
				tag: "origin,<:fallback",
				opt: "<",
			},
			want:  "fallback",
			want1: true,
		},
		{
			name: "with decode keys",
			args: args{
				tag: "origin,<:fallback1 fallback2",
				opt: "<",
			},
			want:  "fallback1 fallback2",
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			extension := &multipleKeysExtension{}
			got, got1 := extension.extractOptionValue(tt.args.tag, tt.args.opt)
			if got != tt.want {
				t.Errorf("extractOptionValue() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("extractOptionValue() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func init() {
	SupportMultipleKeys()
}

func Test_decode_multiple_keys(t *testing.T) {

	should := require.New(t)

	t.Run("fallback 1 name", func(t *testing.T) {
		type TestObject struct {
			Name string `json:"name,<:first_name"`
			Age  int64  `json:"age"`
		}

		var obj TestObject

		err := jsoniter.Unmarshal([]byte(`{"name":"name"}`), &obj)
		should.NoError(err)
		should.Equal("name", obj.Name)

		err = jsoniter.Unmarshal([]byte(`{"first_name":"first_name"}`), &obj)
		should.NoError(err)
		should.Equal("first_name", obj.Name)
	})

	t.Run("fallback 2 names", func(t *testing.T) {
		type TestObject struct {
			Name string `json:"name,<:first_name legal_name"`
			Age  int64  `json:"age"`
		}

		var obj TestObject

		err := jsoniter.Unmarshal([]byte(`{"legal_name":"name"}`), &obj)
		should.NoError(err)
		should.Equal("name", obj.Name)

		err = jsoniter.Unmarshal([]byte(`{"first_name":"first_name"}`), &obj)
		should.NoError(err)
		should.Equal("first_name", obj.Name)
	})
}

func Test_encode_multiple_keys(t *testing.T) {

	should := require.New(t)

	t.Run("1 more keys", func(t *testing.T) {
		type TestObject struct {
			Name string `json:"name,>:first_name"`
			Age  int64  `json:"age"`
		}

		var obj TestObject

		obj.Name = "test"

		data, err := jsoniter.Marshal(obj)
		should.NoError(err)

		should.Equal(`{"name":"test","first_name":"test","age":0}`, string(data))

	})
}
