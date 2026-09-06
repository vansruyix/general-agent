package xzip

import (
	"general-agent/app/model/xconst"
	"testing"
)

func TestCompressDirAndEncrypt(t *testing.T) {
	err := CompressDir("C:\\Docs\\My\\jetbrains-agent\\test_dir.zip", xconst.ExportZipPassword, "C:\\Docs\\My\\jetbrains-agent")
	if err != nil {
		t.Error(err)
	}
}
func TestDeCompressAndDecrypt(t *testing.T) {
	err := DeCompress("C:\\Docs\\My\\jetbrains-agent\\test_dir.zip", "C:\\Docs\\My\\jetbrains-agent", xconst.ExportZipPassword)
	if err != nil {
		t.Error(err)
	}
}
