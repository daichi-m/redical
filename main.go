package main

import (
	"fmt"
	"io"

	"github.com/daichi-m/go-prompt"
	"github.com/fatih/color"
	"go.uber.org/zap"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

// // global is the global RedicalConf object to store all global parameters
// var global RedicalConf
func main() {
	// Setup logging
	logFile, logger := SetupLogger()
	defer logFile.Close()
	defer func() {
		_ = logger.Sync()
	}()
	undo := zap.ReplaceGlobals(logger)
	defer undo()

	banner()

	// Initialize the configs
	config, err := InitializeRedical()
	if err != nil {
		panic(fmt.Sprintf("Could not initialize config: %s", err.Error()))
	}
	defer config.Close()

	// Create redis connection
	err = config.redisDB.InitializeRedis()
	if err != nil {
		zap.S().Warn("Redis connection not initialized, starting CLI without active redis")
	}
	defer config.redisDB.TearDownRedis()

	// Initialize the set of emojis
	config.initEmojis()

	// Setup CLI prompt
	p := SetupPrompt(config)
	logger.Info("Prompt setup complete, initialize prompt now")
	p.Run()
}

// SetupPrompt sets up the CLI Prompt to run with proper prompt.CompletionManager and prompt.Executor
func SetupPrompt(r *Redical) *prompt.Prompt {
	p := prompt.New(r.Execute,
		func(d prompt.Document) []prompt.Suggest {
			/*cmds := r.supported.completions
			return r.CmdSuggestions(cmds, d)*/
			return nil
		},
		prompt.OptionTitle("redical"),
		prompt.OptionLivePrefix(func() (string, bool) {
			return r.PromptPrefix(), true
		}),
		prompt.OptionMaxSuggestion(5),
		prompt.OptionStatusBarCallback(statusBar),
		prompt.OptionKeywordColor(color.New(color.FgHiGreen)),
		prompt.OptionKeywords(r.supported.Keywords()),
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

// SetupLogger sets up the logging params
func SetupLogger() (io.WriteCloser, *zap.Logger) {
	logFile := &lumberjack.Logger{
		Filename:   "logs/redical.log",
		MaxSize:    5, // megabytes
		MaxBackups: 3,
		MaxAge:     2,    // days
		Compress:   true, // disabled by default
		LocalTime:  false,
	}
	logger, err := zap.NewDevelopmentConfig().Build(zap.IncreaseLevel(zap.DebugLevel))
	if err != nil {
		panic("Could not initialize logger")
	}
	return logFile, logger
}

/*
func statusBarPrefSuf() (string, string) {
	smile := turtle.Search("smile")[0]
	rocket := turtle.Search("rocket")[0]
	return smile.Char, rocket.Char
}*/
