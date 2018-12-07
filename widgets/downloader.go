package widgets

import (
	"errors"
	"net/http"
	"os"
	"sync"

	"github.com/therecipe/qt/core"

	"schannel-qt5/crawler"
)

// HTTPDownloader 通过HTTP下载文件
// 配合goroutine和QProgressBar/QProgressDialog使用
type HTTPDownloader struct {
	core.QObject

	// updateProgress 更新下载进度，size为已下载的byte数
	// updateMax 通知主线程下载文件的大小
	// failed 下载失败
	// done 下载完成
	_ func(size int)  `signal:"updateProgress"`
	_ func(max int)   `signal:"updateMax"`
	_ func(err error) `signal:"failed"`
	_ func()          `signal:"done"`

	// stop停止下载并使Download返回
	_ func() `slot:"stop,auto"`

	// request 缓存下载请求
	// client 发起下载请求
	request *http.Request
	client  *http.Client

	// 获取请求结果
	// resp保存请求结果
	responses chan *http.Response
	resp      *http.Response

	// 控制isStopped标志，代表下载是否取消
	lock      *sync.Mutex
	isStopped bool
}

// NewHTTPDownloader2 创建下载器
// url为下载地址
// file为保存的本地文件地址
// referer为HTTP Header的Referer，可为空
// proxy为代理，可设置为空
// cookies用户身份凭证，可为空
func NewHTTPDownloader2(url, referer, proxy string,
	cookies []*http.Cookie) (*HTTPDownloader, error) {
	downloader := NewHTTPDownloader(nil)
	downloader.lock = &sync.Mutex{}
	var err error
	downloader.client, err = crawler.GenClientWithProxy(proxy)
	if err != nil {
		return nil, err
	}

	downloader.request, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	crawler.SetRequestHeader(downloader.request, cookies, referer, "")

	downloader.responses = make(chan *http.Response, 1)
	go func() {
		response, err := downloader.client.Do(downloader.request)
		if err != nil {
			downloader.Failed(err)
			return
		}

		downloader.responses <- response
		close(downloader.responses)
	}()

	return downloader, nil
}

// TotalSize 下载文件的总大小
func (d *HTTPDownloader) TotalSize() (int, error) {
	if d.resp == nil {
		var ok bool
		d.resp, ok = <-d.responses
		if !ok {
			return 0, errors.New("responses has been closed")
		}
	}

	return int(d.resp.ContentLength), nil
}

const (
	chunk = 1024 * 32 // 一次下载的数据块大小(byte)
)

// Download 下载文件，每下载一个chunk长度出发一次UpdateProgress信号
// 下载被取消时删除已下载的部分文件
func (d *HTTPDownloader) Download(file string) {
	totalSize, err := d.TotalSize()
	if err != nil {
		d.Failed(err)
	}
	defer d.resp.Body.Close()
	d.UpdateMax(totalSize)

	f, err := os.OpenFile(file, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		d.Failed(err)
	}
	defer f.Close()

	wrote := 0
	buf := make([]byte, chunk)
	for wrote < totalSize {
		d.lock.Lock()
		if d.isStopped {
			d.lock.Unlock()
			os.Remove(file)
			return
		}

		n, err := d.resp.Body.Read(buf)
		if err != nil {
			d.Failed(err)
		}

		// golang的writer屏蔽了部分写，因此只需要检查err
		_, err = f.Write(buf[:n])
		if err != nil {
			d.Failed(err)
		}

		wrote += n
		d.UpdateProgress(wrote)
		d.lock.Unlock()
	}

	d.Done()
}

// stop 设置isStopped为true停止下载
// fixme: 不会立刻停止，可能会延迟至下一次读写
func (d *HTTPDownloader) stop() {
	d.lock.Lock()
	d.isStopped = true
	d.lock.Unlock()
}
