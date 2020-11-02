package main

import (
	"fmt"

	"github.com/c-bata/go-prompt"
	"github.com/gomodule/redigo/redis"
)

// RedicalConf is the global configuration struct to encapsulate all global parameters
type RedicalConf struct {
	config    DBConfig
	supported CommandList
	redis     *redis.Conn
}

// ModifyConfig modifies the DBConfig for redis and refreshes the global redis client.
func (rc *RedicalConf) ModifyConfig(config DBConfig) error {
	rc.config = config
	r, err := config.InitializeRedis()
	if err != nil {
		return err
	}
	rc.redis = &r
	return nil
}

// global is the global RedicalConf object to store all global parameters
var global RedicalConf

func main() {

	var err error
	global = RedicalConf{}

	if err = SetupLogger(); err != nil {
		panic("Could not set up logging")
	}
	defer TearDownLogger()

	InitCmds()
	ParseConfig()
	r, err := global.config.InitializeRedis()
	global.redis = &r

	p := SetupPrompt()
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
	logger.Println("Received input for action ", txt)
	fmt.Println(txt)
}
