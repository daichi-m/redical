package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/daichi-m/go-prompt"
	"github.com/fatih/color"
	"github.com/hackebrot/turtle"
	"github.com/kpango/glg"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

// // global is the global RedicalConf object to store all global parameters
// var global RedicalConf
func main() {
	// Setup logging
	logFile, logger := SetupLogger()
	defer logFile.Close()

	banner()

	// Initialize the configs
	config, err := InitializeRedical()
	if err != nil {
		panic(fmt.Sprintf("Could not initialize config: %s", err.Error()))
	}
	defer config.Close()

	// Create redis connection
	config.redisDB.InitializeRedis()
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
	kw := make([]string, 0)
	for k := range r.supported.keywordCommands {
		kw = append(kw, strings.Fields(k)...)
	}
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
		prompt.OptionKeywords(kw),
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
func SetupLogger() (io.WriteCloser, *glg.Glg) {
	logFile := &lumberjack.Logger{
		Filename:   "logs/redical.log",
		MaxSize:    5, // megabytes
		MaxBackups: 3,
		MaxAge:     2,    // days
		Compress:   true, // disabled by default
		LocalTime:  false,
	}
	logger := glg.Get().
		SetMode(glg.WRITER).
		InitWriter().
		AddWriter(logFile).
		// AddLevelWriter(glg.DEBG, logFile).
		// AddLevelWriter(glg.INFO, logFile).
		// AddLevelWriter(glg.WARN, logFile).
		// AddLevelWriter(glg.ERR, logFile).
		DisableJSON().
		DisableColor()
	return logFile, logger
}

func statusBarPrefSuf() (string, string) {
	smile := turtle.Search("smile")[0]
	rocket := turtle.Search("rocket")[0]
	return smile.Char, rocket.Char
}
