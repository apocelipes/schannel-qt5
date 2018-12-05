package widgets

import (
	"log"
	"net/http"
	"sync"
	"time"

	"schannel-qt5/crawler"
	"schannel-qt5/parser"
)

const (
	// 日志子前缀
	dataBridgePrefix = "data_bridge: "
)

// UserDataBridge 传递用户数据到界面组件
type UserDataBridge interface {
	sync.Locker
	// ServiceInfos 获取服务信息
	ServiceInfos() []*parser.Service
	// SSRInfos 根据service信息获取ssr使用和节点信息
	SSRInfos(ser *parser.Service) *parser.SSRInfo
	// Invoices 获取账单信息
	Invoices() []*parser.Invoice
	// GetLogger 获取logger
	GetLogger() *log.Logger
	// GetCookies 获取用户的cookie
	GetCookies() []*http.Cookie
}

// accountDataProxy 用于获取和缓存用户数据的代理类
type accountDataProxy struct {
	*sync.Mutex
	// 缓存过期时间，默认20min
	cached time.Time
	// 用户cookies，用于数据交互
	cookies []*http.Cookie
	// 记录日志
	logger *log.Logger
	// 代理URL
	proxy string
	// 用户数据
	ssrInfos []*parser.SSRInfo
	invoices []*parser.Invoice
}

// NewDataBridge 生成用户数据接口
func NewDataBridge(cookies []*http.Cookie, proxy string, logger *log.Logger) UserDataBridge {
	u := &accountDataProxy{}
	u.Mutex = &sync.Mutex{}
	u.cookies = cookies
	u.proxy = proxy
	u.logger = logger

	return u
}

// checkCacheExpired 检查缓存是否过期，如果过期就更新
// 虽然非并发安全，但是不公开，只能由公开接口调用，调用公开接口前加锁则不会发生数据竞争
func (a *accountDataProxy) checkCacheExpired() {
	if time.Now().Sub(a.cached) < 20*time.Minute {
		return
	}

	servicesHTML, err := crawler.GetServiceHTML(a.cookies, a.proxy)
	if err != nil {
		a.logger.Printf(dataBridgePrefix+"%v\n", err)
		return
	}

	servicesList := parser.GetService(servicesHTML)
	tmp := make([]*parser.SSRInfo, 0, len(servicesList))
	for _, ser := range servicesList {
		infoHTML, err := crawler.GetSSRInfoHTML(ser, a.cookies, a.proxy)
		if err != nil {
			a.logger.Printf(dataBridgePrefix+"%v\n", err)
			return
		}

		ssrInfo := parser.GetSSRInfo(infoHTML, ser)
		tmp = append(tmp, ssrInfo)
	}
	a.ssrInfos = tmp

	invoiceHTML, err := crawler.GetInvoiceHTML(a.cookies, a.proxy)
	if err != nil {
		a.logger.Printf(dataBridgePrefix+"%v\n", err)
		return
	}
	a.invoices = parser.GetInvoices(invoiceHTML)

	a.cached = time.Now()
}

// ServiceInfo 获取服务信息
// 并发安全，因为只会修改slice而不会修改其中item的具体数据
func (a *accountDataProxy) ServiceInfos() []*parser.Service {
	a.Lock()
	defer a.Unlock()
	a.checkCacheExpired()

	sers := make([]*parser.Service, 0, len(a.ssrInfos))
	for i := range a.ssrInfos {
		sers = append(sers, a.ssrInfos[i].Service)
	}

	return sers
}

// SSRInfos 根据给出的Service返回ssr服务和节点信息
// 并发安全，因为只有slice会被修改，其中的item不会被修改，因此没有数据竞争
func (a *accountDataProxy) SSRInfos(ser *parser.Service) *parser.SSRInfo {
	a.Lock()
	defer a.Unlock()
	a.checkCacheExpired()

	for _, v := range a.ssrInfos {
		if *v.Service == *ser {
			return v
		}
	}

	return nil
}

// Invoices 返回所有账单信息
// 并发安全
func (a *accountDataProxy) Invoices() []*parser.Invoice {
	a.Lock()
	defer a.Unlock()
	a.checkCacheExpired()

	invoices := make([]*parser.Invoice, len(a.invoices))
	copy(invoices, a.invoices)

	return invoices
}

// GetLogger 获取共享logger
func (a *accountDataProxy) GetLogger() *log.Logger {
	a.Lock()
	defer a.Unlock()

	return a.logger
}

// GetCookies 获取登录后的用户身份cookies
func (a *accountDataProxy) GetCookies() []*http.Cookie {
	a.Lock()
	defer a.Unlock()
	cookies := make([]*http.Cookie, len(a.cookies))
	copy(cookies, a.cookies)

	return cookies
}
