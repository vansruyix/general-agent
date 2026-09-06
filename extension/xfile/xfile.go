package xfile

import (
	"fmt"
	"general-agent/extension/errorx"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

// InternalFileNameRegexp 内部文件名合法正则
var InternalFileNameRegexp = regexp.MustCompile(`^[a-zA-z0-9._-]+$`)
var AppPath = ""

var UploadPath = filepath.Join(os.TempDir(), "upload")

func init() {
	if runtime.GOOS == "windows" {
		AppPath = GetProjectRoot()
	} else {
		exePath, err := os.Executable()
		if err != nil {
			fmt.Println("Failed to get executable path:", err)
			return
		}
		AppPath = filepath.Dir(exePath)
	}
	fmt.Println("App root path: " + AppPath)
}

// GetProjectRoot 获取项目的根路径（开发阶段）
func GetProjectRoot() string {
	// 获取当前工作目录
	wd, err := os.Getwd()
	if err != nil {
		return ""
	}

	// 找到包含 go.mod 文件的目录作为项目根路径
	for {
		// 检查当前目录是否包含 go.mod 文件
		if _, err := os.Stat(filepath.Join(wd, "go.mod")); err == nil {
			return wd
		}

		// 到达文件系统的根目录，未找到 go.mod 文件
		if wd == filepath.Dir(wd) {
			return ""
		}

		// 向上一级目录继续查找
		wd = filepath.Dir(wd)
	}
}

// CopyFile 复制单个文件从 src 到 dst
func CopyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	// 确保目标文件有与源文件相同的权限
	err = os.Chmod(dst, os.FileMode(0644))
	if err != nil {
		return err
	}

	return out.Close()
}

// CopyDir 复制整个目录从 src 到 dst
func CopyDir(src string, dst string) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	// 确保目标目录存在
	err = os.MkdirAll(dst, os.ModePerm)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		fi, err := os.Lstat(srcPath)
		if err != nil {
			return err
		}

		switch mode := fi.Mode(); {
		case mode.IsDir():
			// 递归复制子目录
			if err := CopyDir(srcPath, dstPath); err != nil {
				return err
			}
		case mode.IsRegular():
			// 复制文件
			if err := CopyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// Mkdir 文件夹创建
func Mkdir(dst string) error {
	// 确保目标目录存在
	return os.MkdirAll(dst, os.ModePerm)
}
func ReadFile(filename string) (string, error) {
	// 打开文件
	file, err := os.Open(filename)
	if err != nil {
		return "", fmt.Errorf("could not open file: %w", err)
	}
	defer file.Close()

	// 读取文件内容
	content, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("could not read file: %w", err)
	}
	return string(content), nil
}

func Read(filePath string) ([]byte, error) {
	_, err := os.Stat(filePath)
	if err != nil {
		return nil, errorx.ErrDefault.WithError(err).WithMessage("file not exist")
	}
	return os.ReadFile(filePath)
}
func Write(filePath string, data []byte) error {
	return os.WriteFile(filePath, data, 0644)
}

func Append(filePath string, data []byte) error {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// 写入数据
	_, err = file.Write(data)
	return err
}

func CheckDirExistAndCreate(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, 0644); err != nil {
			fmt.Println("mkdir error", err)
		}
	}
}

// ReadEnvToMap 读取环境变量文件并返回键值对映射
func ReadEnvToMap(filePath string) (map[string]string, error) {
	content, err := ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	envMap := make(map[string]string)
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		// 跳过空行和注释行
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// 分割键值对
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			// 去掉可能存在的引号
			value = strings.Trim(value, "\"'")
			envMap[key] = value
		}
	}

	return envMap, nil
}
