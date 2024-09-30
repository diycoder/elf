package generator

import "testing"

func TestErrCodeGen(t *testing.T) {
	gen := NewErrorCodeGen(
		WithPackageName("errcode"),
		WithInputPath("/Users/diycoder/work/person/elf/kit/error-gen/generator/error.yaml"),
		WithOutputPath("/Users/diycoder/work/person/elf/kit/error-gen/errcode/error.go"),
	)
	if err := gen.Run(); err != nil {
		t.Error(err)
		return
	}
}
