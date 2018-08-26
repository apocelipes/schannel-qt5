package parser

import (
	"time"
)

// Service 购买的服务信息
type Service struct {
	Name    string
	Link    string
	Price   string
	Expires time.Time
	State   string
}
