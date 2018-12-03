package valexa
	
import (
	"testing"
	"net/http"
	"net/http/httptest"
	"bufio"
	"bytes"
	"os"
	"io"
	"encoding/json"
//	"fmt"
)

func init(){
	os.Chdir("./test/data")
}


func Test_Alexa_1(testingT *testing.T){
	echoApp := &EchoApplication{
		ValidReqTimestamp	: testValidTime,
		CertFolder			: "./AmazonCertFile",
	}
	alexa := &Alexa{
		AppIdAttr	: "appid",
	}
	alexa.SetEchoApp("amzn1.ask.skill.594232d8-4095-499b-9ba3-11701107834d", echoApp)
	
	w := httptest.NewRecorder()
	r, _ := http.ReadRequest(bufio.NewReader(bytes.NewReader(requestBody)))
	alexa.ServeHTTP(w, r)
}

func EchoHTTPHandler(w http.ResponseWriter, r *http.Request) {
	inf := r.Context().Value(r.URL.Path)
	echoIntent, ok := inf.(*EchoIntent)
	if !ok {
		http.Error(w, "该路径没有绑定有效的处理程序", 500)
	}
	echoIntent.Response.OutputSpeech("Hello world from my new Echo test app!").SimpleCard("Hello World", "This is a test card.")
	io.WriteString(w, echoIntent.Response.String())
}

func Test_Alexa_3(testingT *testing.T){
	//这里也不测试了，原因是requestBody的参数过期了
	//有空再更新了
	return
	echoApp := &EchoApplication{
		ValidReqTimestamp	: testValidTime,
		CertFolder			: "./AmazonCertFile",
		HandleFunc     		: EchoHTTPHandler,
	}
	alexa := &Alexa{
		AppIdAttr	: "appid",
	}
	alexa.SetEchoApp("amzn1.ask.skill.594232d8-4095-499b-9ba3-11701107834d", echoApp)
	
	w := httptest.NewRecorder()
	r, _ := http.ReadRequest(bufio.NewReader(bytes.NewReader(requestBody)))
	alexa.ServeHTTP(w, r)
	
	var echoResp *EchoResponse
	err := json.NewDecoder(w.Body).Decode(&echoResp)
	if err != nil {
		testingT.Fatalf("错误:%v", err)
	}
	if echoResp.Response.Card.Title != "Hello World" {
		testingT.Fatalf("错误的:%s", echoResp.Response.Card.Title)
	}
	if echoResp.Response.Card.Content != "This is a test card." {
		testingT.Fatalf("错误的:%s", echoResp.Response.Card.Content)
	}
	if echoResp.Response.OutputSpeech.Text != "Hello world from my new Echo test app!" {
		testingT.Fatalf("错误的:%s", echoResp.Response.Card.Content)
	}
}

