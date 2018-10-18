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
			src: "US_test",
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
