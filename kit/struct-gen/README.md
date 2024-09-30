## MySQL 表转 struct 生成工具

### 功能

根据配置自动将MySQL 表结构转成 struct 的工具

### 安装
```bazaar
go install github.com/diycoder/elf/kit/struct-gen@dev
```

### 使用

```bazaar
Usage of struct-gen:
  -dsn string
    	数据库dsn配置. 如：root:123456@tcp(127.0.0.1:3306)/test
  -enable_json_tag
    	是否添加json的tag,默认false.
  -file string
    	保存路径. 如：/Users/diycoder/work
  -h	帮助.
  -help
    	帮助.
  -method string
    	获取结构体对应的表名. 如：Test
  -package_name string
    	生成的struct包名. 如：model (default "model")
  -show_sql
    	是否添加查询SQL.
  -table string
    	要迁移的表. 如：test
  -tag_key string
    	字段tag的key. 如：db、orm等 (default "db")
  -v	版本号.
  -version
    	版本号.
```

```bazaar
struct-gen -file="/Users/diycoder/work/sample/dbconvert/model" -dsn="root:123456@tcp(127.0.0.1:3306)/blog_center"
```