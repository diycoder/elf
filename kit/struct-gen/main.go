package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/diycoder/elf/kit/struct-gen/generator"
)

func main() {
	dsn := flag.String("dsn", "", "数据库dsn配置. 如：root:123456@tcp(127.0.0.1:3306)/test")
	file := flag.String("file", "", "保存路径. 如：/Users/diycoder/work")
	table := flag.String("table", "", "要迁移的表. 如：test")
	methodName := flag.String("method", "", "获取结构体对应的表名. 如：Test")
	packageName := flag.String("package_name", "model", "生成的struct包名. 如：model")
	tagKey := flag.String("tag_key", "db", "字段tag的key. 如：db、orm等")
	showSQL := flag.Bool("show_sql", false, "是否添加查询SQL.")
	version := flag.Bool("version", false, "版本号.")
	enableJsonTag := flag.Bool("enable_json_tag", false, "是否添加json的tag,默认false.")
	v := flag.Bool("v", false, "版本号.")
	h := flag.Bool("h", false, "帮助.")
	help := flag.Bool("help", false, "帮助.")

	// 开始
	flag.Parse()

	if *h || *help {
		flag.Usage()
		return
	}

	// 版本号
	if *version || *v {
		fmt.Println(fmt.Sprintf("\n version: %s\n %s\n using -h param for more help \n",
			generator.VERSION, generator.VersionText))
		return
	}

	gen := generator.NewStructGen(
		generator.WithPath(*file),
		generator.WithTable(*table),
		generator.WithPackageName(*packageName),
		generator.WithDsn(*dsn),
		generator.WithShowSQL(*showSQL),
		generator.WithTableNameMethod(*methodName),
		generator.WithEnableJsonTag(*enableJsonTag),
		generator.WithTagKey(*tagKey),
	)

	if err := gen.Run(); err != nil {
		log.Println(err.Error())
	}
}
