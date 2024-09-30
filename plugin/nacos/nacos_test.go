package nacos

import (
	"testing"
)

func TestNaocsCfg(t *testing.T) {
	err := load(&Options{
		Address:     "127.0.0.1:8848",
		SecretKey:   "",
		AccessKey:   "",
		WatchConfig: "../../env/nacos_watch.yaml",
		NamespaceId: "5fe7e4db-424c-40b9-bb6f-ac893f049d24",
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("--default--a-->>", Get("default", "a").String(""))
}
