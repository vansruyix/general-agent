package rule_crypto

import (
	"fmt"
	"io"
	"os"
)

// 验证码，用于检查文件是否已经正确加密
const validationNumber = 0xFEDCBA98

// 解密算法中使用的异或密钥
const byMagic = 0xA5

// doDecodeIpsLib 函数用于解密文件
func doDecodeIpsLib(fileName, dstFile string) int {
	var uTemp int
	var uSum int
	var buf [1024]byte
	var nextBuf [100]byte

	// 打开待解密的文件
	psFile, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return -1
	}
	defer psFile.Close()

	// 创建输出文件
	pdFile, err := os.Create(dstFile)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return -1
	}
	defer pdFile.Close()

	// 获取文件大小
	fileInfo, err := psFile.Stat()
	if err != nil {
		fmt.Println("Error getting file size:", err)
		return -1
	}
	fileSize := int(fileInfo.Size())

	// 读取待解密文件的最后两个双字
	if _, err := psFile.Seek(int64(fileSize-2*4), 0); err != nil {
		fmt.Println("Error reading file:", err)
		return -1
	}

	temp := make([]byte, 4)
	if _, err := psFile.Read(temp); err != nil {
		fmt.Println("Error reading file:", err)
		return -1
	}
	uSum = int(temp[0])<<24 | int(temp[1])<<16 | int(temp[2])<<8 | int(temp[3])
	if _, err := psFile.Read(temp); err != nil {
		fmt.Println("Error reading file:", err)
		return -1
	}
	uTemp = int(temp[0])<<24 | int(temp[1])<<16 | int(temp[2])<<8 | int(temp[3])

	// 验证文件是否已正确加密
	if uSum+uTemp != validationNumber {
		fmt.Printf("file encode error, Expected to be %d but actually %d\n", validationNumber, uSum+uTemp)
		//return -1
	}

	// 定位文件指针到文件开头
	if _, err := pdFile.Seek(0, 0); err != nil {
		fmt.Println("Error seeking file:", err)
		return -1
	}
	if _, err := psFile.Seek(0, 0); err != nil {
		fmt.Println("Error seeking file:", err)
		return -1
	}

	uTemp = 0

	// 读取并解密文件的前 91 个字节
	if _, err := io.ReadFull(psFile, buf[:91]); err != nil {
		fmt.Println("Error reading file:", err)
		return -1
	}
	if _, err := io.ReadFull(psFile, nextBuf[:1]); err != nil {
		fmt.Println("Error reading file:", err)
		return -1
	}
	if _, err := io.ReadFull(psFile, nextBuf[:91]); err != nil {
		fmt.Println("Error reading file:", err)
		return -1
	}

	// 解密数据
	uTemp += xorDecrypt(buf[:], 91)
	for i := 0; i < 91; i++ {
		buf[i] ^= nextBuf[i]
	}

	// 写入解密后的数据到输出文件
	if _, err := pdFile.Write(buf[:91]); err != nil {
		fmt.Println("Error writing file:", err)
		return -1
	}

	uWrite := 91

	// 继续读取文件并解密
	for {
		nRet, err := psFile.Read(buf[:])
		if err != nil && err != io.EOF {
			fmt.Println("Error reading file:", err)
			return -1
		}
		if nRet == 0 {
			break
		}
		uWrite += nRet
		if uWrite+2*4 >= fileSize {
			nValidLen := fileSize - (uWrite - nRet) - 2*4
			if nValidLen < 0 {
				break
			}
			uTemp += xorDecrypt(buf[:nValidLen], nValidLen)
			if _, err := pdFile.Write(buf[:nValidLen]); err != nil {
				fmt.Println("Error writing file:", err)
				return -1
			}
			break
		} else {
			uTemp += xorDecrypt(buf[:nRet], nRet)
			if _, err := pdFile.Write(buf[:nRet]); err != nil {
				fmt.Println("Error writing file:", err)
				return -1
			}
		}
	}

	// 校验解密后的数据
	if uTemp != uSum {
		fmt.Printf("file decode after error, Expected to be %d but actually %d\n", uSum, uTemp)
		//return -1
	}

	return 0
}

// xorDecrypt 函数用于异或解密数据
func xorDecrypt(pBuf []byte, dwLen int) int {
	uSum := 0
	for i := 0; i < dwLen; i++ {
		pBuf[i] ^= byMagic
		uSum += int(pBuf[i] & 0xff)
	}
	return uSum
}
