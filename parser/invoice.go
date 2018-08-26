package parser

import (
	"time"
)

// PaymentState 表示账单的付款状态
type PaymentState int

const (
	// NeedPay 未付款
	NeedPay PaymentState = iota
	// FinishedPay 已付款
	FinishedPay
)

// Invoice 账单信息
type Invoice struct {
	// 账单编号
	Number string
	// 账单链接
	Link string
	// 账单开始日期
	StartDate time.Time
	// 账单结束日期
	ExpireDate time.Time
	// 支付金额
	Payment int64
	// 付款状态
	State PaymentState
}
