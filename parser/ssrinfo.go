package parser

// SSRInfo ssr套餐信息
type SSRInfo struct {
	*Service

	// 节点的端口和密码
	Port   int64
	Passwd string

	// 可用数据总量
	TotalData string
	// 已用数据总量
	UsedData string
	// 下载用量
	Download string
	// 上传用量
	Upload string

	// 可用节点信息
	Nodes []*SSRNode
}

// NewSSRInfo 生成SSRInfo
func NewSSRInfo(ser *Service) *SSRInfo {
	s := new(SSRInfo)
	s.Service = ser
	s.Nodes = make([]*SSRNode, 0)
	return s
}
