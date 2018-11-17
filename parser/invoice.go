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

// GetStatus 返回账单的状态
// 未付款会返回false
// time.Now()超过ExpireDate将视为账单过期
func (i *Invoice) GetStatus() (string, bool) {
	msg := ""
	flag := false

	if i.State == NeedPay {
		msg += "需要付款"
	} else if i.State == FinishedPay {
		msg += "无需付款"
		flag = true
	}

	current := GetCurrentDay()
	if current.After(i.ExpireDate) {
		msg += "，账单过期"
	}

	return msg, flag
}

// GetCurrentDay 返回当前的时间，精确到day
func GetCurrentDay() time.Time {
	timeFormat := "2006-01-02"
	dayText := time.Now().Format(timeFormat)
	day, _ := time.ParseInLocation(timeFormat, dayText, time.Local)
	return day
}
