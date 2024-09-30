package json

import (
	jsoniter "github.com/json-iterator/go"
)

var jsoniters = jsoniter.ConfigCompatibleWithStandardLibrary

func Marshal(data interface{}) ([]byte, error) {
	return jsoniters.Marshal(data)
}

func Unmarshal(data []byte, v interface{}) error {
	return jsoniters.Unmarshal(data, v)
}

func MarshalToString(data interface{}) (string, error) {
	return jsoniters.MarshalToString(data)
}

func UnmarshalFromString(str string, data interface{}) error {
	return jsoniters.UnmarshalFromString(str, data)
}
