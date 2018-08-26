package parser

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"schannel-qt5/urls"
)

func TestGetTotal(t *testing.T) {
	data := `使用报表 (流量：50GB)`
	res := getTotal.FindStringSubmatch(data)
	if res[1] != "50GB" {
		t.Error("regexp getTotal has some problem.")
	}
}

func TestGetDataInfo(t *testing.T) {
	data1 := `已使用 (16.14GB)`
	data2 := `上传 (14.66MB)`
	data3 := `下载 (16.12GB)`

	if getDataInfo.FindStringSubmatch(data1)[1] != "16.14GB" {
		t.Error("regexp getDataInfo has some problem on getting used.")
	}

	if getDataInfo.FindStringSubmatch(data2)[1] != "14.66MB" {
		t.Error("regexp getDataInfo has some problem on getting upload.")
	}

	if getDataInfo.FindStringSubmatch(data3)[1] != "16.12GB" {
		t.Error("regexp getDataInfo has some problem on getting download.")
	}
}

func TestGetInvoice(t *testing.T) {
	// 应该被解析出来的信息
	correctRes := []Invoice{
		{
			Number:  "12345",
			Link:    urls.RootPath + "test1",
			Payment: 10,
			State:   NeedPay,
		}, {
			Number:  "2345",
			Link:    urls.RootPath + "test2",
			Payment: 10,
			State:   FinishedPay,
		}, {
			Number:  "345",
			Link:    urls.RootPath + "test3",
			Payment: 10,
			State:   FinishedPay,
		}, {
			Number:  "4390",
			Link:    urls.RootPath + "test4",
			Payment: 10,
			State:   FinishedPay,
		},
	}

	correctRes[0].StartDate, _ = time.ParseInLocation("2006-01-02", "2018-04-10", time.Local)
	correctRes[0].ExpireDate, _ = time.ParseInLocation("2006-01-02", "2018-04-11", time.Local)

	correctRes[1].StartDate, _ = time.ParseInLocation("2006-01-02", "2018-04-30", time.Local)
	correctRes[1].ExpireDate, _ = time.ParseInLocation("2006-01-02", "2018-04-30", time.Local)

	correctRes[2].StartDate, _ = time.ParseInLocation("2006-01-02", "2018-05-28", time.Local)
	correctRes[2].ExpireDate, _ = time.ParseInLocation("2006-01-02", "2018-05-29", time.Local)

	correctRes[3].StartDate, _ = time.ParseInLocation("2006-01-02", "2018-06-28", time.Local)
	correctRes[3].ExpireDate, _ = time.ParseInLocation("2006-01-02", "2018-06-29", time.Local)

	f, err := os.Open("testdata/invoice.html")
	if err != nil {
		t.Error(err)
	}
	defer f.Close()
	testData, err := ioutil.ReadAll(f)
	if err != nil {
		t.Error(err)
	}

	res := GetInvoices(string(testData))
	if len(res) != len(correctRes) {
		t.Errorf("解析到的数据量不正确，期望%d个，实际%d个\n", len(correctRes), len(res))
	}

	for i := range res {
		if res[i].Number != correctRes[i].Number {
			t.Errorf("订单号错误，期望%v，实际%v\n", correctRes[i].Number, res[i].Number)
		}
		if res[i].Link != correctRes[i].Link {
			t.Errorf("链接错误，期望%v，实际%v\n", correctRes[i].Link, res[i].Link)
		}
		if res[i].ExpireDate != correctRes[i].ExpireDate {
			t.Errorf("过期日期错误，期望%v，实际%v\n", correctRes[i].ExpireDate, res[i].ExpireDate)
		}
		if res[i].StartDate != correctRes[i].StartDate {
			t.Errorf("开始日期错误， 期望%v，实际%v\n", correctRes[i].StartDate, res[i].StartDate)
		}
		if res[i].Payment != correctRes[i].Payment {
			t.Errorf("支付金额错误，期望%v，实际%v\n", correctRes[i].Payment, res[i].Payment)
		}
		if res[i].State != correctRes[i].State {
			t.Errorf("状态错误，期望%v，实际%v\n", correctRes[i].State, res[i].State)
		}
	}
}
