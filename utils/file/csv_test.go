package file

import (
	"os"
	"testing"
)

func TestSaveCsv(t *testing.T) {
	header := []string{"用户id", "用户名称", "部门"}
	data := make([][]string, 0)
	data = append(data, []string{"1", "diycoder", "商业变现部"})
	csv := NewCSV(
		WithHeader(header),
		WithFileDir(os.TempDir()),
		WithFileName("department.csv"),
		WithData(data),
	)
	path, err := csv.Export()
	if err != nil {
		t.Errorf("export error: %v", err)
		return
	}
	defer func() {
		err = csv.Remove()
		if err != nil {
			t.Errorf("remove error: %v", err)
			return
		}
	}()
	t.Logf("path:%v", path)
}
