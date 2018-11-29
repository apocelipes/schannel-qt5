package widgets

import (
	"errors"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"

	"schannel-qt5/config"
)

const (
	// KB 流量单位，1kb
	KB = 1
	// MB 流量单位，1024kb
	MB = 1024 * KB
	// GB 1024mb
	GB = MB * 1024
	// HighRatio 流量使用量警告阀值
	HighRatio = 0.9
)

var (
	// Matcher 匹配数字和流量单位
	usedMatcher = regexp.MustCompile(`(.+)(GB|MB|KB)`)
	// 匹配linux kernel版本 A.B.C
	versionMatcher = regexp.MustCompile(`(\d+\.\d+\.\d+)`)
	// 名字对应的国家/地区/城市缩写
	geoAreaName = map[string]string{
		"US":    "美国",
		"SG":    "新加坡",
		"Tokyo": "日本",
		"EUR":   "欧洲",
		"AMSD":  "阿姆斯特丹",
		"LA":    "洛杉矶",
		"ALT":   "亚特兰大",
		"FRK":   "美因河畔法兰克福",
	}
)

// convertToKb 将string类型的流量数据转换成以kb为单位的int
// data格式为 number[KB|MB|GB]
func convertToKb(data string) int {
	tmp := usedMatcher.FindStringSubmatch(data)
	if tmp == nil {
		return -1
	}
	// tmp[0]为完整字符串，1为数字，2为容量单位
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

// checkEmptyPath 检查路径是否为绝对路径且不为目录，允许空路径
func checkEmptyPath(path string) error {
	if path == "" {
		return nil
	}

	return checkPath(path)
}

// getGeoName 根据节点名字获取对应地区信息
func getGeoName(name string) string {
	areas := strings.Split(name, "_")
	// 名字不符合area_number
	if len(areas) <= 1 {
		return "Unknown"
	}

	areaName := make([]string, 0, len(areas)-1)
	// 最后一个元素为节点编号，忽略
	for _, v := range areas[:len(areas)-1] {
		if n, ok := geoAreaName[v]; ok {
			areaName = append(areaName, n)
		} else {
			areaName = append(areaName, "Unknown")
		}
	}

	return strings.Join(areaName, "-")
}

// computeSizeUnit 计算图表适用的size单位，为KB，MB或GB
// 因为月底网站会清空上月使用数据，所以选择最大的作为计量单位选择的依据
// 返回单位名称和相对KB的换算倍率
func computeSizeUnit(dataSet []int) (int, string) {
	sort.Ints(dataSet)
	max := dataSet[len(dataSet)-1]
	// 判断单位
	if max/GB != 0 {
		return GB, "GB"
	} else if max/MB != 0 {
		return MB, "MB"
	}

	return KB, "KB"
}

// computeRange 计算坐标轴的range，四舍五入为1位小数后+/-0.5，使折线平滑
func computeRange(dataSet []int, ratio int, unit string) (float64, float64) {
	var tuning float64
	switch unit {
	case "GB":
		tuning = 0.1
	case "MB":
		tuning = 0.5
	case "KB":
		tuning = 5
	}

	data := make([]float64, 0, len(dataSet))
	for _, v := range dataSet {
		value := float64(v) / float64(ratio)
		data = append(data, value)
	}
	sort.Float64s(data)

	max := math.Trunc(data[len(data)-1] * 10 + 0.5) / 10
	min := math.Trunc(data[0] * 10 + 0.5) / 10

	return math.Max(min - tuning, 0), max + tuning
}
