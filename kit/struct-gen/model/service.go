package model

type Service struct {
	ID          int64  `db:"id"`           // 服务id
	ServiceName string `db:"service_name"` // 服务名称
	Comment     string `db:"comment"`      // 服务备注
	CreateTime  string `db:"create_time"`  // 创建时间
	UpdateTime  string `db:"update_time"`  // 更新时间
}
