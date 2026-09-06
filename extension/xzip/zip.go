package xzip

import (
	"github.com/alexmullins/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Compress 压缩文件或目录到指定目标
// dest: 压缩文件保存路径
// password: 压缩密码，为空表示不加密
// filePaths: 要压缩的文件或目录路径
func Compress(dest string, password string, filePaths ...string) error {
	// 获取目标文件的绝对路径，用于排除自身
	destAbs, err := filepath.Abs(dest)
	if err != nil {
		return err
	}

	// 创建压缩文件
	zipFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	// 使用支持加密的zip.Writer
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// 添加文件到压缩包
	for _, path := range filePaths {
		err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// 获取当前文件的绝对路径
			fileAbs, err := filepath.Abs(filePath)
			if err != nil {
				return err
			}

			// 排除目标压缩文件本身，避免递归包含
			if fileAbs == destAbs {
				return nil
			}

			// 创建文件头
			header, err := zip.FileInfoHeader(info)
			if err != nil {
				return err
			}

			// 设置压缩包中的文件路径（保持相对路径结构）
			relPath, err := filepath.Rel(filepath.Dir(path), filePath)
			if err != nil {
				return err
			}
			header.Name = relPath

			// 如果是目录，需要在名称末尾添加斜杠
			if info.IsDir() {
				header.Name += "/"
			} else {
				// 设置压缩方法
				header.Method = zip.Deflate
			}

			// 设置密码（如果有）
			if password != "" {
				header.SetPassword(password)
			}

			// 创建压缩文件条目
			writer, err := zipWriter.CreateHeader(header)
			if err != nil {
				return err
			}

			// 如果是目录，没有内容需要写入
			if info.IsDir() {
				return nil
			}

			// 打开并写入文件内容
			file, err := os.Open(filePath)
			if err != nil {
				return err
			}
			defer file.Close()

			_, err = io.Copy(writer, file)
			return err
		})

		if err != nil {
			return err
		}
	}

	return nil
}

// CompressDir 压缩整个目录到指定目标
// dest: 压缩文件保存路径
// password: 压缩密码，为空表示不加密
// dir: 要压缩的目录路径
func CompressDir(dest string, password string, dir string) error {
	return Compress(dest, password, dir)
}

// DeCompress 解压ZIP文件到指定目录
// zipFile: 要解压的ZIP文件路径
// dest: 解压目标目录
// password: 解压密码，为空表示无密码
func DeCompress(zipFile, dest, password string) error {
	// 打开ZIP文件
	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer reader.Close()

	// 获取目标目录的绝对路径，用于路径遍历检查
	destAbs, err := filepath.Abs(dest)
	if err != nil {
		return err
	}

	// 遍历ZIP文件中的每个文件/目录
	for _, file := range reader.File {
		// 设置密码（如果有）
		if password != "" {
			file.SetPassword(password)
		}

		// 将 ZIP 中的路径分隔符统一为 '/'，再转换为系统分隔符
		name := strings.ReplaceAll(file.Name, `\`, `/`) // 先替换反斜杠
		name = filepath.Clean(name)                     // 清理路径

		// 构建目标路径
		targetPath := filepath.Join(dest, name)

		// 检查路径遍历漏洞
		targetAbs, err := filepath.Abs(targetPath)
		if err != nil {
			return err
		}
		if !strings.HasPrefix(targetAbs, destAbs+string(os.PathSeparator)) && targetAbs != destAbs {
			return &os.PathError{Op: "extract", Path: file.Name, Err: os.ErrPermission}
		}

		// 判断是否是目录：检查是否以 '/' 结尾 或 FileInfo().IsDir()
		isDir := file.FileInfo().IsDir()
		if !isDir && strings.HasSuffix(name, "/") {
			isDir = true
		}

		if isDir {
			// 创建目录
			err = os.MkdirAll(targetPath, os.ModePerm)
			if err != nil {
				return err
			}
			continue
		}

		// 确保父目录存在
		if err = os.MkdirAll(filepath.Dir(targetPath), os.ModePerm); err != nil {
			return err
		}

		// 打开ZIP中的文件
		fileReader, err := file.Open()
		if err != nil {
			return err
		}

		// 创建目标文件
		targetFile, err := os.Create(targetPath)
		if err != nil {
			fileReader.Close()
			return err
		}

		// 复制文件内容
		_, err = io.Copy(targetFile, fileReader)
		fileReader.Close()
		targetFile.Close()

		if err != nil {
			return err
		}

		// 设置文件权限
		if err = os.Chmod(targetPath, file.Mode()); err != nil {
			return err
		}

		// 设置文件修改时间
		// 注意: 在Windows上修改时间可能不完全准确
		if err = os.Chtimes(targetPath, file.ModTime(), file.ModTime()); err != nil {
			return err
		}
	}

	return nil
}
