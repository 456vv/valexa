package valexa

import (
	"os"
	"testing"
//	"fmt"
	"bufio"
	"bytes"
	"net/http"
)
func init(){
	os.Chdir("./test/data")
}

//注意
//下面这requestBody是签名的，不能改动他，生成日期是 2017-12-01T09:00:10Z
//在此项目里设置测试有效时间是在 testValidTime 变量里
//如果你测试这个项目发生报错，应该是testValidTime过期了
//1，你可以增加 testValidTime 数值，不能超出 int 类型允许大小的限制。
//2，你可以自己生成一个requestBody内容更新它。
var testValidTime int = 60*60*24*365*10 	//单位为秒,默认10年
var requestBody = []byte(`POST /echo/helloworld?appid=amzn1.ask.skill.594232d8-4095-499b-9ba3-11701107834d&type=1 HTTP/1.1
Content-Type: application/json; charset=utf-8
Accept: application/json
Accept-Charset: utf-8
Signature: L+sBVB8FrP0lvR0MelzmsNivlw6dFYb5p0FU865mSPszIAyHyJ02Eg0GKCATOV25KYC2VpLgVD33tdgSQM56RifFukPnh8jJRCJRP36GHchszW1sIeBsd3ey/MTO7DW4QnYwLtDxhuIaDIifWwSkgT7I2IiqcxhUZPwECYLtzG51HU7Azwj/ECb5gew8wR2NlPKlbdzIO6938pF8veU3JMVlRkFs7dZfLxglcSk+sCcf0qnzCasocMHrO/p70szCN9X2vRt9y3Jur377Xncxb0vz2t5N8yR5KGDctw/J2yZHqhgJtLvKolnxa8wW2CjaqSb6y4mA95VM1JMl5Abjdw==
SignatureCertChainUrl: https://s3.amazonaws.com/echo.api/echo-api-cert-5.pem
Content-Length: 1138
Host: alexatest.xxx.com.cn
Connection: Keep-Alive
User-Agent: HttpClient

 {"session": 
{"sessionId":"SessionId.b67984b5-01e4-4a2d-8cf1-c4692d304dde","application":{"applicationId":"amzn1.ask.skill.594232d8-4095-499b-9ba3-11701107834d"},"attributes":{},"user":{"userId":"amzn1.ask.account.AF3E6FFKSBVMJK6ZOGSEWKEDLD2EXJCZUFJX2ZNWO3R55COS5ZLILGRLM7WTJWKQYPRAOUZWFS2ZZP6ULJRALRA3CVIDCVZ7W5VUMZZMWREW3UWZRIF3XWJMXG5HV4LZZ5ZCDYXKM56BQUKOWVPYP4CWH3TP3SQABYCVATYLG56PNUHV2VON3RAY54LELDRHWBZ2JI6O6VN6LWY","accessToken":null},"new":false},
"request":
{"intent":{"name":"RecipeIntent","slots":{"Item":{"name":"Item","value":"map"}}},"requestId":"EdwRequestId.09c3feaa-a54d-4833-bca0-2f2e4dc1ae5e","type":"IntentRequest","locale":"en-US","timestamp":"2017-12-01T09:00:10Z"},"context":{"AudioPlayer":{"playerActivity":"IDLE"},"System":{"application":{"applicationId":"amzn1.ask.skill.594232d8-4095-499b-9ba3-11701107834d"},"user":{"userId":"amzn1.ask.account.AF3E6FFKSBVMJK6ZOGSEWKEDLD2EXJCZUFJX2ZNWO3R55COS5ZLILGRLM7WTJWKQYPRAOUZWFS2ZZP6ULJRALRA3CVIDCVZ7W5VUMZZMWREW3UWZRIF3XWJMXG5HV4LZZ5ZCDYXKM56BQUKOWVPYP4CWH3TP3SQABYCVATYLG56PNUHV2VON3RAY54LELDRHWBZ2JI6O6VN6LWY"},"device":{"supportedInterfaces":{}}}},"version":"1.0"}`)

func Test_CheckRequestBody_verifyBody(testingT *testing.T){
	//这里就不测试了，除非你可以得到 Signature。
	//SignatureCertChainUrl 的连接也是过期的，amazonaws常更换连接。
	//不提供测试，不代表 verifyBody 方法不能使用。
	return
	var tests = []struct{
		crb checkRequestBody
		err bool
	}{
		{crb:checkRequestBody{R:func() (req *http.Request) {
			req, _ = http.ReadRequest(bufio.NewReader(bytes.NewReader(requestBody)))
			return req
		}()}},
		{crb:checkRequestBody{R:func() (req *http.Request) {
			req, _ = http.ReadRequest(bufio.NewReader(bytes.NewReader(requestBody)))
			req.Header.Set("Signature", "111")
			return req
		}()}, err:true},
		{crb:checkRequestBody{R:func() (req *http.Request) {
			req, _ = http.ReadRequest(bufio.NewReader(bytes.NewReader(requestBody)))
			req.Header.Set("SignatureCertChainUrl", "https://s3.amazonaws.com/echo.api/echo-api-cert-5.pem.xxxxxxxxx")
			return req
		}()}, err:true},
	}
	app := &EchoApplication{
		CertFolder:	"./AmazonCertFile",
	}
	for _, test := range tests {
		_, err := test.crb.verifyBody(app)
		if err != nil && !test.err {
			testingT.Fatal(err)
		}
	}
}

func Test_CheckRequestBody_echoRequest(testingT *testing.T){
	var tests = []struct{
		crb checkRequestBody
		timestamp int
		err bool
	}{
		{crb:checkRequestBody{R:func() (req *http.Request) {
			req, _ = http.ReadRequest(bufio.NewReader(bytes.NewReader(requestBody)))
			return req
		}()},timestamp: 150, err:true},
		{crb:checkRequestBody{R:func() (req *http.Request) {
			req, _ = http.ReadRequest(bufio.NewReader(bytes.NewReader(requestBody)))
			return req
		}()}, timestamp:testValidTime, err:false},
	}
	
	app := &EchoApplication{}

	for _, test := range tests {
		app.ValidReqTimestamp = test.timestamp
		_, err := test.crb.echoRequest(app)
		if err != nil && !test.err {
			testingT.Fatal(err)
		}
	}
}
