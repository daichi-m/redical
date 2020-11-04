package main

import (
	"fmt"

	"github.com/c-bata/go-prompt"
	"github.com/fatih/color"
)

// global is the global RedicalConf object to store all global parameters
var global RedicalConf

func main() {

	color.New(color.FgHiRed, color.BgCyan, color.Bold).Println("Hello, World\n")

	var err error
	global = RedicalConf{}
	global.redisDB.DBConfig = ParseConfig()

	if err = SetupLogger(); err != nil {
		panic("Could not set up logging")
	}
	defer TearDownLogger()
	logger.Info("Logger initialized")

	if err := global.redisDB.InitializeRedis(); err != nil {
		panic("Redis could not be initialized")
	}
	defer global.redisDB.TearDownRedis()
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
		prompt.OptionSelectedDescriptionBGColor(prompt.Turquoise),
		prompt.OptionSelectedDescriptionTextColor(prompt.Black),

		prompt.OptionSuggestionBGColor(prompt.Cyan),
		prompt.OptionSuggestionTextColor(prompt.White),
		prompt.OptionDescriptionBGColor(prompt.DarkGray),
		prompt.OptionDescriptionTextColor(prompt.White),
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
