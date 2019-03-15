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

// GetCountryCity 返回ip的国家/地区和城市名称
// lang为语言名称，e.g. "zh-CN", "en"
func GetCountryCity(ip string, lang string) (string, string, error) {
	dbDir, err := GetGeoIPSavePath()
	if err != nil {
		return "", "", err
	}

	dbPath := filepath.Join(dbDir, DatabaseName)
	db, err := geoip2.Open(dbPath)
	if err != nil {
		return "", "", err
	}
	defer db.Close()

	ipAddr := net.ParseIP(ip)
	record, err := db.City(ipAddr)
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
