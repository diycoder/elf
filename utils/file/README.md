# csv导出插件

## 食用方式

> 通过对struct新增两个tag，其中`csv`标签表示需要导出当前字段的表头，`sort`标签表示当前字段在表头的第几个位置。

> 目前仅支持通过`slice[struct{}]`的形式进行csv导出。
```text
支持导出的数据类型： 
string, Bool
Int, Int8, Int16, Int32, Int64
Uint, Uint8, Uint16, Uint32, Uint64
Float32, Float64

```

### 1. 新建结构体
```go
type DepartmentUser struct {
	ID           int64    `csv:"编号" sort:"1"`
	Name         string   `csv:"名字" sort:"3"` 
	Department   string   `csv:"部门" sort:"2"`
	Online       bool     `csv:"是否在线" sort:"4"`
	Salary       float64  `csv:"薪水" sort:"8"`
}
```

> 备注：sort可以为乱序, 但是不能重复，但是排序仍然按照升序的方式进行排列，由上事例可知，最终排序为`编号`,`部门`,`名字`,`是否在线`,`薪水`。

### 2.使用步骤
```go
func main() {
    list := []DepartmentUser{
        {ID:1, Name:"产品经理xxx", Department: "产品设计部", Online: true, Salary: 888.8},
        {ID:2, Name:"技术开发xxx", Department: "技术开发部", Online: false, Salary: 999.9},
        {ID:3, Name:"测试运维xxx", Department: "测试韵味部", Online: true, Salary: 666.6},
    }   
    // 初始化一个csvHelper，
    csvHelper := NewCsvHelper(list, "部门用户列表")
    // 执行数据导出
    err := csvHelper.Export()
    if err != nil {
        fmt.Printf("export error: %v", err)
        return
    }
    
    // 导出的数据由于在本地存储了临时文件，所以需要移除
    defer func() {
        err = csvHelper.Remove()
        if err != nil {
            fmt.Printf("export error: %v", err)
            return
        }
    }()
}
```

## 错误对照
| 错误提示                                                        | 错误原因                |
|-------------------------------------------------------------|---------------------|
| the struct may have som problem                             | 结构体未定义csv标签         |
| the specified tag is not provided according to the standard | sort排序不存在           |
| just accept slice into this function                        | 导出的多行数据目前仅支持slice传入 |
| just accept struct into this function                       | 多行数据的每一行仅支持为struct  |
| tag for sort is not a number type                           | tag不为数字类型           |
| tag for sort is repeated                                    | sort排序值冲突           |
