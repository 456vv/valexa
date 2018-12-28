package valexa

import (
	"net/http"
	"fmt"
	"sync"
)


type Alexa struct{
	//如：https://www.xxx.com/?appid=amzn1.ask.skill.594232d8-4095-499b-9ba3-11701107834s&type=1
	//这个 appid 是程序属性名称，其中的值是 amzn1.ask.skill.594232d8-4095-499b-9ba3-11701107834s
	//这个 type 是程序类型属性名称，暂时仅支持 Echo，值也就是 1。其它值不支持
	//在Alexa后台，Configuration->Global Fields->Endpoint->Service Endpoint Type->选择HTTPS
	//Default 填入 https://www.xxx.com/?appid=amzn1.ask.skill.594232d8-4095-499b-9ba3-11701107834s&type=1
	AppIdAttr		string									//程序的名称，如：appid
	echoApps		*sync.Map								//map[ID名]*EchoApplication
}

//设置Echo程序
//	alexaAppId string		Alexa程序ID名
//	app *EchoApplication	Echo对象
func (T *Alexa) SetEchoApp(alexaAppId string, app *EchoApplication){
	if T.echoApps == nil {
		T.echoApps = &sync.Map{}
	}
	if app == nil {
		T.echoApps.Delete(alexaAppId)
		return
	}
	T.echoApps.Store(alexaAppId, app)
}

//服务处理
//	w http.ResponseWriter	http响应对象
// 	r *http.Request			http请求对象
func (T *Alexa) ServeHTTP(w http.ResponseWriter, r *http.Request){
	var(
		query 	= r.URL.Query()
		appId 	= query.Get(T.AppIdAttr)
	)
	if appId == "" {
		http.Error(w, fmt.Sprintf("valexa: 请求数据参数不完整，缺少程序ID名（AppIdAttr=%s）", appId), 400)
		return
	}
	
	if T.echoApps == nil {
		http.Error(w, fmt.Sprintf("valexa: 服务器配置为空，请使用 .SetEchoApp 方法来设置配置。"), 400)
		return
	}
	inf, ok := T.echoApps.Load(appId)
	if !ok {
		http.Error(w, fmt.Sprintf("valexa: 该程序Id（AppIdAttr=%s）不是有效的", appId), 400)
		return
	}
	echoApp := inf.(*EchoApplication)

	crb := &checkRequestBody{R:r}
	//如果调试ID不相等，则验证Body
	if !echoApp.IsDevelop {
		//验证数据
		_, err := crb.verifyBody(echoApp)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
	}
	echoReq, err := crb.echoRequest(echoApp)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	//处理语音服务
	eialexa := &EchoIntent{
		Request		: echoReq,
		Response	: NewEchoResponse(),
		App			: echoApp,
	}
	eialexa.ServeHTTP(w, r)
}


//strSliceContains 从切片中查找匹配的字符串
//  参：
//      ss []string     切片
//      s string        需要从切片中查找的字符
func strSliceContains(ss []string, s string) bool {
	for _, v := range ss {
		if v == s {
			return true
		}
	}
	return false
}


