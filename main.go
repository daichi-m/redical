package main

import (
	"github.com/c-bata/go-prompt"
	"github.com/fatih/color"
)

// global is the global RedicalConf object to store all global parameters
var global RedicalConf

func main() {

	banner()

	var err error
	global = RedicalConf{}
	global.redisDB.DBConfig = ParseConfig()

	if err = SetupLogger(); err != nil {
		panic("Could not set up logging")
	}
	defer TearDownLogger()
	logger.Info("Logger initialized")

	if err := global.redisDB.InitializeRedis(); err != nil {
		color.Red("Redis could not be initialized")
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
	p := prompt.New(PromptAction, CmdSuggestions,
		prompt.OptionTitle("redical"),
		prompt.OptionLivePrefix(func() (string, bool) {
			return global.PromptPrefix(), true
		}),
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

func banner() {
	banner := `
*** Welcome to RediCal - the all new replacement to redis-cli ***
   ___           __ _  _____       __
  / _ \ ___  ___/ /(_)/ ___/___ _ / /
 / , _// -_)/ _  // // /__ / _ '// / 
/_/|_| \__/ \_,_//_/ \___/ \_,_//_/  
								   
`
	color.Cyan(banner)
}
