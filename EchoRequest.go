package valexa

import (
	"time"
	"fmt"
	
)



//提供AudioPlayer接口当前状态的对象
//https://developer.amazon.com/docs/custom-skills/request-and-response-json-reference.html#audioplayer-object
type EchoRequestContextAudioPlayer struct {
 	//发送请求时的轨道偏移量（以毫秒为单位）。
 	//如果曲目在开头，这是0。AudioPlayer当您的技能是最近在设备上播放音频的技能时，这只会包含在对象中。
	OffsetInMilliseconds int    										`json:"offsetInMilliseconds"`

	//音频播放的最后已知状态。
	//IDLE：没有什么，没有入队的项目。
	//PAUSED：流已暂停。
	//PLAYING：流正在播放。
	//BUFFER_UNDERRUN：缓冲区不足
	//FINISHED：流完了。STOPPED：流被中断。
	PlayerActivity       string 										`json:"playerActivity"`

	//音频流的不透明标记。发送Play指令时提供此标记。
	//AudioPlayer当您的技能是最近在设备上播放音频的技能时，这只会包含在对象中。
	Token                string 										`json:"token"`
}

type EchoRequestContextSystemApplication struct {
	ApplicationID string 												`json:"applicationId"`
}

type EchoRequestContextSystemDeviceSupportedInterfaces struct {
	AudioPlayer struct{} 												`json:"AudioPlayer"`
}

type EchoRequestContextSystemDevice struct {
	DeviceID            string 												`json:"deviceId"`
	SupportedInterfaces EchoRequestContextSystemDeviceSupportedInterfaces	`json:"supportedInterfaces"`
}

type EchoRequestContextSystemUserPermissions struct {
	ConsentToken string 												`json:"consentToken"`
}

type EchoRequestContextSystemUser struct {
	AccessToken string 													`json:"accessToken"`
	Permissions EchoRequestContextSystemUserPermissions					`json:"permissions"`
	UserID 		string 													`json:"userId"`
}

type EchoRequestContextSystem struct {
	APIAccessToken 	string 												`json:"apiAccessToken"`
	APIEndpoint   	string 												`json:"apiEndpoint"`
	Application    	EchoRequestContextSystemApplication					`json:"application"`
	Device			EchoRequestContextSystemDevice						`json:"device"`
	User  			EchoRequestContextSystemUser						`json:"user"`
}

//Alexa服务和设备的当前状态的信息。用于在会话（的上下文中发送的请求LaunchRequest和IntentRequest）时，
//context对象将复制user和application其也可在信息session对象。
//https://developer.amazon.com/docs/custom-skills/request-and-response-json-reference.html#context-object
type EchoRequestContext struct {
	AudioPlayer EchoRequestContextAudioPlayer							`json:"AudioPlayer"`
	System  	EchoRequestContextSystem		 						`json:"System"`

}

type EchoRequestRequestError struct {
	Message string 														`json:"message"`
	Type    string 														`json:"type"`
}

type EchoRequestRequestIntentSlot struct {
	ConfirmationStatus string   										`json:"confirmationStatus"`
	Resolutions        struct{} 										`json:"resolutions"`
	Name               string   										`json:"name"`
	Value              string   										`json:"value"`
}

type EchoRequestRequestIntent struct {
	ConfirmationStatus string 											`json:"confirmationStatus"`
	Name               string 											`json:"name"`
	Slots              map[string]*EchoRequestRequestIntentSlot			`json:"slots"`
}

type EchoRequestRequest	struct {
	DialogState string 													`json:"dialogState"`
	Error       EchoRequestRequestError									`json:"error"`
	Intent  	EchoRequestRequestIntent								`json:"intent"`
	Locale    	string 													`json:"locale"`
	Reason    	string 													`json:"reason"`
	RequestID 	string 													`json:"requestId"`
	Timestamp 	string 													`json:"timestamp"`
	Type      	string 													`json:"type"`
}

type EchoRequestSessionApplication struct {
	ApplicationID string 												`json:"applicationId"`
}

type EchoRequestSessionUserPermissions struct {
	ConsentToken string 												`json:"consentToken"`
}
type EchoRequestSessionUser struct {
	AccessToken string 													`json:"accessToken"`
	Permissions EchoRequestSessionUserPermissions						`json:"permissions"`
	UserID 		string 													`json:"userId"`
}
type EchoRequestSession struct {
	Application EchoRequestSessionApplication							`json:"application"`
	Attributes	struct{} 												`json:"attributes"`
	New       	bool     												`json:"new"`
	SessionID 	string   												`json:"sessionId"`
	User      	EchoRequestSessionUser									`json:"user"`
}
type EchoRequest struct {
	Context	EchoRequestContext											`json:"context"`
	Request EchoRequestRequest											`json:"request"`
	Session EchoRequestSession											`json:"session"`
	Version string 														`json:"version"`
}

