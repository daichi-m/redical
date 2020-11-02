package main

import (
	"fmt"

	"github.com/c-bata/go-prompt"
)

func main() {

	if err := SetupLogger(); err != nil {
		panic("Could not set up logging")
	}
	defer TearDownLogger()

	InitCmds()

	p := SetupPrompt()
	p.Run()

}

func SetupPrompt() *prompt.Prompt {
	p := prompt.New(Action, CmdSuggestions,
		prompt.OptionTitle("redical"),
		prompt.OptionLivePrefix(livePrefix),
		prompt.OptionMaxSuggestion(5),

		prompt.OptionSelectedSuggestionBGColor(prompt.White),
		prompt.OptionSelectedSuggestionTextColor(prompt.Black),
		prompt.OptionSelectedDescriptionBGColor(prompt.White),
		prompt.OptionSelectedDescriptionTextColor(prompt.DarkRed),

		prompt.OptionSuggestionBGColor(prompt.LightGray),
		prompt.OptionSuggestionTextColor(prompt.Black),
		prompt.OptionDescriptionBGColor(prompt.LightGray),
		prompt.OptionDescriptionTextColor(prompt.Cyan),
	)
	return p
}

func livePrefix() (string, bool) {
	return ">>>", true
}

func Action(txt string) {
	logger.Println("Received input for action ", txt)
	fmt.Println(txt)
}
