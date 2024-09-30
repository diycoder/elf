package generator

import (
	"bufio"
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/go-sql-driver/mysql"
)

// map for converting mysql type to golang types
var dataTypeMap = map[string]string{
	"int":                "int64",
	"integer":            "int64",
	"tinyint":            "int64",
	"smallint":           "int64",
	"mediumint":          "int64",
	"bigint":             "int64",
	"int unsigned":       "int64",
	"integer unsigned":   "int64",
	"tinyint unsigned":   "int64",
	"smallint unsigned":  "int64",
	"mediumint unsigned": "int64",
	"bigint unsigned":    "int64",
	"bit":                "int64",
	"bool":               "bool",
	"enum":               "string",
	"set":                "string",
	"varchar":            "string",
	"char":               "string",
	"tinytext":           "string",
	"mediumtext":         "string",
	"text":               "string",
	"longtext":           "string",
	"blob":               "string",
	"tinyblob":           "string",
	"mediumblob":         "string",
	"longblob":           "string",
	"date":               "string",
	"datetime":           "string",
	"timestamp":          "string",
	"time":               "string",
	"float":              "float64",
	"double":             "float64",
	"decimal":            "float64",
	"binary":             "string",
	"varbinary":          "string",
}

var timeTypeMap = map[string]string{
	"date":      "time.Time", // time.Time
	"datetime":  "time.Time", // time.Time
	"timestamp": "time.Time", // time.Time
	"time":      "time.Time", // time.Time
}

type Options struct {
	dsn              string
	savePath         string
	db               *sql.DB
	table            []string
	tableNameMethod  string // function name
	showSQL          bool   // generate select SQL, default false
	parseTime        bool   // parse mysql time, default false
	enableJsonTag    bool   // enable json tag, default false
	enableMsgpackTag bool   // enable msgpack tag, default false
	packageName      string // the packege name
	tagKey           string // the struct default tag
}

type Option func(options *Options)

func WithDsn(dsn string) Option {
	return func(o *Options) {
		o.dsn = dsn
	}
}

func WithTagKey(tagKey string) Option {
	return func(o *Options) {
		o.tagKey = tagKey
	}
}

func WithShowSQL(showSQL bool) Option {
	return func(o *Options) {
		o.showSQL = showSQL
	}
}

func WithParseTime(parseTime bool) Option {
	return func(o *Options) {
		o.parseTime = parseTime
	}
}

func WithPackageName(packageName string) Option {
	return func(o *Options) {
		o.packageName = packageName
	}
}

func WithPath(path string) Option {
	return func(o *Options) {
		o.savePath = path
	}
}

func WithTable(table string) Option {
	return func(o *Options) {
		o.table = strings.Split(table, ",")
	}
}

func WithEnableJsonTag(enableJsonTag bool) Option {
	return func(o *Options) {
		o.enableJsonTag = enableJsonTag
	}
}

func WithEnableMsgPackTag(enableMsgpackTag bool) Option {
	return func(o *Options) {
		o.enableMsgpackTag = enableMsgpackTag
	}
}

func WithTableNameMethod(tableNameMethod string) Option {
	return func(o *Options) {
		o.tableNameMethod = tableNameMethod
	}
}

func NewStructGen(opts ...Option) *Options {
	options := Options{}
	for _, o := range opts {
		o(&options)
	}
	return &options
}

func (o *Options) Run() error {
	// check dsn
	if o.dsn == "" {
		return errors.New("dsn is empty")
	}
	dsn, err := mysql.ParseDSN(o.dsn)
	if err != nil {
		fmt.Printf("parse dsn:%v error: %v\n", o.dsn, err)
		return err
	}
	db, err := sql.Open("mysql", dsn.FormatDSN())
	if err != nil {
		fmt.Printf("get mysql error: %v\n", err)
		return err
	}
	o.db = db

	if err := checkPath(o.savePath); err != nil {
		return err
	}

	// get schema
	tableColumns, err := o.getColumns(o.table...)
	if err != nil {
		fmt.Printf("get columns error: %v\n", err)
		return err
	}

	// generate struct
	for tableName, columns := range tableColumns {
		if err := o.writeToFile(tableName, o.packageName, columns); err != nil {
			fmt.Printf("write to file:%v error: %v\n", tableName, err)
		}
	}
	return nil
}

