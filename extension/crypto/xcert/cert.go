package xcert

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"general-agent/extension/xcmd"
	"general-agent/extension/xfile"
	"golang.org/x/crypto/pkcs12"
	"strings"
)

func ToPem(certBinary, keyBinary, password string) ([]byte, []byte, error) {
	certBytes, err := base64.StdEncoding.DecodeString(certBinary)
	if err != nil {
		return nil, nil, err
	}
	keyBytes, err := base64.StdEncoding.DecodeString(keyBinary)
	if err != nil {
		return nil, nil, err
	}
	if keyBinary != "" {
		return certBytes, keyBytes, nil
	} else {
		// 分离证书和密钥
		//c, k, err := Pkcs12ToPEM(certBytes, password)
		//if err != nil {
		//	return nil, nil, err
		//}
		err = xfile.Write("/tmp/crt.p12", certBytes)
		if err != nil {
			return nil, nil, err
		}
		// 证书内容
		output, err := xcmd.ExecCombinedOutput("openssl", "pkcs12", "-in", "/tmp/crt.p12", "-passin", "pass:"+password, "-clcerts", "-nokeys", "-out", "/tmp/crt.pem", "-nodes")
		if err != nil {
			return nil, nil, errors.New(string(output))
		}
		crtRaw, err := xcmd.ExecCombinedOutput("cat", "/tmp/crt.pem")
		if err != nil {
			return nil, nil, errors.New(string(crtRaw))
		}
		// 私钥内容
		output, err = xcmd.ExecCombinedOutput("openssl", "pkcs12", "-in", "/tmp/crt.p12", "-passin", "pass:"+password, "-nocerts", "-out", "/tmp/private_key.pem", "-nodes")
		if err != nil {
			return nil, nil, errors.New(string(output))
		}
		keyRaw, err := xcmd.ExecCombinedOutput("cat", "/tmp/private_key.pem")
		if err != nil {
			return nil, nil, errors.New(string(keyRaw))
		}
		return crtRaw, keyRaw, nil
	}
}

// ParseCrt 读取.cer、.crt、.pem证书
func ParseCrt(certData []byte) (*x509.Certificate, error) {
	// 将PEM格式的证书数据解码为x509.Certificate对象
	block, _ := pem.Decode(certData)
	//if len(rest) > 0 {
	//	return nil, errors.New("unexpected trailing data after PEM block")
	//}

	if block == nil || block.Type != "CERTIFICATE" {
		return nil, errors.New("failed to decode PEM block containing the certificate")
		//log.Fatal("Failed to decode PEM block containing the certificate")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}
	return cert, nil
}

// ParsePk12 解析PK12证书
func ParsePk12(pfxData []byte, password string) (*x509.Certificate, error) {
	blocks, err := pkcs12.ToPEM(pfxData, password)
	if err != nil {
		return nil, err
	}

	// 遍历所有解析出的块
	for _, block := range blocks {
		// block.Type == "CERTIFICATE" or "PRIVATE KEY"
		if block.Type == "CERTIFICATE" {
			cert, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				return nil, err
			}
			return cert, err
		}
	}
	return nil, nil
}

func Pkcs12ToPEM(pfxData []byte, password string) ([]byte, []byte, error) {
	blocks, err := pkcs12.ToPEM(pfxData, password)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to convert PFX to PEM: %w", err)
	}
	var cert, privateKey []byte
	// 遍历所有解析出的块
	for _, block := range blocks {
		block.Headers = nil
		if block.Type == "CERTIFICATE" && len(cert) == 0 {
			cert = pem.EncodeToMemory(block)
		}
		if block.Type == "PRIVATE KEY" && len(privateKey) == 0 {
			privateKey = pem.EncodeToMemory(block)
		}
	}
	// 检查是否找到了证书和私钥
	if len(cert) == 0 {
		return nil, nil, errors.New("no certificate found in PFX data")
	}
	if len(privateKey) == 0 {
		return nil, nil, errors.New("no private key found in PFX data")
	}

	return cert, privateKey, nil
}

func SplitPk12(p12Bytes []byte, password string) ([]byte, []byte, error) {
	// 解码 P12 文件
	privateKey, certificate, err := pkcs12.Decode(p12Bytes, password)
	if err != nil {
		fmt.Println("Error decoding PKCS #12:", err)
		return nil, nil, err
	}

	// 将私钥编码为 PEM 格式
	keyPEMBlock := pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey.(*rsa.PrivateKey)),
	}
	keyPEM := pem.EncodeToMemory(&keyPEMBlock)

	// 将证书编码为 PEM 格式
	certPEMBlock := pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certificate.Raw,
	}
	certPEM := pem.EncodeToMemory(&certPEMBlock)

	return certPEM, keyPEM, nil
}

// ValidateCertAndKey 验证证书和密钥是否匹配
func ValidateCertAndKey(certPEM, keyPEM []byte) error {
	// 解析证书
	cert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		return fmt.Errorf("无法解析证书或私钥: %v", err)
	}

	// 解析证书链
	x509Cert, err := x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		return fmt.Errorf("无法解析 x509 证书: %v", err)
	}

	// 检查公钥是否匹配私钥
	switch pub := x509Cert.PublicKey.(type) {
	case *rsa.PublicKey:
		if pub.N.Cmp(cert.PrivateKey.(*rsa.PrivateKey).N) != 0 {
			return fmt.Errorf("RSA 公钥和私钥不匹配")
		}
	case *ecdsa.PublicKey:
		if pub.X.Cmp(cert.PrivateKey.(*ecdsa.PrivateKey).PublicKey.X) != 0 ||
			pub.Y.Cmp(cert.PrivateKey.(*ecdsa.PrivateKey).PublicKey.Y) != 0 {
			return fmt.Errorf("ECDSA 公钥和私钥不匹配")
		}
	default:
		return fmt.Errorf("不支持的公钥类型")
	}

	return nil
}

func FormatPkixName(attr pkix.Name) string {
	var result []string
	for _, name := range attr.Names {
		label := fmt.Sprintf("%s=%v", getAttributeTypeName(name.Type.String()), name.Value)
		result = append(result, label)
	}
	return strings.Join(result, ",")
}

var attributeTypeNames = map[string]string{
	"2.5.4.3":              "CN",                   // 通用名称
	"2.5.4.5":              "SERIALNUMBER",         // 序列号
	"2.5.4.6":              "C",                    // 国家/地区
	"2.5.4.7":              "L",                    // 城市/地点
	"2.5.4.8":              "ST",                   // 省/自治区/直辖市
	"2.5.4.9":              "STREET",               // 街道地址
	"2.5.4.10":             "O",                    // 组织名称
	"2.5.4.11":             "OU",                   // 组织单位名称
	"2.5.4.12":             "T",                    // 职务
	"2.5.4.17":             "POSTALCODE",           // 邮政编码
	"2.5.4.23":             "DC",                   // 域组件
	"1.2.840.113549.1.9.1": "E",                    //邮箱
	"1.2.840.113549.1.9.2": "UNSTRUCTURED NAME",    // 非结构化名称
	"1.2.840.113549.1.9.3": "UNSTRUCTURED ADDRESS", // 非结构化地址
	"1.2.840.113549.1.9.4": "D",                    // 描述
}

func getAttributeTypeName(oid string) string {
	v, ok := attributeTypeNames[oid]
	if ok {
		return v
	}
	return oid
}
