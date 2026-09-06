// Package xfile
// @author: fengyi
// @date: 2024/4/26
// @note:
package xfile

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestGetProjectRoot(t *testing.T) {
	got := GetProjectRoot()
	fmt.Println(got)
}

func TestExecutable(t *testing.T) {
	// 获取当前可执行文件的路径
	exePath, err := os.Executable()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	exeDir := filepath.Dir(exePath)
	fmt.Println("Executable path:", exePath)
	fmt.Println("Executable directory:", exeDir)
}