//验证请求是否有效
//	second int	允许过时多久？
//	bool		true有效，false无效
func (T *EchoRequest) VerifyTimestamp(second int) bool {
	reqTimestamp, err := time.Parse(time.RFC3339, T.Request.Timestamp)
	if err != nil {
		return false
	}
	if time.Since(reqTimestamp) < time.Duration(second)*time.Second {
		return true
	}
	return false
}

//读取语言类型
//	string	语言类型
func (T *EchoRequest) GetLocale() string {
	return T.Request.Locale
}

//读取程序ID
//	string	程序ID
func (T *EchoRequest) GetApplicationID() string {
	return T.Session.Application.ApplicationID
}

//读取本次会话的ID。如果用户结束会话，再发起会话，这个ID将是一个全新ID。
//	string	会话ID
func (T *EchoRequest) GetSessionID() string {
	return T.Session.SessionID
}

//读取用户的ID
//	string	用户ID
func (T *EchoRequest) GetUserID() string {
	return T.Session.User.UserID
}

//读取用户的令牌
//	string	令牌
func (T *EchoRequest) GetAccessToken() string {
	return T.Session.User.AccessToken
}


//读取意图类型，有效：LaunchRequest, IntentRequest, SessionEndedRequest
//	string	意图类型
func (T *EchoRequest) GetRequestType() string {
	return T.Request.Type
}

//读取意图名称，如标注符：{Map}
//	string	意图名称
func (T *EchoRequest) GetIntentName() string {
	if T.GetRequestType() == "IntentRequest" {
		return T.Request.Intent.Name
	}
	return ""
}

//读取出所有Solt对象，所有值，包含空值
//	slots	所有跟踪名称和值
//	err		错误，如果不是意图类型请求，调用这个方法，将会返回错误
func (T *EchoRequest) GetSlots() (slots map[string]*EchoRequestRequestIntentSlot, err error) {
	if T.GetRequestType() != "IntentRequest" {
		return  nil, fmt.Errorf("valexa: 此意图类型不支持读取插槽，仅支持IntentRequest类型。")
	}

	slots = T.Request.Intent.Slots
	if slots == nil {
		return  nil, fmt.Errorf("valexa: 插槽为nil")
	}
	return
}

//读取出所有Solt对象，非空值的
//	slots	所有跟踪名称和值
//	err		错误，如果不是意图类型请求，调用这个方法，将会返回错误
func (T *EchoRequest) GetValueSlots() (slots map[string]*EchoRequestRequestIntentSlot, err error) {
	gslots, err := T.GetSlots()
	if err != nil {
		return
	}
	slots = make(map[string]*EchoRequestRequestIntentSlot)
	for key, slot := range gslots {
		if slot.Value != "" {
			slots[key] = slot
		}
	}
	return slots, nil
}



//读取出所有Solt名称
//	names	所有跟踪名称
func (T *EchoRequest) GetSlotNames() (names []string) {
	names = []string{}
	slots, err := T.GetSlots()
	if err != nil {
		return
	}
	for k, _ := range slots {
		names = append(names, k)
	}
	return
}

//读取Solt值
//	slotName	跟踪名称
//	val			跟踪的值
//	err			错误，如果这个名称不存在，将会返回错误
func (T *EchoRequest) GetSlotValue(slotName string) (val string, err error) {
	slots, err := T.GetSlots()
	if err != nil {
		return
	}
	if slot, ok := slots[slotName]; ok {
		return slot.Value, nil
	}
	return "", fmt.Errorf("valexa: 名称（%s）没有找到匹配 Solt ！", slotName)
}

//读取Solt对象
//	slotName	跟踪名称
//	slot		跟踪的对象
//	err			错误，如果这个名称不存在，将会返回错误
func (T *EchoRequest) GetSlot(slotName string) (slot *EchoRequestRequestIntentSlot, err error) {
	slots, err := T.GetSlots()
	if err != nil {
		return
	}
	slot, ok := slots[slotName]
	if ok {
		return
	}
	return nil, fmt.Errorf("valexa: 名称（%s）没有找到匹配 Solt ！", slotName)
}
