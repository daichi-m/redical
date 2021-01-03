package main

import (
	"os"

	"github.com/daichi-m/go-prompt"
	"github.com/fatih/color"
	"github.com/hackebrot/turtle"
)

// global is the global RedicalConf object to store all global parameters
var global RedicalConf

func main() {

	os.Setenv("GO_PROMPT_ENABLE_LOG", "1")
	var err error
	banner()

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
	kwColor := make(map[string]*color.Color)
	for _, k := range global.supported.keywords {
		kwColor[k] = color.New(color.FgHiGreen, color.Bold)
	}
	p := prompt.New(PromptAction, CmdSuggestions,
		prompt.OptionTitle("redical"),
		prompt.OptionLivePrefix(func() (string, bool) {
			return global.PromptPrefix(), true
		}),
		prompt.OptionMaxSuggestion(5),
		prompt.OptionStatusBarCallback(statusBar),
		prompt.OptionKeywordColor(color.New(color.FgHiGreen)),
		prompt.OptionKeywords(global.supported.keywords),

		/*
			prompt.OptionSelectedSuggestionBGColor(prompt.White),
			prompt.OptionSelectedSuggestionTextColor(prompt.Black),
			prompt.OptionSelectedDescriptionBGColor(prompt.Turquoise),
			prompt.OptionSelectedDescriptionTextColor(prompt.Black),

			prompt.OptionSuggestionBGColor(prompt.Cyan),
			prompt.OptionSuggestionTextColor(prompt.White),
			prompt.OptionDescriptionBGColor(prompt.DarkGray),
			prompt.OptionDescriptionTextColor(prompt.White),*/
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

func statusBar(buf *prompt.Buffer, comp *prompt.CompletionManager) (string, bool) {
	return "All systems go", true
}

func statusBarPrefSuf() (string, string) {
	smile := turtle.Search("smile")[0]
	rocket := turtle.Search("rocket")[0]
	return smile.Char, rocket.Char
}
