package valexa

import (
	"encoding/json"
)


/*
https://developer.amazon.com/zh/docs/custom-skills/request-and-response-json-reference.html

HTTP/1.1 200 OK
Content-Type: application/json;charset=UTF-8
Content-Length:

{
  "version": "string",
  "sessionAttributes": {
    "key": "value"
  },
  "response": {
    "outputSpeech": {
      "type": "PlainText",
      "text": "Plain text string to speak",
      "ssml": "<speak>SSML text string to speak</speak>",
      "playBehavior": "REPLACE_ENQUEUED"      
    },
    "card": {
      "type": "Standard",
      "title": "Title of the card",
      "content": "Content of a simple card",
      "text": "Text content for a standard card",
      "image": {
        "smallImageUrl": "https://url-to-small-card-image...",
        "largeImageUrl": "https://url-to-large-card-image..."
      }
    },
    "reprompt": {
      "outputSpeech": {
        "type": "PlainText",
        "text": "Plain text string to speak",
        "ssml": "<speak>SSML text string to speak</speak>",
        "playBehavior": "REPLACE_ENQUEUED"             
      }
    },
    "directives": [
      {
        "type": "InterfaceName.Directive"
        (...properties depend on the directive type)
      }
    ],
    "shouldEndSession": true
  }
}

Response Examples：
{
  "version": "1.0",
  "sessionAttributes": {
    "supportedHoroscopePeriods": {
      "daily": true,
      "weekly": false,
      "monthly": false
    }
  },
  "response": {
    "outputSpeech": {
      "type": "PlainText",
      "text": "Today will provide you a new learning opportunity.  Stick with it and the possibilities will be endless. Can I help you with anything else?"
    },
    "card": {
      "type": "Simple",
      "title": "Horoscope",
      "content": "Today will provide you a new learning opportunity.  Stick with it and the possibilities will be endless."
    },
    "reprompt": {
      "outputSpeech": {
        "type": "PlainText",
        "text": "Can I help you with anything else?"
      }
    },
    "shouldEndSession": false
  }
}

Response to IntentRequest or Launch Request with Directives Example:
{
  "version": "1.0",
  "sessionAttributes": {},
  "response": {
    "outputSpeech": {
      "type": "PlainText",
      "text": "Playing the requested song."
    },
    "card": {
      "type": "Simple",
      "title": "Play Audio",
      "content": "Playing the requested song."
    },
    "reprompt": {
      "outputSpeech": {
        "type": "PlainText",
        "text": null
      }
    },
    "directives": [
      {
        "type": "AudioPlayer.Play",
        "playBehavior": "ENQUEUE",
        "audioItem": {
          "stream": {
            "token": "this-is-the-audio-token",
            "url": "https://my-audio-hosting-site.com/audio/sample-song.mp3",
            "offsetInMilliseconds": 0
          }
        }
      }
    ],
    "shouldEndSession": true
  }
}

Response to AudioPlayer or PlaybackController Example (Directives Only):
{
  "version": "1.0",
  "response": {
    "directives": [
      {
        "type": "AudioPlayer.Play",
        "playBehavior": "REPLACE_ALL",
        "audioItem": {
          "stream": {
            "token": "track2-long-audio",
            "url": "https://my-audio-hosting-site.com/audio/sample-song-2.mp3",
            "offsetInMilliseconds": 0
          }
        }
      }
    ]
  }
}


*/
type EchoResponseResponseCardImage struct {
	LargeImageURL string 												`json:"largeImageUrl,omitempty"`
	SmallImageURL string 												`json:"smallImageUrl,omitempty"`
}

type EchoResponseResponseCard struct {
	Content string 														`json:"content,omitempty"`
	Image   EchoResponseResponseCardImage								`json:"image,omitempty"`
	Text  	string 														`json:"text,omitempty"`
	Title 	string 														`json:"title,omitempty"`
	Type  	string 														`json:"type,omitempty"`
}

type EchoResponseResponseDirectivesAudioItemStream struct {
	ExpectedPreviousToken string 										`json:"expectedPreviousToken"`
	OffsetInMilliseconds  int    										`json:"offsetInMilliseconds"`
	Token                 string 										`json:"token"`
	URL                   string 										`json:"url"`
}

type EchoResponseResponseDirectivesAudioItem struct {
	Stream  EchoResponseResponseDirectivesAudioItemStream				`json:"stream"`
}

type EchoResponseResponseDirectives struct {
	AudioItem  		EchoResponseResponseDirectivesAudioItem				`json:"audioItem"`
	PlayBehavior 	string 												`json:"playBehavior"`
	Type         	string 												`json:"type"`
}

type EchoResponseResponseOutputSpeech struct {
	Ssml string 														`json:"ssml,omitempty"`
	Text string 														`json:"text,omitempty"`
	Type string 														`json:"type,omitempty"`
}

