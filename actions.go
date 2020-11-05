package main

import (
	"encoding/base64"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/fatih/color"
	"github.com/hackebrot/turtle"
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

func render(reply interface{}) string {

	defer func() {
		if r := recover(); r != nil {
			logger.Error("Faced panic while render %v - Panic %v", reply, r)
			color.Red("Panic while rendering: %v", r)
		}
	}()

	if reply == nil {
		return color.HiGreenString("<nil>")
	}

	switch val := reply.(type) {
	case byte, int8, int16, int32, int64, int:
		return color.GreenString("%d", val)
	case float32, float64:
		return color.GreenString("%f", val)
	case bool:
		return color.GreenString("%t", val)
	case string:
		return color.GreenString("%s", val)
	case []byte:
		if len(val) < 2048 && printable(val) {
			return color.GreenString("%s", string(val))
		}
		b64 := base64.RawStdEncoding.EncodeToString(val)
		return color.GreenString("---BASE64--- %s", b64)
	}
	return renderComplex(reply)
}

func renderComplex(reply interface{}) string {

	switch val := reflect.ValueOf(reply); val.Kind() {

	case reflect.Slice:
		slc := reply.([]interface{})
		return renderSlice(slc)
	case reflect.Map:
		mp := reply.(map[interface{}]interface{})
		return renderMap(mp)
	default:
		return color.HiMagentaString("%T <", val.Interface()) +
			color.GreenString("%v", val.Interface()) +
			color.HiMagentaString(" >")
	}
}

func renderSlice(slc []interface{}) string {
	var sb strings.Builder
	var newline bool
	if len(slc) > 20 {
		newline = true
	}
	sb.WriteString(color.HiMagentaString("[ "))
	for i, x := range slc {
		var delim string

		if i < len(slc)-1 {
			if newline {
				delim = "\n"
			} else {
				delim = ", "
			}
		}
		sb.WriteString(render(x))
		sb.WriteString(color.RedString(delim))
	}
	sb.WriteString(color.HiMagentaString(" ]"))
	return sb.String()
}

func renderMap(mp map[interface{}]interface{}) string {

	var delim string

	if len(mp) > 10 {
		delim = "\n"
	} else {
		delim = ","
	}

	var kvPairs []string
	for k, v := range mp {
		kv := color.GreenString(render(k)) +
			color.HiMagentaString("->") +
			color.GreenString(render(v))
		kvPairs = append(kvPairs, kv)
	}
	return color.HiMagentaString("{") +
		strings.Join(kvPairs, delim) +
		color.HiMagentaString(" }")
}

func printable(bytes []byte) bool {
	if len(bytes) > 2048 {
		return false
	}
	return utf8.Valid(bytes)

}

func printWithEmoji(em string, data interface{}) {
	emoji := turtle.Search(em)
	var emChar string
	if len(emoji) == 0 {
		emChar = em
	} else {
		emChar = emoji[0].Char
	}

	fmt.Printf("%s  %v\n", emChar, data)
}
