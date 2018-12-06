package parser

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"

	"schannel-qt5/urls"
)

var (
	getDataInfo = regexp.MustCompile(`.+ \((.+(?:GB|MB|KB))\)`)
	getTotal    = regexp.MustCompile(`.+ \(流量：(.+(?:GB|MB|KB))\)`)
)

// GetService 返回所有可用的套餐的信息
func GetService(data string) []*Service {
	res := make([]*Service, 0)

	table, _ := goquery.NewDocumentFromReader(strings.NewReader(data))
	// id为tableServicesList的table里有所有的服务信息
	table.Find("#tableServicesList tbody tr").Each(func(i int, s *goquery.Selection) {
		ser := new(Service)
		tds := s.Find("td")
		// 第一列是服务名称
		ser.Name = tds.Eq(0).Text()
		// 第二列是详细信息页面链接和价格
		link, _ := tds.Eq(1).Find("a").Attr("href")
		ser.Link = urls.RootPath + link
		ser.Price = tds.Eq(1).Text()
		// 第三列是服务到期时间
		expire := tds.Eq(2).Find("span").Text()
		ser.Expires, _ = time.ParseInLocation("2006-01-02", expire, time.Local)
		// 第四列是服务状态信息
		ser.State = tds.Eq(3).Text()
		res = append(res, ser)
	})

	return res
}

// GetSSRInfo 获取套餐的详细使用信息
func GetSSRInfo(data string, ser *Service) *SSRInfo {
	res := NewSSRInfo(ser)

	dom, _ := goquery.NewDocumentFromReader(strings.NewReader(data))
	// 信息包含在第2-4个section.panel里
	sections := dom.Find("section.panel")

	// 第2个section的table是端口和密码
	serverInfo := sections.Eq(1)
	portAndPasswd := serverInfo.Find("table").First()
	res.Port, _ = strconv.ParseInt(portAndPasswd.Find("tbody tr").Find("td").First().Text(), 10, 64)
	res.Passwd = portAndPasswd.Find("tbody tr").Find("td").Eq(1).Text()

	// 第3个section的header里是套餐总量
	usageInfo := sections.Eq(2)
	total := usageInfo.Find("header").Text()
	res.TotalData = getTotal.FindStringSubmatch(total)[1]

	usage := usageInfo.Find("#plugin-usage").Find("p")
	res.UsedData = getDataInfo.FindStringSubmatch(usage.First().Text())[1]
	res.Upload = getDataInfo.FindStringSubmatch(usage.Eq(1).Text())[1]
	res.Download = getDataInfo.FindStringSubmatch(usage.Eq(2).Text())[1]

	// 第4个section的table是节点信息表
	sections.Eq(3).Find("table").
		Find("tbody tr").Each(func(i int, s *goquery.Selection) {
		node := new(SSRNode)
		tds := s.Children()

		node.NodeName = tds.First().Text()
		node.Type = tds.Eq(1).Text()
		node.IP = tds.Eq(2).Text()
		node.Crypto = tds.Eq(3).Text()
		node.Proto = tds.Eq(4).Text()
		node.Minx = tds.Eq(5).Text()
		node.Port = res.Port
		node.Passwd = res.Passwd

		res.Nodes = append(res.Nodes, node)
	})

	return res
}

// GetInvoices 返回所有账单信息
func GetInvoices(data string) []*Invoice {
	invoiceList := make([]*Invoice, 0, 2)

	dom, _ := goquery.NewDocumentFromReader(strings.NewReader(data))
	invoiceTable := dom.Find("#tableInvoicesList")

	invoiceTable.Find("tbody tr").Each(func(_ int, s *goquery.Selection) {
		invoice := new(Invoice)
		tds := s.Find("td")

		invoice.Number = tds.First().Text()

		startDate := tds.Eq(1).Find("span").Text()
		expireDate := tds.Eq(2).Find("span").Text()
		invoice.StartDate, _ = time.ParseInLocation("2006-01-02", startDate, time.Local)
		invoice.ExpireDate, _ = time.ParseInLocation("2006-01-02", expireDate, time.Local)

		payment, exists := tds.Eq(3).Attr("data-order")
		// 如果取不到，就用默认值0
		if exists {
			invoice.Payment, _ = strconv.ParseInt(payment, 10, 64)
		}

		if tds.Eq(4).Text() == "已付款" {
			invoice.State = FinishedPay
		} else if tds.Eq(4).Text() == "未付款" {
			invoice.State = NeedPay
		}

		link, _ := tds.Eq(5).Find("a").Attr("href")
		invoice.Link = urls.RootPath + link

		invoiceList = append(invoiceList, invoice)
	})

	return invoiceList
}

// GetInvoiceDownloadURL 获取invoice下载地址
func GetInvoiceDownloadURL(data string) string {
	dom, _ := goquery.NewDocumentFromReader(strings.NewReader(data))
	downloadBtn := dom.Find("i.fa-download").Parent()
	downloadURL, exists := downloadBtn.Attr("href")
	if !exists {
		return ""
	}

	return urls.RootPath + downloadURL
}
