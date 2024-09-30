## 错误码自动生成工具

### 功能

根据配置文件`error.yml`，自动生成`error.go`错误码文件，错误码生成如下

### 安装

```bazaar
go install github.com/diycoder/elf/kit/error-gen@dev
```

### 使用

```bazaar
Usage of error-gen:
  -h	帮助.
  -help
    	帮助.
  -i string
    	错误码配置文件路径. 如：/Users/diycoder/work/sample/error-gen/error.yml
  -n string
    	生成错误码文件的包名.  (default "errcode")
  -o string
    	生成错误码文件目录. 如：/Users/diycoder/work/sample/error-gen/error.go
  -v	版本号.
  -version
    	版本号.
```

`error.yaml` 错误码配置文件内容

```bazaar
errors:
  MissingData:
    code: 100001
    msg: 数据缺失
  ParamIllegal:
    code: 100002
    msg: 参数传入不合法:[%s]

```

命令执行

```bazaar
error-gen -i="/Users/diycoder/work/sample/error-gen/error.yml" -o="/Users/diycoder/work/sample/error-gen/error.go" -n="errcode"
```

生成错误码文件内容

```bazaar
package errcode

/*
*错误码
 */
const (
	MissingData  = 100001
	ParamIllegal = 100002
)

var errorMsg = map[int]string{
	MissingData:  "数据缺失",
	ParamIllegal: "参数传入不合法:[%s]",
}

func errMsg(code int) string {
	return errorMsg[code]
}

```
