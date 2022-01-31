package extra

import (
	"testing"

	"github.com/stretchr/testify/require"

	jsoniter "github.com/json-iterator/go"
)

func Test_decode_only(t *testing.T) {
	SupportEncodeOrDecodeOnly()

	should := require.New(t)

	type TestObject struct {
		Name string `json:"name,<-"`
		Age  int64  `json:"age"`
	}

	t.Run("simple not encode", func(t *testing.T) {

		obj := TestObject{
			Name: "test",
		}

		data, err := jsoniter.Marshal(obj)
		should.NoError(err)
		should.Equal(`{"age":0}`, string(data))
	})

	t.Run("simple do decode", func(t *testing.T) {

		obj := TestObject{
			Name: "test",
		}

		err := jsoniter.Unmarshal([]byte(`{"name":"decode-only"}`), &obj)
		should.NoError(err)
		should.Equal("decode-only", obj.Name)
	})

	t.Run("struct not encode", func(t *testing.T) {

		var obj = struct {
			Obj TestObject `json:"obj,<-"`
		}{}

		obj.Obj.Name = "test"

		data, err := jsoniter.Marshal(obj)
		should.NoError(err)
		should.Equal(`{}`, string(data))
	})

	t.Run("struct do decode", func(t *testing.T) {
		var obj = struct {
			Obj TestObject `json:"obj,<-"`
		}{}

		obj.Obj.Name = "test"

		err := jsoniter.Unmarshal([]byte(`{"obj":{"name":"decode-only"}}`), &obj)
		should.NoError(err)
		should.Equal("decode-only", obj.Obj.Name)
	})
}

func Test_encode_only(t *testing.T) {
	SupportEncodeOrDecodeOnly()

	should := require.New(t)

	type TestObject struct {
		Name string `json:"name,->"`
		Age  int64  `json:"age"`
	}

	t.Run("simple", func(t *testing.T) {

		obj := TestObject{
			Name: "test",
		}

		data, err := jsoniter.Marshal(obj)

		should.NoError(err)

		should.Equal(`{"name":"test","age":0}`, string(data))

		err = jsoniter.Unmarshal([]byte(`{"name":"encode-only"}`), &obj)
		should.NoError(err)

		should.Equal("test", obj.Name)
	})

	t.Run("struct", func(t *testing.T) {

		var obj = struct {
			Obj TestObject `json:"obj,->"`
		}{}

		obj.Obj.Name = "test"

		data, err := jsoniter.Marshal(obj)
		should.NoError(err)

		should.Equal(`{"obj":{"name":"test","age":0}}`, string(data))

		err = jsoniter.Unmarshal([]byte(`{"obj":{"name":"encode-only","age":10}}`), &obj)
		should.NoError(err)

		should.Equal("test", obj.Obj.Name)
		should.Equal(int64(0), obj.Obj.Age)
	})

	t.Run("nested", func(t *testing.T) {

		var obj = struct {
			Obj TestObject `json:"obj"`
		}{}

		obj.Obj.Name = "test"

		data, err := jsoniter.Marshal(obj)
		should.NoError(err)

		should.Equal(`{"obj":{"name":"test","age":0}}`, string(data))

		err = jsoniter.Unmarshal([]byte(`{"obj":{"name":"encode-only","age":10}}`), &obj)
		should.NoError(err)

		should.Equal("test", obj.Obj.Name)
		should.Equal(int64(10), obj.Obj.Age)
	})
}
