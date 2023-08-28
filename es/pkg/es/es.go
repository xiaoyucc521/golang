package es

import (
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
)

var _es *elasticsearch.Client

// Init 创建 es 链接
func Init(addr []string, username, password string) {
	var (
		es  *elasticsearch.Client
		err error
	)

	if es, err = elasticsearch.NewClient(elasticsearch.Config{
		Addresses: addr,
		Username:  username,
		Password:  password,
		Transport: &http.Transport{
			// 所有host的连接池最大连接数量，默认无穷大
			MaxIdleConns: 100,
			// 每个host的连接池最大空闲连接数,默认2
			MaxIdleConnsPerHost: 10,
			// 每个host的最大连接数量
			MaxConnsPerHost: 2,
			// 空闲连接在连接池中保留的时间。
			IdleConnTimeout: time.Hour,
			// 限制读取response header的时间
			ResponseHeaderTimeout: time.Second,
			DialContext:           (&net.Dialer{Timeout: time.Second}).DialContext,
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
			},
		},
	}); err != nil {
		log.Fatalln("es 连接出错", err)
	}

	_es = es
}

func NewESClient() *elasticsearch.Client {
	es := _es
	return es
}
