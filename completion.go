package main

import (
	"regexp"
	"strings"

	"github.com/daichi-m/go-prompt"
	"github.com/kpango/glg"
)

var splitter *regexp.Regexp = regexp.MustCompile("\\s+")

/*
CmdSuggestions returns the list of suggestions based on the current input.
In case of multi word commands like ACL and LATENCY, it splits the input into two parts
and tries to filter based on the current word that is being input.
E.g., if user enters `AC` suggestions comes up as `ACL LOAD`, `ACL LOG`, `ACL SAVE`, `ACL LIST` etc.
but once user enters `ACL L` suggestions changes to `LOAD`, `LOG` and `LIST` only.
*/
func (r *Redical) CmdSuggestions(cmds []prompt.Suggest, d prompt.Document) []prompt.Suggest {
	full := strings.ToUpper(d.TextBeforeCursor())
	curr := strings.ToUpper(d.GetWordBeforeCursor())
	parts := splitter.Split(full, -1)
	full = strings.Join(parts, " ")

	spaced := len(parts) > 1

	if len(full) == 0 {
		// logger.Debug("Input empty returning empty suggest")
		return []prompt.Suggest{}
	}

	filt := prompt.FilterHasPrefix(cmds, full, true)
	if !spaced {
		// logger.Debug("Full Input: %s, Current Word: %s, IsSpaced: %t, Filtered Suggestions: %#v\n",
		// full, curr, spaced, LogSafeSlice(filt))
		return filterComplete(filt, full)
	}
	modFilt := filterMultiWord(filt, full, curr, parts)
	return filterComplete(modFilt, parts[len(parts)-1])
}

func filterMultiWord(filt []prompt.Suggest, full, current string, parts []string) []prompt.Suggest {
	spaced := len(parts) > 1
	var prefix string

	if spaced {
		prefix = parts[0]
	}
	modFilt := make([]prompt.Suggest, 0, len(filt))
	for _, x := range filt {
		sugg := prompt.Suggest{
			Text:        strings.TrimSpace(strings.TrimPrefix(x.Text, prefix)),
			Description: x.Description,
		}
		modFilt = append(modFilt, sugg)
	}
	modFilt = prompt.FilterHasPrefix(modFilt, current, true)
	glg.Debug("Full Input: %s, Current Word: %s, Prefix: %s, IsSpaced: %t, Filtered Suggestions: %#v\n",
		full, current, prefix, spaced, modFilt)
	return modFilt
}

func filterComplete(filt []prompt.Suggest, txt string) []prompt.Suggest {
	if len(filt) != 1 {
		return filt
	}

	filt1 := filt[0]
	glg.Debug("Suggestions: %#v, Full Text: %s\n", filt, txt)
	if filt1.Text == txt {
		return []prompt.Suggest{}
	}
	return filt
}
