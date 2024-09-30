package hashid

import "testing"

func TestHashIDEncode(t *testing.T) {
	number := int64(100086)
	result, err := Encode(number)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(result)
}

func TestHashIDEncode2(t *testing.T) {
	hash := "9n81powl"
	number, err := Decode(hash)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(number)
}
