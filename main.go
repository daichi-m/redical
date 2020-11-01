package main

import (
	"fmt"
	"log"
	"os"

	"github.com/c-bata/go-prompt"
)

func main() {

	var posixWriter *prompt.PosixWriter = prompt.NewStandardOutputWriter().(*prompt.PosixWriter)
	var posixParser *prompt.PosixParser = prompt.NewStandardInputParser()

	if err := SetupLogger(); err != nil {
		panic("Could not set up logging")
	}
	defer TearDownLogger()

	InitCmds()

	/*
		count := 0
		for _, rc := range supported.Cmds {
			if strings.Contains(rc.Name, " ") {
				logger.Println(rc.Name)
				count++
			}
		}
		logger.Println("Count = ", count) */

	kb := prompt.KeyBind{
		Key: prompt.Escape,
		Fn:  esc,
	}

	kb2 := prompt.KeyBind{
		Key: prompt.Escape,
		Fn:  esc,
	}

	p := prompt.New(print, completer,
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

		prompt.OptionWriter(posixWriter),
		prompt.OptionParser(posixParser),

		prompt.OptionAddKeyBind(kb, kb2),
	)
	p.Run()

}

func print(s string) {
	if len(s) != 0 {
		fmt.Println(s)
	}
}

func livePrefix() (string, bool) {
	return fmt.Sprintf("(%d) >>> ", 0), true
}

func esc(b *prompt.Buffer) {

}

func setupLogs() (*log.Logger, *os.File, error) {
	f, err := os.OpenFile("logs/prompter.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644)

	if err != nil {
		log.Panicln("Cannot open logs file")
		return nil, nil, err
	}

	logger := log.New(f, "prompter", log.LstdFlags)
	return logger, f, nil

}