type EchoResponseResponseRepromptOutputSpeech struct {
	Ssml string 														`json:"ssml,omitempty"`
	Text string 														`json:"text,omitempty"`
	Type string 														`json:"type,omitempty"`
}

type EchoResponseResponseReprompt struct {
	OutputSpeech EchoResponseResponseRepromptOutputSpeech 				`json:"outputSpeech,omitempty"`
}

type EchoResponseResponse struct {
	Card  				*EchoResponseResponseCard						`json:"card,omitempty"`
	Directives 			[]EchoResponseResponseDirectives 				`json:"directives,omitempty"`
	OutputSpeech  		*EchoResponseResponseOutputSpeech 				`json:"outputSpeech,omitempty"`
	Reprompt  			*EchoResponseResponseReprompt					`json:"reprompt,omitempty"`
	ShouldEndSession 	bool 											`json:"shouldEndSession"`
}

type EchoResponse struct {
	Response  			EchoResponseResponse							`json:"response"`
	SessionAttributes 	map[string]interface{}							`json:"sessionAttributes,omitempty"`
	Version 			string 											`json:"version"`
}

//默认的响应对象
func NewEchoResponse() *EchoResponse {
	er := &EchoResponse{
		Version: "1.0",
		Response: EchoResponseResponse{
			ShouldEndSession: false,
		},
		SessionAttributes: make(map[string]interface{}),
	}

	return er
}

//设置是否结束本次会活，默认是结束会话
//	ok		true结束会话，fales不结束会话
func (T *EchoResponse) SetEndSession(ok bool) *EchoResponse {
	T.Response.ShouldEndSession = ok
	return T
}

//设置语音输出
//	text			语音文本
//	*EchoResponse	响应对象
func (T *EchoResponse) OutputSpeech(text string) *EchoResponse {
	T.Response.OutputSpeech = &EchoResponseResponseOutputSpeech{
		Type	: "PlainText",
		Text	: text,
	}
	return T
}

//设置语音输出
//更多浏览这里: https://developer.amazon.com/docs/custom-skills/speech-synthesis-markup-language-ssml-reference.html
//	text			语音文本（支持语气表达，等等）
//	*EchoResponse	响应对象
func (T *EchoResponse) OutputSpeechSSML(text string) *EchoResponse {
	T.Response.OutputSpeech = &EchoResponseResponseOutputSpeech{
		Type	: "SSML",
		Ssml	: text,
	}
	return T
}

//设置在屏幕上显示的卡片
//	title			标题
//	content			内容
//	*EchoResponse	响应对象
func (T *EchoResponse) SimpleCard(title string, content string) *EchoResponse {
	T.Response.Card = &EchoResponseResponseCard{
		Type	: "Simple",
		Title	: title,
		Content	: content,
	}
	return T
}

//设置在屏幕上显示的卡片，支持图片
//	title			标题
//	text			内容
//	smallImg		小图片 720w x 480h
//	largeImg		中图片 1200w x 800h
//	*EchoResponse	响应对象
func (T *EchoResponse) StandardCard(title string, text string, smallImg string, largeImg string) *EchoResponse {
	T.Response.Card = &EchoResponseResponseCard{
		Type	: "Standard",
		Title	: title,
		Text	: text,
	}
	if smallImg != "" {
		T.Response.Card.Image.SmallImageURL = smallImg
	}
	if largeImg != "" {
		T.Response.Card.Image.LargeImageURL = largeImg
	}
	return T
}

//设置在屏幕上显示的卡片，仅用于用户认证
//	*EchoResponse	响应对象
func (T *EchoResponse) LinkAccountCard() *EchoResponse {
	T.Response.Card = &EchoResponseResponseCard{
		Type: "LinkAccount",
	}
	return T
}

//设置回复确认输出
//	text			内容
//	*EchoResponse	响应对象
func (T *EchoResponse) RepromptText(text string) *EchoResponse {
	T.Response.Reprompt = &EchoResponseResponseReprompt{
		OutputSpeech: EchoResponseResponseRepromptOutputSpeech{
			Type: "PlainText",
			Text: text,
		},
	}
	return T
}

//设置回复确认输出
//更多浏览这里: https://developer.amazon.com/docs/custom-skills/speech-synthesis-markup-language-ssml-reference.html
//	text			语音文本（支持语气表达，等等）
//	*EchoResponse	响应对象
func (T *EchoResponse) RepromptSSML(text string) *EchoResponse {
	T.Response.Reprompt = &EchoResponseResponseReprompt{
		OutputSpeech: EchoResponseResponseRepromptOutputSpeech{
			Type: "SSML",
			Ssml: text,
		},
	}
	return T
}


func (T *EchoResponse) String() string {
	jsonStr, err := json.Marshal(T)
	if err != nil {
		return "{}"
	}
	return string(jsonStr)
}

