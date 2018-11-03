package parser

import (
	"time"
)

// Service 购买的服务信息
type Service struct {
	// 服务名称
	Name    string
	// 服务详细信息链接
	Link    string
	// 服务价格
	Price   string
	// 服务过期时间
	Expires time.Time
	// 服务状态：是否可用/是否需要付费
	State   string
}
