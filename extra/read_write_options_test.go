package extra

import (
	"testing"

	"github.com/stretchr/testify/require"

	jsoniter "github.com/json-iterator/go"
)

func Test_read_only(t *testing.T) {
	SupportReadWriteOptions()

	should := require.New(t)

	type TestObject struct {
		Name string `json:"name,<-"`
		Age  int64  `json:"age"`
	}

	obj := TestObject{
		Name: "test",
	}

	data, err := jsoniter.Marshal(obj)

	should.NoError(err)

	should.Equal(`{"age":0}`, string(data))

	err = jsoniter.Unmarshal([]byte(`{"name":"read-only"}`), &obj)
	should.NoError(err)

	should.Equal("read-only", obj.Name)
}

func Test_write_only(t *testing.T) {
	SupportReadWriteOptions()

	should := require.New(t)

	type TestObject struct {
		Name string `json:"name,->"`
		Age  int64  `json:"age"`
	}

	obj := TestObject{
		Name: "test",
	}

	data, err := jsoniter.Marshal(obj)

	should.NoError(err)

	should.Equal(`{"name":"test","age":0}`, string(data))

	err = jsoniter.Unmarshal([]byte(`{"name":"write-only"}`), &obj)
	should.NoError(err)

	should.Equal("test", obj.Name)
}
