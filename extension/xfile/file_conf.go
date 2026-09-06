// Package xfile
// @author: fengyi
// @date: 2024/6/21
// @note:
package xfile

import (
	"general-agent/extension/errorx"
	"os"
)

type FileConf struct {
	filePath string
}

func NewFileConf(filePath string) *FileConf {
	return &FileConf{filePath: filePath}
}

func (jf *FileConf) Read() ([]byte, error) {
	_, err := os.Stat(jf.filePath)
	if err != nil {
		return nil, errorx.ErrDefault.WithError(err).WithMessage("file not exist")
	}
	return os.ReadFile(jf.filePath)
}
func (jf *FileConf) Write(content string) error {
	return os.WriteFile(jf.filePath, []byte(content), 0644)
}
