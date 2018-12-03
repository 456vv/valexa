package valexa
	
import (
	"net/http"
	"context"
	"io"
)

//Echo意图处理
type EchoIntent struct{
	Request		*EchoRequest		//Echo请求
	Response	*EchoResponse		//Echo响应
	App			*EchoApplication	//Echo配置
}

//服务处理
//	w http.ResponseWriter	http响应对象
// 	r *http.Request			http请求对象
func (T *EchoIntent) ServeHTTP(w http.ResponseWriter, r *http.Request){
	if T.App.HandleFunc != nil {
		r = r.WithContext(context.WithValue(r.Context(), r.URL.Path, T))
		T.App.HandleFunc(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	io.WriteString(w, T.Response.String())
}