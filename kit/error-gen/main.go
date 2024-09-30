package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/diycoder/elf/kit/error-gen/generator"
)

func main() {
	// 生成错误码
	infile := flag.String("i", "", "错误码配置文件路径. 如：/Users/diycoder/work/sample/error-gen/error.yaml")
	outfile := flag.String("o", "", "生成错误码文件目录. 如：/Users/diycoder/work/sample/error-gen/error.go")
	packageName := flag.String("n", "errcode", "生成错误码文件的包名. ")

	version := flag.Bool("version", false, "版本号.")
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

	gen := generator.NewErrorCodeGen(
		generator.WithPackageName(*packageName),
		generator.WithInputPath(*infile),
		generator.WithOutputPath(*outfile),
	)
	if err := gen.Run(); err != nil {
		log.Println(err)
	}
}
