package widgets

import (
	"errors"
	"os"
	"path/filepath"
)

// SavePathRecorder 缓存用户上一次保存文件的目录，通过服务名称获取目录数据
type SavePathRecorder struct {
	recorder map[string]string
}

// 记录用户上次保存文件的目录
var pathRecorder = NewSavePathRecorder()

func NewSavePathRecorder() *SavePathRecorder {
	recorder := &SavePathRecorder{
		recorder: make(map[string]string),
	}

	return recorder
}

// LastSavePath 根据service返回上次保存文件的目录，service不存在则返回false
func (r *SavePathRecorder) LastSavePath(service string) (string, bool) {
	path, ok := r.recorder[service]
	return path, ok
}

// SetLastSavePath 设置service上一次保存文件的目录
// service不能为空
func (r *SavePathRecorder) SetLastSavePath(service, path string) error {
	if service == "" {
		return errors.New("empty service")
	}

	path = filepath.Dir(path)
	r.recorder[service] = path
	return nil
}

// defaultSavePath 返回文件默认存储位置
// 优先选择上次保存文件的目录(仅限本次会话期间)
// 否则选择$HOME
func defaultSavePath(service, fileName string) (string, error) {
	savePath, ok := pathRecorder.LastSavePath(service)
	if !ok {
		home := os.Getenv("HOME")
		if home == "" {
			return "", errors.New("无法获取HOME")
		}

		savePath = filepath.Join(home, fileName)
	} else {
		savePath = filepath.Join(savePath, fileName)
	}

	return savePath, nil
}
