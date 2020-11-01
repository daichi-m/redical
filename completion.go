package main

import (
	"strings"

	"github.com/c-bata/go-prompt"
)

func completer(d prompt.Document) []prompt.Suggest {
	cmds := supported.completions
	full := strings.TrimSpace(d.TextBeforeCursor())
	curr := d.GetWordBeforeCursor()
	spaced := strings.Contains(full, " ")

	if len(full) == 0 {
		logger.Println("Input empty returning empty suggest")
		return []prompt.Suggest{}
	}

	prefix := ""
	if spaced {
		parts := strings.Fields(full)
		prefix = parts[0]
		full = strings.Join(parts, " ")
	}
	filt := prompt.FilterHasPrefix(cmds, full, true)
	if !spaced {
		logger.Printf("Full Input: %s, Current Word: %s, IsSpaced: %t, Filtered Suggestions: %#v\n",
			full, curr, spaced, LogSafeSlice(filt))
		return filt
	}
	modFilt := make([]prompt.Suggest, 0, len(filt))
	for _, x := range filt {
		sugg := prompt.Suggest{
			Text:        strings.TrimSpace(strings.TrimPrefix(x.Text, prefix)),
			Description: x.Description,
		}
		modFilt = append(modFilt, sugg)
	}
	modFilt = prompt.FilterHasPrefix(modFilt, curr, true)
	logger.Printf("Full Input: %s, Current Word: %s, Prefix: %s, IsSpaced: %t, Filtered Suggestions: %#v\n",
		full, curr, prefix, spaced, LogSafeSlice(modFilt))
	return modFilt
}
