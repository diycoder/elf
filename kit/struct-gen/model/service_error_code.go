package model

type ServiceErrorCode struct {
	ID         int64  `db:"id"`          // 错误码id
	ServiceID  string `db:"service_id"`  // 服务id
	ErrorCode  string `db:"error_code"`  // 错误码
	ErrorMsg   string `db:"error_msg"`   // 错误消息
	CreateTime string `db:"create_time"` // 创建时间
	UpdateTime string `db:"update_time"` // 更新时间
}
