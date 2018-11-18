package widgets

import (
	"testing"

	"math"
)

func TestGetGeoName(t *testing.T) {
	testData := []*struct {
		src string
		res string
	}{
		{
			src: "Tokyo_0",
			res: "日本",
		},
		{
			src: "US_LA_1",
			res: "美国-洛杉矶",
		},
		{
			src: "US_1",
			res: "美国",
		},
		{
			src: "EUR_0",
			res: "欧洲",
		},
		{
			src: "EUR_AMSD_0",
			res: "欧洲-阿姆斯特丹",
		},
		{
			src: "test",
			res: "Unknown",
		},
		{
			src: "US_test_1",
			res: "美国-Unknown",
		},
	}

	for _, v := range testData {
		geoName := getGeoName(v.src)
		if geoName != v.res {
			format := "error when:\t%v\n\twant:\t%v\n\thave:\t%v\n"
			t.Errorf(format, v.src, v.res, geoName)
		}
	}
}

func TestFastOpenAble(t *testing.T) {
	testData := []*struct {
		version string
		res     bool
	}{
		{
			version: "2.6.32.1",
			res:     false,
		},
		{
			version: "3.6.2",
			res:     false,
		},
		{
			version: "4.17.19",
			res:     true,
		},
		{
			version: "3.7.0",
			res:     true,
		},
	}

	for _, v := range testData {
		if able := fastOpenAble(v.version); able != v.res {
			format := "fastOpenAble error: %v\n\twant: %v\n\thave: %v\n"
			t.Errorf(format, v.version, v.res, able)
		}
	}
}

func TestConvertToKb(t *testing.T) {
	testData := []*struct {
		data string
		res  int
	}{
		{
			data: "10KB",
			res:  10,
		},
		{
			data: "10MB",
			res:  1024 * 10,
		},
		{
			data: "10GB",
			res:  1024 * 1024 * 10,
		},
		{
			data: "10.41KB",
			// 10.41 * 1
			res: 10,
		},
		{
			data: "10.41MB",
			// 10.41 * 1024
			res: 10659,
		},
		{
			data: "10.41GB",
			// 10.41 * 1024 * 1024
			res: 10915676,
		},
		{
			data: "",
			res:  -1,
		},
		{
			data: "test",
			res:  -1,
		},
	}

	for _, v := range testData {
		res := convertToKb(v.data)
		if res != v.res {
			format := "convertToKb error: %v\n\twant: %v\n\thave: %v\n"
			t.Errorf(format, v.data, v.res, res)
		}
	}
}

func TestComputeSizeUnit(t *testing.T) {
	testData := []*struct {
		data  []int
		ratio int
		unit  string
	}{
		{
			data:  []int{1, 10 * GB, 5 * MB, 1000 * MB, 100 * KB},
			ratio: GB,
			unit:  "GB",
		},
		{
			data:  []int{4000, 50, 546, 1},
			ratio: MB,
			unit:  "MB",
		},
		{
			data:  []int{0, 100 * KB, 56 * KB, 1000 * KB},
			ratio: KB,
			unit:  "KB",
		},
		{
			data:  []int{0},
			ratio: KB,
			unit:  "KB",
		},
	}

	for _, v := range testData {
		ratio, unit := computeSizeUnit(v.data)
		if ratio != v.ratio || unit != v.unit {
			format := "computeSizeUnit error: %v\n\twant: (%v, %v)\n\thave: (%v, %v)\n"
			t.Errorf(format, v.data, v.ratio, v.unit, ratio, unit)
		}
	}
}

func TestComputeRange(t *testing.T) {
	testData := []*struct {
		data     []int
		ratio    int
		min, max float64
	}{
		{
			data:  []int{1, 10, 2},
			ratio: KB,
			min:   0.5,
			max:   10.5,
		},
		{
			data:  []int{1900 * KB, 100 * MB},
			ratio: MB,
			min:   1.4,
			max:   100.5,
		},
		{
			data:  []int{1782580 * KB, 2590 * MB, 2 * GB},
			ratio: GB,
			min:   1.2,
			max:   3.0,
		},
	}

	for _, v := range testData {
		min, max := computeRange(v.data, v.ratio)
		if !floatEqual(min, v.min) {
			format := "range min wrong:\n\twant: %v\n\thave: %v\n"
			t.Errorf(format, v.min, min)
		}
		if !floatEqual(max, v.max) {
			format := "range max wrong:\n\twant: %v\n\thave: %v\n"
			t.Errorf(format, v.max, max)
		}
	}
}

func floatEqual(a, b float64) bool {
	EPSILON := 0.00000001
	return math.Abs(a-b) < EPSILON
}
