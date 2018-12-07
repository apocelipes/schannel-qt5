package crawler

import (
	"compress/gzip"
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"schannel-qt5/parser"
	"schannel-qt5/urls"
)

// GetCSRFToken 返回本次回话所使用的CSRFToken
func GetCSRFToken(proxy string) (string, []*http.Cookie, error) {
	client, err := GenClientWithProxy(proxy)
	if err != nil {
		return "", nil, err
	}

	request, err := http.NewRequest("GET", urls.AccountPath, nil)
	if err != nil {
		return "", nil, err
	}
	SetRequestHeader(request, nil, urls.RootPath, "gzip")

	resp, err := client.Do(request)
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()

	htmlReader, err := gzip.NewReader(resp.Body)
	if err != nil {
		return "", nil, err
	}
	defer htmlReader.Close()

	dom, err := goquery.NewDocumentFromReader(htmlReader)
	if err != nil {
		return "", nil, err
	}

	CSRFToken, exists := dom.Find("input[type='hidden'][name='token']").Eq(0).Attr("value")
	if !exists {
		return "", nil, errors.New("CSRFToken doesn't exist")
	}

	u2, _ := url.Parse(urls.RootPath)
	return CSRFToken, client.Jar.Cookies(u2), nil
}

// GetAuth 登录schannel并返回登陆成功后获得的cookies
// 这些cookies在后续的页面访问中需要使用
func GetAuth(user, passwd, proxy string) ([]*http.Cookie, error) {
	client, err := GenClientWithProxy(proxy)
	if err != nil {
		return nil, err
	}

	CSRFToken, session, err := GetCSRFToken(proxy)
	if err != nil {
		return nil, err
	}

	form := url.Values{}
	form.Set("token", CSRFToken)
	form.Set("username", user)
	form.Set("password", passwd)
	getLogin, err := http.NewRequest("POST", urls.LoginPath, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}

	SetRequestHeader(getLogin, session, urls.AccountPath, "gzip")
	getLogin.Header.Set("content-type", "application/x-www-form-urlencoded")

	resp, err := client.Do(getLogin)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 验证登录是否成功，如果incorrect的值为“true”，则登录失败
	incorrect := resp.Request.FormValue("incorrect")
	if incorrect == "true" {
		return nil, errors.New("登录验证失败")
	}

	u2, _ := url.Parse(urls.RootPath)
	cookies := make([]*http.Cookie, 0, 2)
	cookies = append(cookies, client.Jar.Cookies(u2)...)
	// 添加cfduid会话cookie
	for _, c := range session {
		if c.Name == "__cfduid" {
			cookies = append(cookies, c)
		}
	}

	return cookies, nil
}

// GetServiceHTML 获取所有已购买服务的状态信息，包括详细页面的地址
func GetServiceHTML(cookies []*http.Cookie, proxy string) (string, error) {
	return getPage(urls.ServiceListPath, urls.AccountPath, cookies, proxy)
}

// GetInvoiceHTML 获取账单页面的HTML,包含未付款和已付款账单
// 未付款账单显示在最前列
// 现只支持获取第一页
func GetInvoiceHTML(cookies []*http.Cookie, proxy string) (string, error) {
	return getPage(urls.InvoicePath, urls.AccountPath, cookies, proxy)
}

// GetSSRInfoHTML 获取服务详细信息页面的HTML，包含使用情况和节点信息
func GetSSRInfoHTML(service *parser.Service, cookies []*http.Cookie, proxy string) (string, error) {
	return getPage(service.Link, urls.ServiceListPath, cookies, proxy)
}

// GetInvoiceInfoHTML 获取账单详情页面内容
func GetInvoiceInfoHTML(invoice *parser.Invoice, cookies []*http.Cookie, proxy string) (string, error) {
	return getPage(invoice.Link, urls.InvoicePath, cookies, proxy)
}
