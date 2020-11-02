package main

import (
	"regexp"
	"strings"

	"github.com/c-bata/go-prompt"
)

var splitter *regexp.Regexp = regexp.MustCompile("\\s+")

func completer(d prompt.Document) []prompt.Suggest {
	cmds := supported.completions
	full := strings.ToUpper(d.TextBeforeCursor())
	curr := strings.ToUpper(d.GetWordBeforeCursor())
	parts := splitter.Split(full, -1)
	full = strings.Join(parts, " ")

	spaced := len(parts) > 1

	if len(full) == 0 {
		logger.Println("Input empty returning empty suggest")
		return []prompt.Suggest{}
	}

	filt := prompt.FilterHasPrefix(cmds, full, true)
	if !spaced {
		logger.Printf("Full Input: %s, Current Word: %s, IsSpaced: %t, Filtered Suggestions: %#v\n",
			full, curr, spaced, LogSafeSlice(filt))
		return isComplete(filt, full)
	}
	modFilt := filterMulti(filt, full, curr, parts)
	return isComplete(modFilt, parts[len(parts)-1])
}

func filterMulti(filt []prompt.Suggest, full, current string, parts []string) []prompt.Suggest {
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
	logger.Printf("Full Input: %s, Current Word: %s, Prefix: %s, IsSpaced: %t, Filtered Suggestions: %#v\n",
		full, current, prefix, spaced, LogSafeSlice(modFilt))
	return modFilt
}

func isComplete(filt []prompt.Suggest, txt string) []prompt.Suggest {
	if len(filt) != 1 {
		return filt
	}

	filt1 := filt[0]
	logger.Printf("Suggestions: %#v, Full Text: %s\n", filt, txt)
	if filt1.Text == txt {
		return []prompt.Suggest{}
	}
	return filt
}
