package valexa

import (
	"net/http"
)

//Echo程序的配置
type EchoApplication struct {
	Version					[]string									// 支持的版本号
    DevId                   string                                      // 调试ID，这是在go程序里的调试标识符，有关：Alexa.AttrDevId
	ValidReqTimestamp		int											// 有效时间，秒为单位
	CertFolder				string										// 证书保存目录
	HandleFunc     			http.HandlerFunc							// 原生处理函数
}