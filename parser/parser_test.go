package parser

import (
	"testing"

	"io/ioutil"
	"os"
	"time"

	"schannel-qt5/urls"
)

func TestGetTotal(t *testing.T) {
	testData := []*struct {
		data string
		// parse得到的值
		res string
	}{
		{
			data: `使用报表 (流量：50GB)`,
			res:  "50GB",
		},
		{
			data: `使用报表 (流量：25.25GB)`,
			res:  "25.25GB",
		},
	}

	for _, v := range testData {
		res := getTotal.FindStringSubmatch(v.data)[1]
		if res != v.res {
			format := "regexp getTotal failed: %v\n\twant: %v\n\thave: %v\n"
			t.Errorf(format, v, v.res, res)
		}
	}
}

func TestGetDataInfo(t *testing.T) {
	testData := []*struct {
		data string
		// parse得到的值
		res string
	}{
		{
			data: `已使用 (16.14GB)`,
			res:  "16.14GB",
		},
		{
			data: `上传 (14.66MB)`,
			res:  "14.66MB",
		},
		{
			data: `下载 (16.12GB)`,
			res:  "16.12GB",
		},
		{
			data: `下载 (5.74KB)`,
			res:  "5.74KB",
		},
	}

	for _, v := range testData {
		res := getDataInfo.FindStringSubmatch(v.data)[1]
		if res != v.res {
			format := "regexp getDataInfo failed: %s\n\twant: %s\n\thave: %s\n"
			t.Errorf(format, v.data, v.res, res)
		}
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
		format := "解析到的数据量不正确，期望%d个，实际%d个\n"
		t.Errorf(format, len(correctRes), len(res))
	}

	for i := range res {
		if res[i].Number != correctRes[i].Number {
			format := "订单号错误，期望%v，实际%v\n"
			t.Errorf(format, correctRes[i].Number, res[i].Number)
		}
		if res[i].Link != correctRes[i].Link {
			format := "链接错误，期望%v，实际%v\n"
			t.Errorf(format, correctRes[i].Link, res[i].Link)
		}
		if res[i].ExpireDate != correctRes[i].ExpireDate {
			format := "过期日期错误，期望%v，实际%v\n"
			t.Errorf(format, correctRes[i].ExpireDate, res[i].ExpireDate)
		}
		if res[i].StartDate != correctRes[i].StartDate {
			format := "开始日期错误， 期望%v，实际%v\n"
			t.Errorf(format, correctRes[i].StartDate, res[i].StartDate)
		}
		if res[i].Payment != correctRes[i].Payment {
			format := "支付金额错误，期望%v，实际%v\n"
			t.Errorf(format, correctRes[i].Payment, res[i].Payment)
		}
		if res[i].State != correctRes[i].State {
			format := "状态错误，期望%v，实际%v\n"
			t.Errorf(format, correctRes[i].State, res[i].State)
		}
	}
}
