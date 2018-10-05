package widgets

import (
	"errors"
	"regexp"
	"schannel-qt5/config"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const (
	// KB 流量单位，1kb
	KB = 1
	// MB 流量单位，1024kb
	MB = 1024 * KB
	// GB 1024mb
	GB = MB * 1024
)

var (
	// Matcher 匹配数字和流量单位
	usedMatcher = regexp.MustCompile(`(.+)(GB|MB)`)
	// 匹配linux kernel版本 A.B.C
	versionMatcher = regexp.MustCompile(`(\d+\.\d+\.\d+)`)
)

// convertToKb 将string类型的流量数据转换成以kb为单位的int
func convertToKb(data string) int {
	tmp := usedMatcher.FindStringSubmatch(data)
	number, err := strconv.ParseFloat(tmp[1], 64)
	if err != nil {
		return -1
	}

	var ratio int
	switch tmp[2] {
	case "GB":
		ratio = GB
	case "MB":
		ratio = MB
	case "KB":
		ratio = KB
	default:
		return -1
	}

	res := number * float64(ratio)
	return int(res)
}

// computeRatio 根据ratio计算阀值
func computeRatio(total int) int {
	res := float64(total) * HighRatio
	return int(res)
}

// time2string 将日期转换为2006-01-02的格式
func time2string(t time.Time) string {
	return t.Format("2006-01-02")
}

// kernelVersion 获取linux kernel version
func kernelVersion() (string, error) {
	uname := syscall.Utsname{}
	if err := syscall.Uname(&uname); err != nil {
		return "", err
	}

	ver := arr2str(uname.Release)
	version := versionMatcher.FindStringSubmatch(ver)[1]

	return version, nil
}

// arr2str 将[65]int8转换成string
func arr2str(data [65]int8) string {
	var buf [65]byte
	for i, b := range data {
		buf[i] = byte(b)
	}

	str := string(buf[:])
	// 截断\0
	if i := strings.Index(str, "\x00"); i != -1 {
		str = str[:i]
	}
	return str
}

// fastOpenAble  检查是否支持特tcp-fast-open
// linux kernel version > 3.7
func fastOpenAble(version string) bool {
	// index 0是主要版本，1是发布版本，2是修复版本
	// 分别与a，b，c对应
	ver := strings.Split(version, ".")
	a, _ := strconv.Atoi(ver[0])
	b, _ := strconv.Atoi(ver[1])

	if a >= 3 {
		if a > 3 || (a == 3 && b >= 7) {
			return true
		}
	}

	return false
}

// showErrorMsg 控制error label的显示
// err为nil则代表没有错误发生，如果label可见则设为隐藏
// err不为nil时设置label可见
// 设置label可见时返回true，否则返回false（不受label原有状态影响）
func showErrorMsg(label *ColorLabel, err error) bool {
	if err != nil {
		label.Show()
		return true
	}

	label.Hide()
	return false
}

// 检查路径是否是绝对路径，且不是目录
func checkPath(path string) error {
	jpath := config.JSONPath{Data: path}
	if _, err := jpath.AbsPath(); err != nil {
		return err
	} else if path[len(path)-1] == '/' {
		return errors.New("dir is not allowed")
	}

	return nil
}
