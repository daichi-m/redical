package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Command struct {
	Arguments []Argument `json:"arguments"`
	Name      string     `json:"name"`
}
type Argument struct {
	Command  string   `json:"command,omitempty"`
	Name     []string `json:"name,omitempty"`
	Optional bool     `json:"optional,omitempty"`
	Enum     []string `json:"enum,omitempty"`
	Multiple bool     `json:"multiple,omitempty"`
}

type CommandList struct {
	Cmds []Command `json:"redisCommands"`
}

func main() {
	f, err := os.Open("redis-commands-golang-mini.json")
	if err != nil {
		panic(err)
	}
	var rc CommandList
	bts, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(bts, &rc)
	if err != nil {
		panic(err)
	}
	for _, r := range rc.Cmds {
		fmt.Printf("%+v \n", r)
	}
}
