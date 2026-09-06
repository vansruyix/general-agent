// Package xjson
// @author: fengyi
// @date: 2024/6/21
// @note:
package xjson

type JsonFile struct {
	filePath string
}

func NewJsonFile(filePath string) *JsonFile {
	return &JsonFile{filePath: filePath}
}

func (jf *JsonFile) Unmarshal(v any) error {
	return FileUnmarshal(jf.filePath, v)
}
func (jf *JsonFile) Marshal(v any) error {
	return FileMarshal(jf.filePath, v)
}