func (o *Options) writeToFile(tableName, packageName string, column []column) error {
	file := fmt.Sprintf("%s%s%s.go", o.savePath, string(os.PathSeparator), tableName)
	f, err := os.Create(file)
	if err != nil {
		fmt.Printf("create file:%v error: %v\n", file, err)
		return err
	}
	defer func() {
		if err = f.Close(); err != nil {
			fmt.Printf("file:%v close error: %v\n", file, err)
			return
		}
		if e := recover(); e != nil {
			err = fmt.Errorf("Got a panic: %+v. ", e)
			return
		}
	}()
	w := bufio.NewWriter(f)
	structName := caseToCamel(tableName)
	depth := 1

	// general package and import
	o.generalPackageAndImport(w, packageName)
	// general select SQL
	o.generalSQL(column, w)
	// general struct
	o.generalStruct(w, structName, column, depth)
	// general function
	o.generalFunction(w, structName, tableName, depth)

	if err = w.Flush(); err != nil {
		fmt.Printf("file:%v flush error: %v\n", file, err)
		return err
	}

	// format go file
	if err = exec.Command("gofmt", "-w", file).Run(); err != nil {
		return err
	}
	fmt.Println(fmt.Sprintf("generate file %v success", file))
	return nil
}

func (o *Options) generalPackageAndImport(w *bufio.Writer, packageName string) {
	// general packege
	output(w, "package", packageName)

	// parse time
	if o.parseTime && o.enableMsgpackTag {
		output(w, `import (`)

		output(w, `"time"`, "\n")
		output(w, `"github.com/shamaton/msgpack/v2"`)

		output(w, `)`)
	} else if o.parseTime {
		output(w, `import "time"`)
	}
}

func (o *Options) generalStruct(w *bufio.Writer, structName string, column []column, depth int) {
	output(w, "type", structName, " struct {")
	for _, v := range column {
		comment := ""
		if v.tableColumn.ColumnComment != "" {
			comment += " // " + v.tableColumn.ColumnComment
		}
		output(w, tab(depth), v.structColumn.Name, " ", v.structColumn.DataType, " ", v.structColumn.Tag, comment)
	}
	output(w, tab(depth-1), "}", "\n")
}

func (o *Options) generalFunction(w *bufio.Writer, structName, tableName string, depth int) {
	// get table name method
	if o.tableNameMethod != "" {
		output(w, fmt.Sprintf("func (*%s) %s() string {", structName, o.tableNameMethod))
		output(w, fmt.Sprintf("%sreturn \"%s\"", tab(depth), tableName))
		output(w, "}", "\n")
	}
	// enable msgpack
	if o.enableMsgpackTag {
		o.generateMarshalMethod(structName, w)
		o.generateUnmarshalMethod(structName, w)

		structListName := o.getStructListName(structName)
		output(w, fmt.Sprintf(`type %s []*%s`, structListName, structName))
		o.generateMarshalMethod(structListName, w)
		o.generateUnmarshalMethod(structListName, w)
	}
}

func (o *Options) generateMarshalMethod(structName string, w *bufio.Writer) {
	firstWord := o.getFirstWord(structName)
	output(w, fmt.Sprintf(`func (%s *%s) MarshalBinary() ([]byte, error) {`, firstWord, structName))
	output(w, fmt.Sprintf(`return msgpack.Marshal(%s)`, firstWord))
	output(w, `}`, "\n")
}

func (o *Options) generateUnmarshalMethod(structName string, w *bufio.Writer) {
	firstWord := o.getFirstWord(structName)
	output(w, fmt.Sprintf(`func (%s *%s) UnmarshalBinary(data []byte) error {`, firstWord, structName))
	output(w, fmt.Sprintf(`return msgpack.Unmarshal(data, %s)`, firstWord))
	output(w, `}`, "\n")
}

func (o *Options) getStructListName(structName string) string {
	return fmt.Sprintf("%sList", structName)
}

func (o *Options) getFirstWord(structName string) string {
	lower := strings.ToLower(structName)
	return lower[:1]
}

func (o *Options) generalSQL(column []column, w *bufio.Writer) {
	if !o.showSQL {
		return
	}
	columns := "select "
	for i, v := range column {
		if v.tableColumn.ColumnNullable == "YES" {
			switch dataTypeMap[v.tableColumn.ColumnDataType] {
			case "int64", "float64":
				columns += "ifnull(" + v.tableColumn.ColumnName + ",0) as " + "`" + v.tableColumn.ColumnName + "`,"
			case "string":
				columns += "ifnull(" + v.tableColumn.ColumnName + ",'') as " + "`" + v.tableColumn.ColumnName + "`,"
			case "bool":
				columns += "ifnull(" + v.tableColumn.ColumnName + ",false) as " + "`" + v.tableColumn.ColumnName + "`,"
			}
		} else {
			columns += "`" + v.tableColumn.ColumnName + "`,"
		}
		if i == len(column)-1 {
			columns = columns[:len(columns)-1]
			columns += " from " + "`" + v.tableColumn.TableName + "`"
		}
	}
	output(w, " // ", columns)
}

