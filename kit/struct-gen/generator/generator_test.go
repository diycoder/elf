package generator

import "testing"

func TestDBConvert(t *testing.T) {
	gen := NewStructGen(
		WithPath("/Users/diycoder/work/person/elf/kit/struct-gen/model"),
		WithPackageName("model"),
		WithDsn("root:123456@tcp(localhost:3306)/tmp"),
		WithTable("service,service_error_code"),
		WithEnableJsonTag(false),
		WithEnableMsgPackTag(false),
		WithShowSQL(false),
		WithParseTime(false),
		WithTagKey("db"),
	)
	err := gen.Run()
	if err != nil {
		t.Error(err)
		return
	}
}
