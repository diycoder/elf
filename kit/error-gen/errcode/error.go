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
