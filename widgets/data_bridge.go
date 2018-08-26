package widgets

import (
  "net/http"
  "sync"
  "time"
  "log"

  "schannel-qt5/parser"
  "schannel-qt5/crawler"
)

// 日志子前缀
const (
  prefix = "data_bridge: "
)

// UserDataBridge 传递用户数据到界面组件
type UserDataBridge interface {
  sync.Locker
  // ServiceInfos 获取服务信息
  ServiceInfos() []*parser.Service
  // SSRInfos 获取ssr使用和节点信息
  SSRInfos() []*parser.SSRInfo
  // Invoices 获取账单信息
  Invoices() []*parser.Invoice
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
// 虽然非并发安全，但是不出口，只能由公开接口调用，调用公开接口前加锁则不会发生数据竞争
func (a *accountDataProxy) checkCacheExpired() {
  if time.Now().Sub(a.cached) < 20 * time.Minute {
    return
  }

  servicesHTML, err := crawler.GetServiceHTML(a.cookies, a.proxy)
  if err != nil {
    a.logger.Printf(prefix + "%v\n", err)
    return
  }

  servicesList := parser.GetService(servicesHTML)
  tmp := make([]*parser.SSRInfo, 0)
  for _, ser := range servicesList {
    infoHTML, err := crawler.GetSSRInfoHTML(ser, a.cookies, a.proxy)
    if err != nil {
      a.logger.Printf(prefix + "%v\n", err)
      return
    }

    ssrInfo := parser.GetSSRInfo(infoHTML, ser)
    tmp = append(tmp, ssrInfo)
  }
  a.ssrInfos = tmp

  invoiceHTML, err := crawler.GetInvoiceHTML(a.cookies, a.proxy)
  if err != nil {
    a.logger.Printf(prefix + "%v\n", err)
    return
  }
  a.invoices = parser.GetInvoices(invoiceHTML)

  a.cached = time.Now()
}

// ServiceInfo 获取服务信息
// 非并发安全，调用前需要先加锁
func (a *accountDataProxy) ServiceInfos() []*parser.Service {
  a.checkCacheExpired()

  sers := make([]*parser.Service, 0, len(a.ssrInfos))
  for i := range a.ssrInfos {
    sers = append(sers, a.ssrInfos[i].Service)
  }

  return sers
}

// SSRInfos 返回ssr服务和节点信息
// 非并发安全，调用前需要先加锁
func (a *accountDataProxy) SSRInfos() []*parser.SSRInfo {
  a.checkCacheExpired()

  return a.ssrInfos
}

// Invoices 返回所有账单信息
// 非并发安全，调用前需要先加锁
func (a *accountDataProxy) Invoices() []*parser.Invoice {
  a.checkCacheExpired()

  return a.invoices
}
