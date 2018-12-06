package crawler

import (
	"compress/gzip"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"golang.org/x/net/publicsuffix"
)

const (
	// UA is Chrome's User-Agent
	UA = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.79 Safari/537.36"
	// AcceptType is Chrome's Accept-Type
	AcceptType = "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8"
)

// genClientWithProxy 生成httpclient，并设置代理为proxy指定的代理服务器
// proxy url支持http，https和socks5协议
func genClientWithProxy(proxy string) (*http.Client, error) {
	client := new(http.Client)
	// all cookiejar's users should use "golang.org/x/net/publicsuffix"
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return nil, err
	}
	client.Jar = jar

	if proxy != "" {
		proxyURL, err := url.Parse(proxy)
		if err != nil {
			return nil, err
		}
		// 设置proxy
		client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyURL)}
	}

	return client, nil
}

// setRequestHeader 设置请求头信息
// cookies为nil时将被忽略，referer为“”时同样被忽略
func setRequestHeader(request *http.Request, cookies []*http.Cookie, referer string) {
	request.Header.Set("accept", AcceptType)
	request.Header.Set("accept-encoding", "gzip")
	if referer != "" {
		request.Header.Set("referer", referer)
	}
	request.Header.Set("user-agent", UA)

	for _, c := range cookies {
		request.AddCookie(c)
	}
}

// getPage 获取url指定的各种账户管理页面信息, cookies用于身份认证
func getPage(url, referer string, cookies []*http.Cookie, proxy string) (string, error) {
	client, err := genClientWithProxy(proxy)
	if err != nil {
		return "", err
	}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	setRequestHeader(request, cookies, referer)

	resp, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	htmlReader, err := gzip.NewReader(resp.Body)
	if err != nil {
		return "", err
	}
	defer htmlReader.Close()

	data, err := ioutil.ReadAll(htmlReader)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
