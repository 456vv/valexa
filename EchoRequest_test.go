package valexa

import (
	"testing"
	"time"
	
)

func Test_EchoRequest_VerifyTimestamp(testingT *testing.T){
	er := &EchoRequest{Request:EchoRequestRequest{}}
	tests := []struct{
		timestamp 	string
		limit		int
		result		bool
		str			string
	}{
		//限制10s，设置了错误时间
		{timestamp:"2017-12-01T09:00:10Z", limit:10, result: false, str:"111"},
		//限制15s，设置了-14s
		{timestamp:time.Now().UTC().Add(time.Second*(-14)).Format(time.RFC3339), limit:15, result: true, str:"222"},
		//限制15s，设置了-15s
		{timestamp:time.Now().UTC().Add(time.Second*(-15)).Format(time.RFC3339), limit:15, result: false, str:"333"},
		//限制15s，设置了-16s
		{timestamp:time.Now().UTC().Add(time.Second*(-16)).Format(time.RFC3339), limit:15, result: false, str:"444"},
		//限制0s，设置了-16s
		{timestamp:time.Now().UTC().Add(time.Second*(-16)).Format(time.RFC3339), limit:0, result: false, str:"555"},
		//限制0s，设置了-1s
		{timestamp:time.Now().UTC().Add(time.Second*(-1)).Format(time.RFC3339), limit:0, result: false, str:"666"},
		//限制1s，设置了-0s
		{timestamp:time.Now().UTC().Add(time.Second*(0)).Format(time.RFC3339), limit:1, result: true, str:"777"},
	}
	for _, test := range tests {
		er.Request.Timestamp = test.timestamp
		if er.VerifyTimestamp(test.limit) != test.result {
			testingT.Fatal("err", test.str)
		}
	}
}

func Test_EchoRequest_GetApplicationID(testingT *testing.T){
	er := &EchoRequest{}
	if er.GetApplicationID() != "" {
		testingT.Fatal("err")
	}
	er.Session.Application.ApplicationID="1234"
	if er.GetApplicationID() != "1234" {
		testingT.Fatal("err")
	}
}

func Test_EchoRequest_GetSessionID(testingT *testing.T){
	er := &EchoRequest{}
	if er.GetSessionID() != "" {
		testingT.Fatal("err")
	}
	er.Session.SessionID="1234"
	if er.GetSessionID() != "1234" {
		testingT.Fatal("err")
	}
}

func Test_EchoRequest_GetUserID(testingT *testing.T){
	er := &EchoRequest{}
	if er.GetUserID() != "" {
		testingT.Fatal("err")
	}
	er.Session.User.UserID="1234"
	if er.GetUserID() != "1234" {
		testingT.Fatal("err")
	}
}

func Test_EchoRequest_GetAccessToken(testingT *testing.T){
	er := &EchoRequest{}
	if er.GetAccessToken() != "" {
		testingT.Fatal("err")
	}
	er.Session.User.AccessToken="1234"
	if er.GetAccessToken() != "1234" {
		testingT.Fatal("err")
	}
}



func Test_EchoRequest_GetRequestType(testingT *testing.T){
	er := &EchoRequest{}
	if er.GetRequestType() != "" {
		testingT.Fatal("err")
	}
	er.Request.Type="1234"
	if er.GetRequestType() != "1234" {
		testingT.Fatal("err")
	}
}

func Test_EchoRequest_GetIntentName(testingT *testing.T){
	er := &EchoRequest{}
	if name := er.GetIntentName(); name != "" {
		testingT.Fatal("err")
	}
	
	er.Request.Type="IntentRequest"
	er.Request.Intent.Name="Iiem"
	
	if intent := er.GetIntentName(); intent != "Iiem" {
		testingT.Fatal("err")
	}
}

func Test_EchoRequest_GetSlots(testingT *testing.T){
	er := &EchoRequest{}
	er.Request.Type="IntentRequest"
	if _, err := er.GetSlots(); err == nil {
		testingT.Fatal("err")
	}
	
	er.Request.Intent.Slots = make(map[string]*EchoRequestRequestIntentSlot)
	if slots, err := er.GetSlots(); err != nil || slots == nil {
		testingT.Fatal("err")
	}
	
	er.Request.Intent.Slots = make(map[string]*EchoRequestRequestIntentSlot)
	er.Request.Intent.Slots["Item"] = &EchoRequestRequestIntentSlot{
		Name	: "A",
		Value	: "B",
	}
	slots, err := er.GetSlots()
	if err != nil || slots == nil {
		testingT.Fatal("err")
	}

	if	slot, ok := slots["Item"]; !ok || slot.Name != "A" || slot.Value != "B" {
		testingT.Fatal("err")
	}
}

func Test_EchoRequest_GetSlotNames(testingT *testing.T){
	er := &EchoRequest{}
	names := er.GetSlotNames()
	if len(names) > 0 {
		testingT.Fatal("err")
	}
	er.Request.Intent.Slots = make(map[string]*EchoRequestRequestIntentSlot)
	names = er.GetSlotNames()
	if len(names) > 0 {
		testingT.Fatal("err")
	}
	er.Request.Type="IntentRequest"
	er.Request.Intent.Slots["Item"] = &EchoRequestRequestIntentSlot{
		Name	: "A",
		Value	: "B",
	}
	names = er.GetSlotNames()
	if len(names) != 1 || names[0] != "Item" {
		testingT.Fatal("err")
	}
}

func Test_EchoRequest_GetSlotValue(t *testing.T){
	er := &EchoRequest{}
	val, err := er.GetSlotValue("Item")
	if err == nil {
		t.Fatal("err")
	}
	er.Request.Type="IntentRequest"
	er.Request.Intent.Slots = make(map[string]*EchoRequestRequestIntentSlot)
	er.Request.Intent.Slots["Item"] = &EchoRequestRequestIntentSlot{
		Name	: "A",
		Value	: "B",
	}
	val, err = er.GetSlotValue("Item")
	if err != nil || val != "B"  {
		t.Fatal("err")
	}
}

func Test_EchoRequest_GetSlot_1(testingT *testing.T){
	er := &EchoRequest{}
	slot, err := er.GetSlot("Item")
	if err == nil {
		testingT.Fatal("err")
	}
	er.Request.Intent.Slots = nil
	slot, err = er.GetSlot("Item")
	if err == nil {
		testingT.Fatal("err")
	}
	er.Request.Type="IntentRequest"
	er.Request.Intent.Slots = make(map[string]*EchoRequestRequestIntentSlot)
	er.Request.Intent.Slots["Item"] = &EchoRequestRequestIntentSlot{
		Name	: "A",
		Value	: "B",
	}
	slot, err = er.GetSlot("Item")
	if err != nil || slot.Name != "A" || slot.Value != "B"{
		testingT.Fatal("err")
	}
}