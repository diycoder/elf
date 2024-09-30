package generator

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

/**
生成的代码格式

package errcode

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

*/

type metadata struct {
	// map类型
	Errors map[string]*statusData `yml:"errors"`
}

type statusData struct {
	Code string `yml:"code"`
	Msg  string `yml:"msg"`
}

type Options struct {
	inputPath   string // 错误码配置文件路径
	outputPath  string // 生成错误码文件路径
	packageName string // 包名错误码文件包名
}

type Option func(options *Options)

func WithInputPath(path string) Option {
	return func(o *Options) {
		o.inputPath = path
	}
}

func WithOutputPath(outputPath string) Option {
	return func(o *Options) {
		o.outputPath = outputPath
	}
}

func WithPackageName(packageName string) Option {
	return func(o *Options) {
		o.packageName = packageName
	}
}

func NewErrorCodeGen(opts ...Option) *Options {
	options := Options{}
	for _, o := range opts {
		o(&options)
	}

	return &options
}

func (o *Options) Run() error {
	if err := checkPath(o.inputPath); err != nil {
		return err
	}
	metadata, err := getStatusDatas(o.inputPath)
	if err != nil {
		return err
	}
	if err := o.writeMapToFile(metadata.Errors); err != nil {
		fmt.Println("generate file ", o.outputPath, " fail", "error is:", err)
	}

	fmt.Println("generate file ", o.outputPath, " success")
	return nil
}

// output handle file println's error
func output(w io.Writer, a ...interface{}) {
	if _, err := fmt.Fprintln(w, a...); err != nil {
		panic(err)
	}
}

func checkPath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return errors.New("path is directory")
	}
	return nil
}

// WriteMapToFile StatusData map组装成文件gen.go
func (o *Options) writeMapToFile(m map[string]*statusData) error {
	fn := o.outputPath
	f, err := os.Create(fn)
	if err != nil {
		fmt.Printf("create map file error: %v\n", err)
		return err
	}
	defer func() {
		if err = f.Close(); err != nil {
			return
		}
		if e := recover(); e != nil {
			err = fmt.Errorf("Got a panic: %+v. ", e)
			return
		}
	}()

	w := bufio.NewWriter(f)

	// 生成包信息
	output(w, "package", o.packageName, "\n")

	// 生成code常量字段
	output(w, "/*")
	output(w, "*错误码")
	output(w, " */")
	output(w, "const (")
	codeSort := make([]int, 0)
	for _, v := range m {
		codeI, errs := strconv.Atoi(v.Code)
		if err != nil {
			return errs
		}
		codeSort = append(codeSort, codeI)
	}
	sort.Ints(codeSort)
	for _, cs := range codeSort {
		for k, v := range m {
			vcode, errs := strconv.Atoi(v.Code)
			if errs != nil {
				return errs
			}
			if cs == vcode {
				lineStr := fmt.Sprintf("	%s	=	%s", k, v.Code)
				output(w, lineStr)
			}
		}
	}
	output(w, ")\n")

	// 生成message map字段
	output(w, "var errorMsg = map[int]string{")

	for _, cs := range codeSort {
		for k, v := range m {
			vcode, err := strconv.Atoi(v.Code)
			if err != nil {
				return err
			}
			if cs == vcode {
				lineStr := fmt.Sprintf("	%s	:	\"%s\",", k, v.Msg)
				output(w, lineStr)
			}
		}
	}
	output(w, "}\n")

	// 生成根据code获取Message 函数
	output(w, "func ErrMsg(code int) string {")
	output(w, "	return errorMsg[code]")
	output(w, "}")
	if err = w.Flush(); err != nil {
		return err
	}

	// TODO: move out from here
	err = format(fn)

	return nil
}

func format(fn string) error {
	if err := exec.Command("gofmt", "-w", fn).Run(); err != nil {
		return err
	}
	return nil
}

// 读取Yaml配置文件,
// 并转换成StatusData对象列表  struct结构
func getStatusDatas(inputPath string) (*metadata, error) {
	// 获取绝对地址
	absFilename, _ := filepath.Abs(inputPath)
	yamlFile, err := ioutil.ReadFile(absFilename)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	metadata := new(metadata)
	err = yaml.Unmarshal(yamlFile, metadata)

	s := string(yamlFile)
	// 设置注释加到code字段后面。#装换成//
	tempStr := ""
	for _, lineStr := range strings.Split(s, "\n") {
		lineStr = strings.TrimSpace(lineStr)
		if lineStr == "" {
			continue
		}
		if strings.HasPrefix(lineStr, "#") {
			tempStr = strings.Replace(lineStr, "#", "//", 1)
		}

		for _, v := range metadata.Errors {
			if strings.Contains(lineStr, v.Code) {
				v.Code = v.Code + tempStr
				continue
			}
		}
	}

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return metadata, nil
}
