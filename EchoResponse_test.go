package valexa
	
import (
	"testing"
)

func Test_EchoResponse_SetShouldEndSession(testingT *testing.T){
	er := &EchoResponse{}
	er.SetEndSession(true)
	if !er.Response.ShouldEndSession {
		testingT.Fatal("err")
	}
}

func Test_EchoResponse_OutputSpeech(testingT *testing.T){
	er := &EchoResponse{}
	er.OutputSpeech("1234")
	if er.Response.OutputSpeech.Text != "1234" {
		testingT.Fatal("err")
	}
}

func Test_EchoResponse_OutputSpeechSSML(testingT *testing.T){
	er := &EchoResponse{}
	er.OutputSpeechSSML("1234")
	if er.Response.OutputSpeech.Ssml != "1234" {
		testingT.Fatal("err")
	}
}

func Test_EchoResponse_SimpleCard(testingT *testing.T){
	er := &EchoResponse{}
	er.SimpleCard("title","1234")
	if er.Response.Card.Title != "title" || er.Response.Card.Content != "1234" {
		testingT.Fatal("err")
	}
}

func Test_EchoResponse_StandardCard(testingT *testing.T){
	er := &EchoResponse{}
	er.StandardCard("title","1234", "https://y.x.z/1.jpg", "https://y.x.z/2.jpg")
	if er.Response.Card.Title != "title" || 
		er.Response.Card.Text != "1234" || 
			er.Response.Card.Image.SmallImageURL != "https://y.x.z/1.jpg" || 
			er.Response.Card.Image.LargeImageURL != "https://y.x.z/2.jpg"{
		testingT.Fatal("err")
	}
}

func Test_EchoResponse_LinkAccountCard(testingT *testing.T){
	er := &EchoResponse{}
	er.LinkAccountCard()
	if er.Response.Card.Type != "LinkAccount" {
		testingT.Fatal("err")
	}
}

func Test_EchoResponse_RepromptText(testingT *testing.T){
	er := &EchoResponse{}
	er.RepromptText("1234")
	if er.Response.Reprompt.OutputSpeech.Text != "1234" {
		testingT.Fatal("err")
	}
}

func Test_EchoResponse_RepromptSSML(testingT *testing.T){
	er := &EchoResponse{}
	er.RepromptSSML("1234")
	if er.Response.Reprompt.OutputSpeech.Ssml != "1234" {
		testingT.Fatal("err")
	}
}

func Test_EchoResponse_String(testingT *testing.T){
	er := &EchoResponse{}
	testingT.Log(er.String())
}



