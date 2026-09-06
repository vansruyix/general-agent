// Package xfile
// @author: fengyi
// @date: 2024/7/10
// @note:
package xfile

import "mime/multipart"

func ReadFileHeader(keyFile *multipart.FileHeader) ([]byte, error) {
	bytes := make([]byte, keyFile.Size)
	file, err := keyFile.Open()
	if err != nil {
		return nil, err
	}
	if _, err := file.Read(bytes); err != nil {
		return nil, err
	}
	return bytes, nil
}
