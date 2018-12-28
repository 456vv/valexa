# valexa [![Build Status](https://travis-ci.org/456vv/valexa.svg?branch=master)](https://travis-ci.org/456vv/valexa)
golang valexa，亚马逊 Echo Alexa 自定义版本

# **列表：**
```go
type Alexa struct{                                                                              // Alexa服务
    AppIdAttr       string                                                                          // 程序的名称，如：appid
}
    func (T *Alexa) SetEchoApp(alexaAppId string, app *EchoApplication)                             // 设置echoApp
    func (T *Alexa) ServeHTTP(w http.ResponseWriter, r *http.Request)                               // 服务器启用
func NewEchoResponse() *EchoResponse                                                            // Echo响应
    func (T *EchoResponse) SetEndSession(ok bool)                                                   // 结束对话
    func (T *EchoResponse) OutputSpeech(text string) *EchoResponse                                  // 输出语音
    func (T *EchoResponse) OutputSpeechSSML(text string) *EchoResponse                              // 输出语音（有语气）
    func (T *EchoResponse) SimpleCard(title string, content string) *EchoResponse                   // 纯文本卡片
    func (T *EchoResponse) StandardCard(title string, text string, smallImg string, largeImg string) *EchoResponse      // 带图文本卡片
    func (T *EchoResponse) LinkAccountCard() *EchoResponse                                          // 登录卡片
    func (T *EchoResponse) RepromptText(text string) *EchoResponse                                  // 回应
    func (T *EchoResponse) RepromptSSML(text string) *EchoResponse                                  // 回应（有语气）
    func (T *EchoResponse) String() string                                                          // 字符串
type EchoRequest struct {                                                                       // Echo请求
    Context EchoRequestContext                                                                      // 上下文
    Request EchoRequestRequest                                                                      // 请求
    Session EchoRequestSession                                                                      // 会话
    Version string                                                                                  // 版本
}
    func (T *EchoRequest) VerifyTimestamp(second int) bool                                              // 验证请求时间
    func (T *EchoRequest) GetLocale() string															// 语言类型
    func (T *EchoRequest) GetApplicationID() string                                                     // 程序ID
    func (T *EchoRequest) GetSessionID() string                                                         // 会话ID
    func (T *EchoRequest) GetUserID() string                                                            // 用户ID
    func (T *EchoRequest) GetAccessToken() string                                                       // 用户Token
    func (T *EchoRequest) GetRequestType() string                                                       // 请求类型
    func (T *EchoRequest) GetIntentName() string                                                        // 意图名称
    func (T *EchoRequest) GetSlots() (slots map[string]*EchoRequestRequestIntentSlot, err error)        // 所有类型意图对象（含空值）
	func (T *EchoRequest) GetValueSlots() (slots map[string]*EchoRequestRequestIntentSlot, err error)	// 所有类型意图对象（不含空值）
    func (T *EchoRequest) GetSlotNames() (names []string)                                               // 所有意图名称
    func (T *EchoRequest) GetSlotValue(slotName string) (val string, err error)                         // 意图值
    func (T *EchoRequest) GetSlot(slotName string) (slot *EchoRequestRequestIntentSlot, err error)      // 意图对象
type EchoIntent struct{                                                                             // Echo意图
    Request     *EchoRequest                                                                            // Echo请求
    Response    *EchoResponse                                                                           // Echo响应
    App         *EchoApplication                                                                        // Echo配置
}
func (T *EchoIntent) ServeHTTP(w http.ResponseWriter, r *http.Request)
type EchoApplication struct {                                                                       // Echo程序的配置
    Version                 []string                                                                    // 支持的版本号
    IsDevelop               bool                                                                      	// 调试
    ValidReqTimestamp       int                                                                         // 有效时间，秒为单位
    CertFolder              string                                                                      // 证书保存目录
    HandleFunc              http.HandlerFunc                                                            // 原生处理函数
}
```