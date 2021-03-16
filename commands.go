package main

import (
	"encoding/json"
	"fmt"
	"strings"

	_ "embed"

	"github.com/daichi-m/go-prompt"
	"go.uber.org/zap"
)

//go:embed resources/commands-golangcompat.json
var data []byte

// TODO: Better documentation comments

// Argument are the args that a redis command can take.
type Argument struct {
	Command  string   `json:"command,omitempty"`
	Name     []string `json:"name,omitempty"`
	Enum     []string `json:"enum,omitempty"`
	Optional bool     `json:"optional,omitempty"`
	Multiple bool     `json:"multiple,omitempty"`
}

// String gives the string representation of the args.
func (arg *Argument) String() string {
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

// Command encapsulates a single redis command instance.
type Command struct {
	Arguments []Argument `json:"arguments"`
	Name      string     `json:"name"`
	Return    string     `json:"return"`
	Summary   string     `json:"summary"`
	suggest   prompt.Suggest
}

// String gives the string representation of the redis command.
func (c *Command) String() string {
	var sb strings.Builder
	sb.WriteString(c.Name)
	sb.WriteString(" ")
	for _, a := range c.Arguments {
		sb.WriteString(a.String())
		sb.WriteString(" ")
	}
	if len(c.Return) > 0 {
		sb.WriteString(fmt.Sprintf(" -> %25s", c.Return))
	}
	if len(c.Summary) > 0 {
		sb.WriteString(fmt.Sprintf("\t @%s", c.Summary))
	}
	return sb.String()
}

// InitSuggest returns a promt.InitSuggest object for a RedisCommand.
func (c *Command) InitSuggest() prompt.Suggest {
	var args strings.Builder
	for _, a := range c.Arguments {
		args.WriteString(a.String())
		args.WriteString(" ")
	}
	c.suggest = prompt.Suggest{
		Text:        c.Name,
		Description: args.String(),
	}
	return c.suggest
}

// CommandList is the list of redis commands supported by redical.
type CommandList struct {
	Commands []Command `json:"redisCommands"`
	kwCmd    map[string]*Command
	multikey map[string]bool
}

// Keywords gets the list of keywords for this command list. This is an O(n) call,
// so client should be careful about calling this multiple times.
func (cl *CommandList) Keywords() []string {
	kw := make([]string, 0, len(cl.kwCmd))
	for x := range cl.kwCmd {
		xs := strings.Fields(x)
		kw = append(kw, xs...)
	}
	return kw
}

// extractCommand extracts the command and the params from a given line string.
func (cl *CommandList) extractCommand(line string) (cmd string, params []string, err error) {
	parts := strings.Fields(line)
	l := len(parts)
	err = nil

	if cl.multikey[parts[0]] {
		if l <= 1 {
			err = fmt.Errorf("Cannot split %s into command and params", line)
			return
		}
		cmd = parts[0] + " " + parts[1]
		if l < 3 {
			params = []string{}
		} else {
			params = parts[2:]
		}
		return
	}
	cmd = parts[0]
	if l < 2 {
		params = []string{}
	} else {
		params = parts[1:]
	}
	return
}

// InitCmds initializes the list of redis commands supported by redical.
func InitCmds() (*CommandList, error) {
	var cmds CommandList
	if err := json.Unmarshal(data, &(cmds.Commands)); err != nil {
		return nil, err
	}

	cmds.kwCmd = make(map[string]*Command, len(cmds.Commands))
	cmds.multikey = make(map[string]bool, len(cmds.Commands))
	for i := range cmds.Commands {
		zap.S().Debugf("Command: %v", cmds.Commands[i])
		cmds.Commands[i].InitSuggest()
		cmds.kwCmd[cmds.Commands[i].Name] = &(cmds.Commands[i])
		if strings.Contains(cmds.Commands[i].Name, " ") {
			first := strings.Fields(cmds.Commands[i].Name)[0]
			cmds.multikey[first] = true
		}
	}

	for _, c := range cmds.Commands {
		zap.S().Debugf("Command: %v", c)
	}
	return &cmds, nil
}
