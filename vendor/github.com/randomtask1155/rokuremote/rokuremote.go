package rokuremote

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Player struct {
	Address  string
	Port     string
	NickName string
	Client   *http.Client
}

// http://yourrokuaddress:8060/query/apps
const (
	Hulu        = 2285
	Netflix     = 12
	HBO         = 61322
	Plex        = 13535
	AmazonPrime = 13
	USTVNow     = 2026
	PBS         = 23333
	DisneyPlus  = 291097
	NickJR      = 66595
	YouTube     = 837
	CBS         = 27536
)

func Connect(ip string) (Player, error) {
	return Player{ip, "8060", "", &http.Client{}}, nil
}

// Given player nickname connect and return player
func ConnectName(ip, name string) (Player, error) {

	return Player{ip, "8060", name, &http.Client{}}, nil
}

// Get send get requst to rokue and return response
func (p *Player) Get(path string) ([]byte, error) {
	req, err := http.NewRequest("GET", "http://"+p.Address+":8060"+path, new(bytes.Buffer))
	if err != nil {
		return nil, err
	}
	resp, err := p.Client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("ResponseError:%s: %s\n", path, resp.Status)
	}
	return ioutil.ReadAll(resp.Body)
}

func (p *Player) Post(path string) error {
	req, err := http.NewRequest("POST", "http://"+p.Address+":8060"+path, new(bytes.Buffer))
	if err != nil {
		return err
	}
	resp, err := p.Client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("ResponseError:%s: %s\n", path, resp.Status)
	}

	return nil
}

func (p *Player) Home() error {
	return p.Post("/keypress/home")
}
func (p *Player) Up() error {
	return p.Post("/keypress/up")
}
func (p *Player) Down() error {
	return p.Post("/keypress/down")
}
func (p *Player) Left() error {
	return p.Post("/keypress/left")
}
func (p *Player) Right() error {
	return p.Post("/keypress/right")
}
func (p *Player) Select() error {
	return p.Post("/keypress/select")
}
func (p *Player) Enter() error {
	return p.Post("/keypress/enter")
}
func (p *Player) Back() error {
	return p.Post("/keypress/back")
}
func (p *Player) Search() error {
	return p.Post("/keypress/search")
}
func (p *Player) Play() error {
	return p.Post("/keypress/play")
}

func (p *Player) StartChannel(chanid int) error {
	return p.Post("/launch/" + strconv.Itoa(chanid))
}
