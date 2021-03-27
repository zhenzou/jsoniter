package test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	jsoniter "github.com/json-iterator/go"
)

func Test_readonly(t *testing.T) {
	type User struct {
		Name  string `json:"name"`
		Pwd   string `json:"pwd,read"`
		Extra struct {
			Age int64 `json:"age"`
		} `json:"extra,read"`
	}
	var user User

	json := jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal([]byte(`{"name":"test","pwd":"123456","extra":{"age":10}}`), &user)

	assert.NoError(t, err)
	assert.Equal(t, "test", user.Name)
	assert.Equal(t, "123456", user.Pwd)
	assert.Equal(t, int64(10), user.Extra.Age)

	data, err := json.Marshal(user)
	assert.NoError(t, err)

	assert.Equal(t, `{"name":"test"}`, string(data))
}

func Test_writeonly(t *testing.T) {
	type User struct {
		Name  string `json:"name"`
		Pwd   string `json:"pwd,write"`
		Extra struct {
			Age int64 `json:"age"`
		} `json:"extra,write"`
	}
	var user User

	json := jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal([]byte(`{"name":"test","pwd":"123456","extra":{"age":10}}`), &user)

	assert.NoError(t, err)
	assert.Equal(t, "test", user.Name)
	assert.Equal(t, "", user.Pwd)
	assert.Equal(t, int64(0), user.Extra.Age)

	user.Pwd = "123456"
	user.Extra.Age = 10
	data, err := json.Marshal(user)
	assert.NoError(t, err)

	assert.Equal(t, `{"name":"test","pwd":"123456","extra":{"age":10}}`, string(data))
}
