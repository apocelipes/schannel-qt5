package config

import "testing"

func TestIsURL(t *testing.T) {
	testData := []struct {
		data string
		res  bool
	}{
		{
			data: "http://example.com",
			res:  true,
		},
		{
			data: "https://example.com:8080",
			res:  true,
		},
		{
			data: "socks5://127.0.0.1:8000/",
			res:  true,
		},
		{
			data: "ftp://test.org/",
			res:  false,
		},
		{
			data: "127.0.0.1:1025/",
			res:  false,
		},
	}

	for _, v := range testData {
		proxy := JSONProxy{v.data}
		if proxy.IsURL() != v.res {
			t.Errorf("failed with %s\n", v.data)
		}
	}
}
