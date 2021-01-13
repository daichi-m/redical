package main

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/daichi-m/go-prompt"
	"github.com/daichi-m/redical/assets"
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
	Commands    []Command `json:"redisCommands"`
	completions []prompt.Suggest
	keywords    []string
}

// InitSuggests initializes the prompt.Suggest slice for the redis command list
func (cl *CommandList) InitSuggests() error {
	if len(cl.completions) > 0 {
		return errors.New("Already initialized the suggestions")
	}

	compl := make([]prompt.Suggest, 0, len(cl.Commands))
	for _, c := range cl.Commands {
		sg := c.Suggest()
		compl = append(compl, sg)
	}
	cl.completions = compl
	return nil
}

// InitCmds initializes the list of redis commands supported by redical
func InitCmds() (*CommandList, error) {
	data, err := assets.Asset("resources/redis-commands-golang.json")
	if err != nil {
		return nil, err
	}
	var cmds CommandList
	if err = json.Unmarshal(data, &cmds); err != nil {
		return nil, err
	}
	if err = cmds.InitSuggests(); err != nil {
		return nil, err
	}
	for _, cmd := range cmds.Commands {
		p := strings.Fields(cmd.Name)
		cmds.keywords = append(cmds.keywords, p...)
	}
	return &cmds, nil
}
