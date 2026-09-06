package rule_crypto

import (
	"fmt"
	"io"
	"os"
)

func decryptPacket_v2(inputFile, outputFile string) {
	// 打开输入文件
	inFile, err := os.Open(inputFile)
	if err != nil {
		fmt.Println("Error opening input file:", err)
		return
	}
	defer inFile.Close()

	// 创建输出文件
	outFile, err := os.Create(outputFile)
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer outFile.Close()

	// 创建缓冲区
	bufSize := 1024 // 可以根据需要调整缓冲区大小
	buf := make([]byte, bufSize)
	nextBuf := make([]byte, bufSize)

	// 缓存文件的开头 91 个字节
	n, err := inFile.Read(buf[:91])
	if err != nil {
		fmt.Println("Error reading input file:", err)
		return
	}
	if n != 91 {
		fmt.Println("Error: Unexpected EOF while reading input file")
		return
	}
	// 缓存文件的下一个 91 个字节
	n, err = inFile.Read(nextBuf[:1]) // 空一个字节
	n, err = inFile.Read(nextBuf[:91])
	if err != nil {
		fmt.Println("Error reading input file:", err)
		return
	}
	if n != 91 {
		fmt.Println("Error: Unexpected EOF while reading input file")
		return
	}

	// 让两部分异或，结果放入开头
	for i := 0; i < 91; i++ {
		b := buf[i]
		nb := nextBuf[i]
		// xor 加密
		outFile.Write([]byte{b ^ nb})
	}

	// 计算被加密文件前91个字节的checksum
	//var dwSum uint32
	//for i := 0; i < 91; i++ {
	//	dwSum += uint32(buf[i] ^ nextBuf[i])
	//}

	// 读取并解密文件的剩余部分
	inFile.Seek(91, 0)
	for {
		n, err := inFile.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Error reading input file:", err)
			return
		}
		// 根据c++代码中的说明，指针跳转到开头offset 91
		for _, b := range buf[:n] {
			// 计算checksum unsigned char
			//dwSum += uint32(b)
			// xor加密文件 这里用signed char
			outFile.Write([]byte{b ^ magicByte})
		}
	}

	//// 写入checksum
	//checksumBytes := make([]byte, 4)
	//binary.LittleEndian.PutUint32(checksumBytes, dwSum)
	//outFile.Write(checksumBytes)
	//
	//// 计算magicDWORD - dwSum
	//magicDWORDMinusSum := magicDWORD - dwSum
	//magicDWORDBytes := make([]byte, 4)
	//binary.LittleEndian.PutUint32(magicDWORDBytes, magicDWORDMinusSum)
	//outFile.Write(magicDWORDBytes)

	fmt.Println("\nDone!")
}
