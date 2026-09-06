package rule_crypto

import (
	"fmt"
	"general-agent/extension/xzip"
	"testing"
)

func TestDe(t *testing.T) {
	fileName := "example/WagDeviceDb_x86.dat"
	dstFile := "example/db.zip"
	//if doDecodeIpsLib(fileName, dstFile) != 0 {
	//	fmt.Println("Decoding failed!")
	//} else {
	//	fmt.Println("Decoding succeeded!")
	//}
	decryptPacket_v2(fileName, dstFile)

	err := xzip.DeCompress(dstFile, "example/db/", "")
	if err != nil {
		fmt.Println(err)
	}
}
