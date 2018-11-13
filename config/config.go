package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	// 默认配置文件路径
	configPath = ".local/share/schannel-qt5.json"
)

var (
	// ErrHOME 无法查找$HOME
	ErrHOME = errors.New("can't find $HOME in your environments")
	// ErrNotAbs 路径无法解析为绝对路径
	ErrNotAbs = errors.New("path is not an abs path")
)

// UserConfig 用户配置
type UserConfig struct {
	// client config
	Proxy   JSONProxy `json:"proxy_url"`
	LogFile JSONPath  `json:"log_file"`

	// ssr config
	SSRNodeConfigPath   JSONPath `json:"ssr_node_config_path"`
	SSRClientConfigPath JSONPath `json:"ssr_client_config_path"`

	// ssr client bin path
	SSRBin JSONPath `json:"ssr_bin"`

	// ssr client config的实体数据
	SSRClientConfig ClientConfig `json:"-"`
}

// ConfigPath 返回`～`被替换为$HOME的config path
func ConfigPath() (string, error) {
	home, exist := os.LookupEnv("HOME")
	if !exist {
		return "", ErrHOME
	}

	return filepath.Join(home, configPath), nil
}

// StoreConfig 将配置存储进ConfigPath路径的文件
func (u *UserConfig) StoreConfig() error {
	storePath, err := ConfigPath()
	if err != nil {
		return err
	}

	f, err := os.OpenFile(storePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0664)
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := json.MarshalIndent(u, "", "\t")
	if err != nil {
		return err
	}

	if _, err = f.Write(data); err != nil {
		return err
	}

	clientConfigPath, err := u.SSRClientConfigPath.AbsPath()
	if err != nil {
		return err
	}
	if err := u.SSRClientConfig.Store(clientConfigPath); err != nil {
		return err
	}

	return nil
}

// LoadConfig 从ConfigPath给出的配置文件路径读出配置
func (u *UserConfig) LoadConfig() error {
	loadPath, err := ConfigPath()
	if err != nil {
		return err
	}

	f, err := os.Open(loadPath)
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(data, u); err != nil {
		return err
	}

	clientConfigPath, err := u.SSRClientConfigPath.AbsPath()
	if err != nil {
		return err
	}
	if err := u.SSRClientConfig.Load(clientConfigPath); err != nil {
		return err
	}

	return nil
}
