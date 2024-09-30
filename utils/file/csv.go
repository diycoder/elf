package file

import (
	"encoding/csv"
	"io/ioutil"
	"os"
	"strings"
)

type CSV struct {
	header   []string
	data     [][]string
	dir      string
	fileName string
}

type Option func(opt *CSV)

func NewCSV(opts ...Option) *CSV {
	opt := &CSV{}
	for _, o := range opts {
		o(opt)
	}
	return opt
}

func WithHeader(header []string) Option {
	return func(opt *CSV) {
		opt.header = header
	}
}

func WithData(data [][]string) Option {
	return func(opt *CSV) {
		opt.data = data
	}
}

func WithFileDir(fileDir string) Option {
	return func(opt *CSV) {
		opt.dir = fileDir
	}
}

func WithFileName(fileName string) Option {
	return func(opt *CSV) {
		opt.fileName = fileName
	}
}

// Export 保存到csv文件
func (c *CSV) Export() (string, error) {
	savePath := c.getFileUrl()
	var f *os.File
	f, err := os.Create(savePath)
	if err != nil {
		return "", err
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)
	_, err = f.WriteString("\xEF\xBB\xBF")
	if err != nil {
		return "", err
	}
	w := csv.NewWriter(f)
	err = w.Write(c.header)
	if err != nil {
		return "", err
	}
	err = w.WriteAll(c.data)
	if err != nil {
		return "", err
	}
	w.Flush()
	return savePath, nil
}

func (c *CSV) getFileUrl() string {
	separator := string(os.PathSeparator)
	if strings.HasSuffix(c.dir, separator) {
		return c.dir + c.fileName
	}
	return c.dir + string(os.PathSeparator) + c.fileName
}

// Remove 删除文件
func (c *CSV) Remove() error {
	err := os.Remove(c.getFileUrl())
	if err != nil {
		return err
	}
	return nil
}

// ReadCsvFile 按行读取文件
func ReadCsvFile(path string) ([][]string, error) {
	var records [][]string
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	reader := csv.NewReader(file)
	records, err = reader.ReadAll()
	if err != nil {
		return nil, err
	}
	return records, nil
}

// IsDir 所给路径是否为文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// IsFile 所给路径是否为文件
func IsFile(path string) bool {
	return !IsDir(path)
}

// Remove 删除文件
func Remove(path string) (bool, error) {
	err := os.Remove(path)
	if err != nil {
		return false, err
	}
	return true, nil
}

// Copy 文件复制操作
func Copy(oldPath, newPath string) (bool, error) {
	data, err := ioutil.ReadFile(oldPath)
	if err != nil {
		return false, err
	}
	err = ioutil.WriteFile(newPath, data, 0666)
	if err != nil {
		return false, err
	}
	return true, nil
}
