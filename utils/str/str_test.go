package str

import "testing"

func TestStringToByte(t *testing.T) {
	s := "hello world"
	b := StringToByte(s)
	convertStr := ByteToString(b)
	t.Log(convertStr == s)
}
