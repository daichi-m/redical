package main

import (
	"fmt"

	"github.com/c-bata/go-prompt"
)

// global is the global RedicalConf object to store all global parameters
var global RedicalConf

func main() {

	var err error
	global = RedicalConf{}
	global.config = ParseConfig()

	if err = SetupLogger(); err != nil {
		panic("Could not set up logging")
	}
	defer TearDownLogger()
	logger.Info("Logger initialized")

	r, err := global.config.InitializeRedis()
	if err != nil {
		panic("Redis could not be initialized")
	}
	global.redis = &r
	logger.Info("CLI inputs parsed and redis-client initialized")

	InitCmds()
	p := SetupPrompt()
	logger.Info("Prompt setup complete, initialize prompt now")
	p.Run()

}

// SetupPrompt sets up the CLI Prompt to run with proper prompt.CompletionManager and prompt.Executor
func SetupPrompt() *prompt.Prompt {
	p := prompt.New(action, CmdSuggestions,
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
	return ">>> ", true
}

func action(txt string) {
	logger.Debug("Received input for action %s", txt)
	fmt.Println(txt)
}
