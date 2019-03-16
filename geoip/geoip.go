package geoip

import (
	"errors"
	"net"
	"os"
	"path/filepath"

	"github.com/oschwald/geoip2-golang"
)

const (
	// GeoIP Database的存放路径
	geoIPSavePath = ".local/share/data/schannel-qt5/GeoIP"
	// GeoIP Database下载地址
	DownloadPath = "https://geolite.maxmind.com/download/geoip/database/GeoLite2-City.mmdb.gz"
	// 数据库文件名
	DatabaseName = "GeoLite2-City.mmdb"
)

// GetGeoIPSavePath 返回完整的数据库路径，默认在$HOME/.local/share/data/schannel-qt5/GeoIP下
func GetGeoIPSavePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", errors.New("cannot find user home dir")
	}

	return filepath.Join(home, geoIPSavePath), nil
}

// getRecord 根据ip返回geoIP查询结果
func getRecord(ip string) (*geoip2.City, error) {
	dbDir, err := GetGeoIPSavePath()
	if err != nil {
		return nil, err
	}

	dbPath := filepath.Join(dbDir, DatabaseName)
	db, err := geoip2.Open(dbPath)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	ipAddr := net.ParseIP(ip)
	record, err := db.City(ipAddr)
	if err != nil {
		return nil, err
	}

	return record, nil
}

// GetCountryCity 返回ip的国家/地区和城市名称
// lang为语言名称，e.g. "zh-CN", "en"
func GetCountryCity(ip string, lang string) (string, string, error) {
	record, err := getRecord(ip)
	if err != nil {
		return "", "", err
	}

	var country, city string
	country = record.Country.Names[lang]
	city, exists := record.City.Names[lang]
	if !exists {
		for _, sub := range record.Subdivisions {
			city, exists = sub.Names[lang]
			if exists {
				break
			}
		}
	}

	return country, city, nil
}

// GetCountryISOCode 获取国家代码
func GetCountryISOCode(ip string) (string, error) {
	record, err := getRecord(ip)
	if err != nil {
		return "", nil
	}

	return record.Country.IsoCode, nil
}
