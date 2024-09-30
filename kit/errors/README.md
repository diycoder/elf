# 错误

### 错误码映射

以下可根据错误码生成工具自动生成

```go
package errors

const (
	Success      = 0
	Error        = 500
	MissingData  = 100001
	DataStatus   = 100002
	ParamIllegal = 100003
)

var errorMsg = map[int]string{
	Success:      "操作成功",
	Error:        "系统异常",
	MissingData:  "数据缺失",
	DataStatus:   "数据参数不正确，请勿非法操作",
	ParamIllegal: "参数传入不合法:[%s]",
}

func ErrorMsg(code int) string {
	return errorMsg[code]
}
```

### 使用示例

```go
package common

import (
	"fmt"
	"time"

	"github.com/diycoder/elf/kit/errors"
)

// 创建错误
func ErrorNew(code int) error {
	return errors.New(errors.ErrorMsg(code), code)
}

// 创建自定义格式错误
func Errorf(code int, args ...interface{}) error {
	return errors.New(fmt.Sprintf(errors.ErrorMsg(code), args...), code)
}
```