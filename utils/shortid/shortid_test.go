package shortid

import "testing"

func TestGenerate(t *testing.T) {
	generateID, err := GetDefault().Generate()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(generateID)
}
