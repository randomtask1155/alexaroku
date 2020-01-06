package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"net/http"

	alexa "github.com/mikeflynn/go-alexa/skillserver"
	roku "github.com/randomtask1155/rokuremote"
)

// Applications defined apps
var Applications map[string]interface{}

var (
	rokuPlayer    roku.Player
	appID         string
	rokuIPAddress string
	rc *RemoteSession
)

func main() {
	var err error

	rc = NewRemoteSession()
	go rc.cleaner()

	Applications = map[string]interface{}{
		"/echo/roku": alexa.EchoApplication{ // Route
			AppID:    os.Getenv("ALEXAAPPID"), // Echo App ID from Amazon Dashboard
			//OnIntent: EchoIntentHandler,
			//OnLaunch: EchoIntentHandler,
			Handler: EchoIntentHandler,
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
	success := true
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
		rokuPlayer.Select() // select is like the ok button
	case "back":
		rokuPlayer.Back()
	case "search":
		rokuPlayer.Search()
	case "pause":
		rokuPlayer.Play() // play pauses too
	case "play":
		rokuPlayer.Play()
	default:
		echoResp.OutputSpeech("that key does not exist")
		success = false
	}
	if success {
		echoResp.OutputSpeech(fmt.Sprintf("pressing %s", key))
	}
}

// SelectChannel opens the given channel
func SelectChannel(echoReq *alexa.EchoRequest, echoResp *alexa.EchoResponse) {
	channel, err := echoReq.GetSlotValue("Channel")
	if err != nil {
		log.Println(err)
		echoResp.OutputSpeech(fmt.Sprintf("i could not find channel %s", channel))
	}
	success := true
	channel = strings.ToLower(channel)
	switch channel {
	case "netflix":
		rokuPlayer.StartChannel(roku.Netflix)
	case "amazon":
		rokuPlayer.StartChannel(roku.AmazonPrime)
	case "hulu":
		rokuPlayer.StartChannel(roku.Hulu)
	case "movies":
		rokuPlayer.StartChannel(roku.HBO)
	case "hbo":
		rokuPlayer.StartChannel(roku.HBO)
	case "h b o":
		rokuPlayer.StartChannel(roku.HBO)
	case "plex":
		rokuPlayer.StartChannel(roku.Plex)
	case "television":
		rokuPlayer.StartChannel(roku.USTVNow)
	case "p b s":
		rokuPlayer.StartChannel(roku.PBS)
	case "pbs":
		rokuPlayer.StartChannel(roku.PBS)
	case "nick":
		rokuPlayer.StartChannel(roku.NickJR)
	case "disney":
		rokuPlayer.StartChannel(roku.DisneyPlus)
	case "youtube":
		rokuPlayer.StartChannel(roku.YouTube)
	case "you tube":
		rokuPlayer.StartChannel(roku.YouTube)
	case "cbs":
		rokuPlayer.StartChannel(roku.CBS)
	case "c b s":
		rokuPlayer.StartChannel(roku.CBS)
	default:
		echoResp.OutputSpeech("I do not know that channel")
		fmt.Printf("Unknown channel: %s\n", channel)
		success = false
	}

	if success {
		echoResp.OutputSpeech(fmt.Sprintf("starting channel %s", channel))
	}
}

// EchoIntentHandler determine intent
//func EchoIntentHandler(echoReq *alexa.EchoRequest, echoResp *alexa.EchoResponse) {
func EchoIntentHandler(w http.ResponseWriter, r *http.Request) {
	echoReq := alexa.GetEchoRequest(r)
	echoResp := alexa.NewEchoResponse().OutputSpeech("I am sorry but roku does not understand your request").EndSession(true)
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
		case "StartRemote":
			echoResp = rc.startController(echoReq, echoReq.GetSessionID())
		case "RemoteControl":
			echoResp = rc.ExecuteCommand(echoReq, echoReq.GetSessionID())
		case "AMAZON.NavigateHomeIntent":
			rokuPlayer.Home() // why? because amazon is stupid
			echoResp.OutputSpeech(fmt.Sprintf("pressing home"))
		default:
			fmt.Printf("Invalid Intent: %s\n", echoReq.GetIntentName())
			echoResp.OutputSpeech("Sorry you must have bad intentions and refuse your request").Card("Failure", "Invalid Intent")
		}
	} else if echoReq.GetRequestType() == "LaunchRequest" {
		echoResp = rc.startController(echoReq, echoReq.GetSessionID())
		log.Println("received launch request")
	} else if echoReq.GetRequestType() == "SessionEndedRequest" {
		rc.ExpireSession(echoReq.GetSessionID())
		echoResp.OutputSpeech("stopping remote control")
	} else {
		log.Println(echoReq.GetRequestType())
		fmt.Printf("%v\n", echoReq)
		echoResp.OutputSpeech("I am sorry but roku does not understand your request").Card("Failure", "Invalid request")
	}
	json, _ := echoResp.String()
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.Write(json)
}
