package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type RedisCommand struct {
	Arguments  []Arguments `json:"arguments"`
	Complexity string      `json:"complexity"`
	Group      string      `json:"group"`
	Name       string      `json:"name"`
	Since      string      `json:"since"`
	Summary    string      `json:"summary"`
}
type Arguments struct {
	Command  string   `json:"command,omitempty"`
	Name     []string `json:"name,omitempty"`
	Optional bool     `json:"optional,omitempty"`
	Type     []string `json:"type"`
	Enum     []string `json:"enum,omitempty"`
}

type RedisCommands struct {
	RedisCmdList []RedisCommand `json:"redisCommands"`
}

func main() {
	f, err := os.Open("redis-commands-golang.json")
	if err != nil {
		panic(err)
	}
	var rc RedisCommands
	bts, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(bts, &rc)
	if err != nil {
		panic(err)
	}
	for _, r := range rc.RedisCmdList {
		fmt.Printf("%+v \n", r)
	}
}
