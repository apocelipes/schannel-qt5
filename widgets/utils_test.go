package widgets

import (
	"testing"
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
