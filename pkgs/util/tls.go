package util

import (
	"crypto/tls"
	"io/ioutil"
	"log"
	"net"

	"golang.org/x/net/http2"
)

// 用于获取TLS配置
func GetTLSConfig(certPemPath, certKeyPath string) *tls.Config {
	var certKeyPair *tls.Certificate
	// 读取 server.pem 凭证文件
	cert, _ := ioutil.ReadFile(certPemPath)
	// 读取 server.key 凭证文件
	key, _ := ioutil.ReadFile(certKeyPath)
	// tls.X509KeyPair 从一对PEM编码的数据中解析公钥/私钥对
	// 成功则返回公钥/私钥对
	pair, err := tls.X509KeyPair(cert, key)
	if err != nil {
		log.Println("TLS KeyPair err: %v\n", err)
	}

	certKeyPair = &pair

	// 用于处理从证书凭证文件 PEM 最终获取tls.Config作为HTTP/2的使用参数
	return &tls.Config{
		// tls.Certificate 返回一个或多个证书
		// 实质我们解析PEM调用的X509KeyPair的函数声明就是
		// func X509KeyPair(certPEMBlock, keyPEMBlock []byte) (Certificate, error)
		// 返回值就是Certificate
		Certificates: []tls.Certificate{*certKeyPair},
		// NextProtoTLS 是谈判期间的NPN/ALPN协议 用于HTTP/2的TLS设置
		NextProtos: []string{http2.NextProtoTLS},
	}
}

func NewTLSListener(inner net.Listener, config *tls.Config) net.Listener {
	// tls.NewListener 会创建一个 Listener
	// -  来自内部Listener的监听器
	// -  tls.Config 必须包含至少一个证书
	return tls.NewListener(inner, config)
}
