package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
)

// PromptAction takes the command from the prompt and runs the command on redis
func PromptAction(cmd string) {
	if len(strings.TrimSpace(cmd)) == 0 {
		return
	}

	if global.redisDB.redisConn == nil {
		color.Red("Redis connection has not been initialized.\n")
		return
	}

	defer func(start time.Time) {
		duration := time.Now().Sub(start)
		durStr := color.HiYellowString("%v", duration)
		printWithEmoji("stopwatch", durStr)

	}(time.Now())

	if err := action(cmd); err != nil {
		errOut := color.RedString("%s", err.Error())
		printWithEmoji("no_entry", errOut)
	}
}

func action(cmd string) error {
	r := global.redisDB.redisConn
	parts := strings.Fields(cmd)
	var rep interface{} = nil
	var err error = nil
	var args []interface{} = make([]interface{}, 0, len(parts))

	for _, a := range parts {
		args = append(args, a)
	}

	if len(args) == 1 {
		ca := inferCustomAct(args[0].(string))
		if ca != nil {
			return ca()
		}
		rep, err = r.Do(args[0].(string))
	} else if len(args) > 1 {
		ca := inferCustomAct(args[0].(string))
		if ca != nil {
			return ca(args[1:]...)
		}
		rep, err = r.Do(args[0].(string), args[1:]...)
	}
	logger.Debug("Reply: %v, Error: %v", rep, err)
	if err == nil {
		rendered := render(rep)
		printWithEmoji("white_check_mark", rendered)
	}
	return err
}

func inferCustomAct(cmd string) customAct {

	switch strings.ToUpper(cmd) {
	case "AUTH":
		return redisAuth
	case "SELECT":
		return redisSelect
	default:
		return nil
	}
}

type customAct func(args ...interface{}) error

// redisAuth is an action that changes the password of the redis connection and refreshes it.
func redisAuth(args ...interface{}) error {
	if len(args) != 1 {
		return fmt.Errorf("Unexpected number of arguments. Expected %d, found %d", 1, len(args))
	}

	password := args[0].(string)
	if err := global.ModifyConfig(&DBConfig{password: password}); err != nil {
		return err
	}
	color.Green("Redis connection refreshed with password %s\n", "****(redacted)")
	return nil
}

// redisSelect action changes the database of the redis connection and refreshes it.
func redisSelect(args ...interface{}) error {

	if len(args) != 1 {
		return fmt.Errorf("Unexpected number of arguments. Expected %d, found %d", 1, len(args))
	}
	dbstr := args[0].(string)
	db, err := strconv.Atoi(dbstr)
	if err != nil {
		return fmt.Errorf("Unable convert db id to integer due to %s", err.Error())
	}
	if err := global.ModifyConfig(&DBConfig{database: db}); err != nil {
		return err
	}
	color.Green("Redis connection refreshed to connect to db %d\n", db)
	return nil
}