// output handle file println error
func output(w io.Writer, a ...interface{}) {
	if _, err := fmt.Fprintln(w, a...); err != nil {
		panic(err)
	}
}

type column struct {
	structColumn structColumn
	tableColumn  tableColumn
}

type structColumn struct {
	Name     string // struct name
	DataType string // struct data type
	Tag      string // struct tag
}

type tableColumn struct {
	TableName      string // table name
	ColumnName     string // column name
	ColumnDataType string // cloumn data type
	ColumnNullable string // column is empty or not
	ColumnComment  string // column comment
}

// Function for fetching schema definition of passed table
func (o *Options) getColumns(tables ...string) (map[string][]column, error) {
	// sql
	query := `SELECT COLUMN_NAME,DATA_TYPE,IS_NULLABLE,TABLE_NAME,COLUMN_COMMENT FROM information_schema.COLUMNS WHERE table_schema = DATABASE()`
	// include table
	tableNames := make([]string, 0)
	for _, table := range tables {
		if table == "" {
			continue
		}
		tableNames = append(tableNames, fmt.Sprintf(`'%s'`, table))
	}
	if len(tableNames) > 0 {
		query += fmt.Sprintf(` AND TABLE_NAME IN (%s)`, strings.Join(tableNames, ","))
	}
	// sql order
	query += ` ORDER BY TABLE_NAME ASC, ORDINAL_POSITION ASC`
	rows, err := o.db.Query(query)
	if err != nil {
		log.Println("Error reading table information: ", err.Error())
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println("rows close err: ", err.Error())
		}
	}(rows)
	return o.parseColumn(rows)
}

func (o *Options) parseColumn(rows *sql.Rows) (map[string][]column, error) {
	tableColumns := make(map[string][]column)
	for rows.Next() {
		col := tableColumn{}
		err := rows.Scan(&col.ColumnName, &col.ColumnDataType, &col.ColumnNullable, &col.TableName, &col.ColumnComment)
		if err != nil {
			log.Println(err.Error())
			return nil, err
		}

		structInfo := structColumn{}
		structInfo.DataType = dataTypeMap[col.ColumnDataType]
		if o.parseTime {
			if val, ok := timeTypeMap[col.ColumnDataType]; ok {
				structInfo.DataType = val
			}
		}
		structInfo.Name = getStructColumnName(col.ColumnName)
		jsonTag, msgpackTag := col.ColumnName, col.ColumnName

		if o.tagKey == "" {
			o.tagKey = "db"
		}
		buf := bytes.Buffer{}
		buf.WriteString("`")
		buf.WriteString(fmt.Sprintf("%s:\"%s\"", o.tagKey, col.ColumnName))
		if o.enableJsonTag {
			buf.WriteString(fmt.Sprintf(" json:\"%s\"", jsonTag))
		}
		if o.enableMsgpackTag {
			buf.WriteString(fmt.Sprintf(" msgpack:\"%s\"", msgpackTag))
		}
		buf.WriteString("`")
		structInfo.Tag = buf.String()

		if _, ok := tableColumns[col.TableName]; !ok {
			tableColumns[col.TableName] = []column{}
		}

		columns := column{}
		columns.structColumn = structInfo
		columns.tableColumn = col
		tableColumns[col.TableName] = append(tableColumns[col.TableName], columns)
	}
	return tableColumns, nil
}

func checkPath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return err
	}
	if !s.IsDir() {
		return errors.New("path is not directory")
	}
	return nil
}

func tab(depth int) string {
	return strings.Repeat("\t", depth)
}

func getStructColumnName(name string) string {
	return replace(caseToCamel(name))
}

func caseToCamel(name string) string {
	name = strings.Replace(name, "_", " ", -1)
	name = strings.Title(name)
	return strings.Replace(name, " ", "", -1)
}

func replace(name string) string {
	if len(name) == 2 {
		return strings.ToUpper(name)
	}
	if strings.HasSuffix(name, "id") || strings.HasSuffix(name, "Id") {
		return name[:len(name)-2] + "ID"
	}
	return name
}
