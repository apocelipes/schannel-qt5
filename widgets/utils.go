package widgets

import (
	"regexp"
	"strconv"
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
	Matcher = regexp.MustCompile(`(.+)(GB|MB)`)
)

// convertToKb 将string类型的流量数据转换成以kb为单位的int
func convertToKb(data string) int {
	tmp := Matcher.FindStringSubmatch(data)
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
