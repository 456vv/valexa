package main

import (
	"github.com/456vv/valexa"
	"fmt"
	"net/http"
	"io"
)


func main(){
	fmt.Println(valexa.Version)

	echoApp := &alexa.EchoApplication{
		ValidReqTimestamp	: 150,
		ValidReqTimestamp	: testValidTime,
		CertFolder			: "./AmazonCertFile",
		HandleFunc        	: EchoHTTPHandler,
	}
	alexa := &valexa.Alexa{}
	alexa.SetEchoApp("amzn1.ask.skill.594232d8-4095-499b-9ba3-11701107834d", echoApp)
	alexa.ServeHTTP(w http.ResponseWriter, r *http.Request)
}

func EchoIntentHandler(echoReq *valexa.EchoRequest, echoResp *valexa.EchoResponse) {
	echoResp.OutputSpeech("Hello world from my new Echo test app!").Card("Hello World", "This is a test card.")
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
