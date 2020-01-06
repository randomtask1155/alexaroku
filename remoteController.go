package main 

import (
	"time"
	"sync"
	"log"

	alexa "github.com/mikeflynn/go-alexa/skillserver"
)

var (
	sessionLock sync.Mutex
)

type Session struct {
	ID string 
	Expire int64 // set a max duration in so system can clean itself up
}


type RemoteSession struct {
	Sessions []Session 
}

func NewRemoteSession() *RemoteSession {
	return &RemoteSession{make([]Session,0)}
}

func (r *RemoteSession) cleaner() {
	for {
		sessionLock.Lock()
		newSessions := make([]Session,0)
		for i := range r.Sessions{
			if r.Sessions[i].Expire < time.Now().Unix(){
				log.Printf("destroying expired session %s\n", r.Sessions[i].ID)
				continue
			} 
			newSessions = append(newSessions, r.Sessions[i])
		}
		r.Sessions = newSessions
		sessionLock.Unlock()
		time.Sleep(60 * time.Second)
	}
}

func (r *RemoteSession) ExpireSession(s string) {
	sessionLock.Lock()
	defer sessionLock.Unlock()
	for i := range r.Sessions{
		if r.Sessions[i].ID == s {
			log.Printf("expiring session %s]n", s)
			r.Sessions[i].Expire = time.Now().Unix()
		}
	}
}

func (r *RemoteSession) setSession(s string) {
	if r.checkSession(s) {
		return // already exists
	}
	sessionLock.Lock()
	defer sessionLock.Unlock()
	newExp := time.Now().Unix() + int64(180 * time.Second) // expire after 3 minutes
	log.Printf("creating new session %s", s)
	r.Sessions = append(r.Sessions, Session{s, newExp})
}

func (r *RemoteSession) checkSession(s string) bool {
	check := r.getSession(s)
	if check == "" {
		return false
	}
	return true
}

func (r *RemoteSession) getSession(s string) string {
	sessionLock.Lock()
	defer sessionLock.Unlock()
	newExp := time.Now().Unix() + int64(180 * time.Second) // expire after 3 minutes
	for i := range r.Sessions {
		if r.Sessions[i].ID == s {
			log.Printf("updating session %s", r.Sessions[i].ID)
			r.Sessions[i].Expire = newExp
			return s
		}
	}
	return ""
}

func (r *RemoteSession) startController(echoReq *alexa.EchoRequest, s string) *alexa.EchoResponse {
	if echoReq.GetIntentName() == "StartRemote" {
		rc.setSession(s)
		log.Printf("starting controller with session %s\n", s)
		return alexa.NewEchoResponse().OutputSpeech("activating remote control").EndSession(false)
	}
	rc.ExpireSession(s)
	return alexa.NewEchoResponse().OutputSpeech("could not activate remote control").EndSession(true)
}

func (r *RemoteSession) ExecuteCommand(echoReq *alexa.EchoRequest, s string) *alexa.EchoResponse{
	if ! r.checkSession(s) {
		return alexa.NewEchoResponse().OutputSpeech("sorry please activate remote").EndSession(true)
	}
	sessionLock.Lock()
	defer sessionLock.Unlock()
	// ensure session is updated
	command, err := echoReq.GetSlotValue("Command")
	if err != nil {
		log.Println(err)
		return alexa.NewEchoResponse().OutputSpeech("sorry try again").EndSession(false)
	}
	log.Printf("Kepressed during session: %s\n", command)
	switch command {
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
	case "close":
		go rc.ExpireSession(s)  // need to put this in go routine or deadlock
		return alexa.NewEchoResponse().OutputSpeech("stopping remote control").EndSession(true)
	case "stop":
		go rc.ExpireSession(s)  // need to put this in go routine or deadlock
		return alexa.NewEchoResponse().OutputSpeech("stopping remote control").EndSession(true)
	case "end":
		go rc.ExpireSession(s) // need to put this in go routine or deadlock
		return alexa.NewEchoResponse().OutputSpeech("stopping remote control").EndSession(true)
	default:
		return alexa.NewEchoResponse().OutputSpeech("that key does not exist").EndSession(false)
	}

	return alexa.NewEchoResponse().OutputSpeech("ok").EndSession(false)

}