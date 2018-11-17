package parser

import (
	"testing"

	"time"
)

func TestInvoiceGetStatus(t *testing.T) {
	// 与过期时间比较
	now := GetCurrentDay()

	testData := []*struct {
		// 账单信息
		i *Invoice
		// err string
		info string
		// 是否付款
		isPayed bool
	}{
		{
			i: &Invoice{
				// 设置未过期时间
				ExpireDate: now.Add(24 * time.Hour),
				State:      FinishedPay,
			},
			info:    "无需付款",
			isPayed: true,
		},
		{
			i: &Invoice{
				// 测试当前时间
				ExpireDate: now,
				State:      FinishedPay,
			},
			info:    "无需付款",
			isPayed: true,
		},
		{
			i: &Invoice{
				// 设置已过期时间
				ExpireDate: now.Add(-24 * time.Hour),
				State:      FinishedPay,
			},
			info:    "无需付款，账单过期",
			isPayed: true,
		},
		{
			i: &Invoice{
				ExpireDate: now.Add(24 * time.Hour),
				State:      NeedPay,
			},
			info:    "需要付款",
			isPayed: false,
		},
		{
			i: &Invoice{
				ExpireDate: now.Add(-24 * time.Hour),
				State:      NeedPay,
			},
			info:    "需要付款，账单过期",
			isPayed: false,
		},
	}

	for _, v := range testData {
		info, state := v.i.GetStatus()
		if info != v.info {
			format := "wrong info:\n\twant :%s\n\thave: %s\n%v\n"
			t.Errorf(format, v.info, info, v)
		}
		if state != v.isPayed {
			format := "wrong payment state:\n\twant: %v\n\thave: %v\n%v\n"
			t.Errorf(format, v.isPayed, state, v)
		}
	}
}
