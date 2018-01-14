package main

import (
	alexa "github.com/mikeflynn/go-alexa/skillserver"
	"os"
	roku "github.com/randomtask1155/rokuremote"
	"fmt"
	"log"
)

var Applications map[string]interface{}

var (
	rokuPlayer roku.Player
	appID string
	rokuIPAddress string
)

func main() {
	var err error

	Applications = map[string]interface{}{
		"/echo/helloworld": alexa.EchoApplication{ // Route
			AppID:   os.Getenv("ALEXAAPPID") , // Echo App ID from Amazon Dashboard
			OnIntent: EchoIntentHandler,
			OnLaunch: EchoIntentHandler,
		},
	}
	rokuIPAddress = os.Getenv("ROKUIP")

	rokuPlayer, err = roku.Connect(rokuIPAddress)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	alexa.Run(Applications, os.Getenv("PORT"))
}

func PerformKeyPress(echoReq *alexa.EchoRequest, echoResp *alexa.EchoResponse) {
	key, err := echoReq.GetSlotValue("Command")
	if err != nil {
		log.Println(err)
		echoResp.OutputSpeech("i do not understand the command")
	}
	log.Printf("Kepressed: %s\n", key)
	switch key {
	case "home":
		rokuPlayer.Home()
	case "up":
		rokuPlayer.Up()
	case "down":
		rokuPlayer.Down()
	case "left":
		rokuPlayer.Left()
	case "right":
		rokuPlayer.Right()
	case "enter":
		rokuPlayer.OK()
	case "back":
		rokuPlayer.Back()
	default:
		echoResp.OutputSpeech("that key does not exist")
	}
}

func SelectChannel(echoReq *alexa.EchoRequest, echoResp *alexa.EchoResponse) {
	channel, err := echoReq.GetSlotValue("Channel")
	if err != nil {
		log.Println(err)
		echoResp.OutputSpeech(fmt.Sprintf("i could not find channel %s", channel))
	}
	switch channel {
	case "netflix":
		rokuPlayer.StartChannel(roku.Netflix)
	case "amazon":
		rokuPlayer.StartChannel(roku.AmazonPrime)
	case "hulu":
		rokuPlayer.StartChannel(roku.Hulu)
	case "movies":
		rokuPlayer.StartChannel(roku.HBO)
	case "plex":
		rokuPlayer.StartChannel(roku.Plex)
	case "television":
		rokuPlayer.StartChannel(roku.USTVNow)
	default:
		echoResp.OutputSpeech("I do not know that channel")
	}
}

func EchoIntentHandler(echoReq *alexa.EchoRequest, echoResp *alexa.EchoResponse) {
	//err := rokuPlayer.Home()
	//if err != nil {
	//	fmt.Println(err)
	//}
	//echoResp.OutputSpeech("Hello world from my new Echo test app!").Card("Hello World", "This is a test card.")

	if echoReq.GetRequestType() == "IntentRequest" {
		log.Println(echoReq.GetIntentName())

		switch echoReq.GetIntentName() {
		case "KeyPress":
			PerformKeyPress(echoReq, echoResp)
		case "PickChannel":
			SelectChannel(echoReq, echoResp)
		default:
			echoResp.OutputSpeech("Sorry you must have bad intentinos and refuse your request").Card("Failure", "Invalid Intent")
		}
	} else {
		echoResp.OutputSpeech("I am sorry but roku does not understand your request").Card("Failure", "Invalid request")
	}

}