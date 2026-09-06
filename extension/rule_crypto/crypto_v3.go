package rule_crypto

import (
	"fmt"
	"os"
)

const (
	magicByte = 0xa5
)

func decryptPacket_v3(inputFile, outputFile string) {
	fileContent, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Println("Error reading input file:", err)
		return
	}

	outFile, err := os.Create(outputFile)
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer outFile.Close()

	// 原数据库的开头91个字节
	buf := fileContent[:91]

	// 原数据库第92到183字节
	nextBuf := fileContent[92:183]

	// 让两部分异或，结果放入开头
	for i := 0; i < len(buf); i++ {
		b := buf[i]
		nb := nextBuf[i]
		// xor 加密
		outFile.Write([]byte{b ^ nb})
	}

	//// 计算被加密文件前91个字节的checksum
	//var dwSum uint32
	//for i := 0; i < 91; i++ {
	//	dwSum += uint32(fileContent[i] ^ fileContent[i+92])
	//}

	// 根据c++代码中的说明，指针跳转到开头offset 91
	for _, b := range fileContent[91:] {
		// 计算checksum unsigned char
		//dwSum += uint32(b)

		// xor加密文件 这里用signed char
		outFile.Write([]byte{b ^ magicByte})
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
