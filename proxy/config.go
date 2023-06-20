package proxy

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"sync"
)


type Config struct{
	Proxy Proxy `json:"proxy"`
	Nodes []Nodes `json:"nodes"`
}

func NewConfig() Config{
	return Config{}
}

type Proxy struct{
	Port string `json:"port"`
}

type Nodes struct{
	syncRw sync.RWMutex
	URL string `json:"url"`
	Up bool
}

func (n *Nodes) getUpState() bool{
	n.syncRw.Lock()
	isUp := n.Up
	n.syncRw.Unlock()
	return isUp
}

func (n *Nodes) setUpState(state bool){
	n.syncRw.Lock()
	n.Up = state
	n.syncRw.Unlock()
}

func(c *Config) GetConfig() Config{
	var config Config
	data , err := ioutil.ReadFile("./config.json")

	if err != nil{
		log.Fatal(err.Error())
	}

	json.Unmarshal(data ,&config)
	return config
	
} 


