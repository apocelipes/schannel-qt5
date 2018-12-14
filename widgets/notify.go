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

// ShowNotification 显示libnotify实现的气泡消息框
// duration == 0时使用默认delay
// 出错信息输出到stderr，不进入log
func ShowNotification(title, text, image string, delay time.Duration) {
	notifyDelay := duration2millisecond(delay)
	if notifyDelay == 0 {
		notifyDelay = duration2millisecond(defaultNotifyDelay)
	}

	libnotify.Init(applicationName)
	defer libnotify.UnInit()

	notify := libnotify.NotificationNew(title, text, image)
	if notify == nil {
		fmt.Fprintf(os.Stderr, "Unable to create a new notification\n")
		return
	}
	defer libnotify.NotificationClose(notify)
	libnotify.NotificationSetTimeout(notify, notifyDelay)

	libnotify.NotificationShow(notify)
	// 保证有足够的时间让notify显示
	time.Sleep(time.Duration(notifyDelay)*time.Millisecond + 1*time.Second)
}

// duration2millisecond 将time.Duration转换成millisecond
// duration不足1ms将返回0
func duration2millisecond(duration time.Duration) int32 {
	res := int32(duration / time.Millisecond)
	if res <= 0 {
		return 0
	}

	return res
}
