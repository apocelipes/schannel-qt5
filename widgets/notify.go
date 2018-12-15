package widgets

import (
	"fmt"
	"os"
	"time"

	libnotify "github.com/mqu/go-notify"
)

const (
	applicationName = "schannel-qt5"
	// 默认气泡框显示时间
	defaultNotifyDelay = 3 * time.Second
)

// ShowNotification 显示org.freedesktop.Notifications气泡消息框
// duration == -1时使用默认delay
// duration == 0表示不设置超时，desktop notification将会一直显示
// 出错信息输出到stderr，不进入log
func ShowNotification(title, text, image string, delay time.Duration) {
	var notifyDelay int32
	if delay == -1 {
		notifyDelay = duration2millisecond(defaultNotifyDelay)
	} else {
		notifyDelay = duration2millisecond(delay)
		// 不合法值(包括duration不足1ms)，使用默认值进行替换
		if notifyDelay == -1 {
			notifyDelay = duration2millisecond(defaultNotifyDelay)
		}
	}

	libnotify.Init(applicationName)

	notify := libnotify.NotificationNew(title, text, image)
	if notify == nil {
		fmt.Fprintf(os.Stderr, "Unable to create a new notification\n")
		return
	}
	notify.SetTimeout(notifyDelay)

	// TODO: also show icon
	notify.Show()
}

// duration2millisecond 将time.Duration转换成millisecond
// duration不足1ms将返回-1
func duration2millisecond(duration time.Duration) int32 {
	res := int32(duration / time.Millisecond)
	if res < 0 {
		return -1
	}

	return res
}
