package main

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/daichi-m/prompter/assets"
)

// TODO: Better documentation comments

// Argument are the args that a redis command can take
type Argument struct {
	Command  string   `json:"command,omitempty"`
	Name     []string `json:"name,omitempty"`
	Optional bool     `json:"optional,omitempty"`
	Enum     []string `json:"enum,omitempty"`
	Multiple bool     `json:"multiple,omitempty"`
}

func (arg Argument) String() string {
	var sb strings.Builder
	if len(arg.Command) != 0 {
		sb.WriteString(strings.ToUpper(arg.Command))
		sb.WriteString(" ")
	}
	if len(arg.Name) != 0 {
		sb.WriteString(strings.Join(arg.Name, " "))
		if arg.Multiple {
			sb.WriteString("...")
		}
	} else if len(arg.Enum) != 0 {
		sb.WriteString(strings.Join(arg.Enum, "|"))
	}

	var s string
	if arg.Optional {
		s = "[ " + sb.String() + " ]"
	} else {
		s = sb.String()
	}
	return s
}

// Command encapsulates a single redis command instance
type Command struct {
	Arguments []Argument `json:"arguments"`
	Name      string     `json:"name"`
}

func (c Command) String() string {
	var sb strings.Builder
	sb.WriteString(c.Name)
	sb.WriteString(" ")
	for _, a := range c.Arguments {
		sb.WriteString(a.String())
		sb.WriteString(" ")
	}
	return sb.String()
}

// Suggest returns a promt.Suggest object for a RedisCommand
func (c Command) Suggest() prompt.Suggest {
	var args strings.Builder
	for _, a := range c.Arguments {
		args.WriteString(a.String())
		args.WriteString(" ")
	}
	return prompt.Suggest{
		Text:        c.Name,
		Description: args.String(),
	}
}

// CommandList is the list of redis commands supported by redical
type CommandList struct {
	Cmds        []Command `json:"redisCommands"`
	completions []prompt.Suggest
}

// InitSuggests initializes the prompt.Suggest slice for the redis command list
func (cl *CommandList) InitSuggests() error {
	if len(cl.completions) > 0 {
		return errors.New("Already initialized the suggestions")
	}

	compl := make([]prompt.Suggest, 0, len(cl.Cmds))
	for _, c := range cl.Cmds {
		sg := c.Suggest()
		compl = append(compl, sg)
	}
	cl.completions = compl
	return nil
}

// supported is the list of supported redis commands
var supported CommandList

// InitCmds initializes the list of redis commands supported by redical
func InitCmds() error {
	data, err := assets.Asset("resources/redis-commands-golang.json")
	if err != nil {
		return err
	}

	if err = json.Unmarshal(data, &supported); err != nil {
		return err
	}
	if err = supported.InitSuggests(); err != nil {
		return err
	}

	return nil
}
