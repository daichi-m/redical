package main

import (
	"fmt"
	"strings"

	"go.uber.org/zap"
)

// RedicalAction is the template function type for executing a command on Redis.
type RedicalAction func(r *Redical, cmd string, params ...string) (string, error)

var (
	emojis    map[string]string
	actionMap map[string]RedicalAction
)

func (r *Redical) initEmojis() {
	emojis = make(map[string]string)
	emojis["ok"] = "✅"
	emojis["fail"] = "🛑"
	emojis["time"] = "⏱"
	emojis["user"] = "👤"

	actionMap = make(map[string]RedicalAction)
	actionMap["SELECT"] = selectAct
	actionMap["AUTH"] = authAct
}

func emojiFor(key string) string {
	if e, ok := emojis[key]; ok {
		return e
	}
	return ""
}

func strIntfConvert(in []string) (out []interface{}) {
	out = make([]interface{}, 0, len(in))
	for _, x := range in {
		out = append(out, x)
	}
	return
}

// Printf for pretty printing in Redical.
func (r *Redical) Printf(format string, params ...interface{}) {
	fmt.Printf("%4s "+format+"\n", "✅", params)
}

// Println for pretty printing in Redical.
func (r *Redical) Println(params ...interface{}) {
	format := strings.Builder{}
	for i := 0; i < len(params); i++ {
		format.WriteString(" %v ")
	}
	r.Printf(format.String(), params)
}

// Errorln for pretty printing errors in Redical.
func (r *Redical) Errorln(err error) {
	fmt.Printf("%4s %s \n", "🛑", err.Error())
}

// Execute executes the given command in cmd with the RedicalConf
func (r *Redical) Execute(line string) {
	defer func() {
		if rec := recover(); r != nil {
			zap.S().Warnw("Panic occured...Recovering now", "panic", rec)
			r.Errorln(fmt.Errorf("%v", rec))
		}
	}()

	if len(strings.TrimSpace(line)) == 0 {
		return
	}
	cmd, params, err := r.supported.extractCommand(line)
	if err != nil {
		r.Errorln(err)
	}
	zap.S().Debugw("Inferred details", "command", cmd, "params", params)

	if action, ok := actionMap[cmd]; ok {
		output, err := action(r, cmd, params...)
		if err != nil {
			r.Errorln(err)
			zap.S().Warnw("Error in execution", "error", err.Error())
			return
		}
		zap.S().Debugw("Output produced", "output", output)
		r.Println(output)
		return
	}

	command, ok := r.supported.kwCmd[cmd]
	if !ok {
		r.Errorln(fmt.Errorf("Did not find command %s", cmd))
		return
	}

	r.redisDB.renderSimpleNumberCommand(cmd, strIntfConvert(params)...)
	r.redisDB.renderSimpleStringCommand(cmd, strIntfConvert(params)...)

	retType := command.Return
	if action, ok := actionMap[retType]; ok {
		output, err := action(r, cmd, params...)
		if err != nil {
			r.Errorln(err)
			zap.S().Errorw("Error in execution", "error", err)
			return
		}
		zap.S().Debugw("Output produced", "output", output)
		r.Println(output)
		return
	}
	zap.S().Warnw("Could not get action from command or return type.", "line", line)
}

// selectAct is an instance of RedicalAction for SELECT call in Redis
func selectAct(r *Redical, cmd string, params ...string) (msg string, err error) {
	if db, ok := ExtractInt(params, 0); ok {
		err = r.SwitchDB(db)
		if err == nil {
			msg = fmt.Sprintf("Switched to DB %d", db)
		}
		return
	}
	return "", fmt.Errorf("Could not extract DB from params")
}

// selectAct is an instance of RedicalAction for AUTH call in Redis
func authAct(r *Redical, cmd string, params ...string) (msg string, err error) {
	if pass, ok := ExtractStr(params, 0); ok {
		err = r.Authenticate(pass)
		if err == nil {
			msg = "Authentication successful"
		}
		return
	}
	return "", fmt.Errorf("Could not extract password from params")
}

// func handleError(err error) {
// 	stack := debug.Stack()
// 	glg.Error("Error occured: %s at %s", err.Error(), string(stack))
// }

func justifyOutput(emoji string, msg string) string {
	return fmt.Sprintf("%2s %-4s %v", "", emoji, msg)
}

func (r *RedisDB) renderSimpleStringCommand(cmd string, params ...interface{}) string {
	zap.S().Debugw("Execution context", "command", cmd, "params", params)
	reply, err := r.redisConn.Do(cmd, params...)
	if err != nil {
		return justifyOutput(emojiFor("fail"), err.Error())
	}
	switch r := reply.(type) {
	case string:
		return justifyOutput(emojiFor("success"), r)
	default:
		panic("Unexpected type of reply, expected string")
	}
}

func (r *RedisDB) renderSimpleNumberCommand(cmd string, params ...interface{}) string {
	zap.S().Debugw("Execution context", "command", cmd, "params", params)
	reply, err := r.redisConn.Do(cmd, params...)
	if err != nil {
		return justifyOutput(emojiFor("error"), err.Error())
	}
	switch r := reply.(type) {
	case byte, int8, int16, int32, int64:
		return justifyOutput(emojiFor("success"), fmt.Sprintf("%5d", r))
	case float32, float64:
		return justifyOutput(emojiFor("success"), fmt.Sprintf("%5f", r))
	default:
		panic("Unexpected type of reply, expected int or float")
	}
}

// // aclLoad runs a ACL LOAD command on redis and responds back with the appropriate response.
// func (r *RedisDB) aclLoad(cmd string, params ...string) (string, error) {
// 	glg.Debug("Executing ACL SAVE, ignoring params")
// 	reply, err := r.redisConn.Do(cmd)
// 	if err != nil {
// 		return "", err
// 	}
// 	return fmt.Sprint(reply), nil
// }

// // aclSave runs a ACL SAVE command on redis and responds back with the appropriate response.
// func (r *RedisDB) aclSave(cmd string, params ...string) (string, error) {
// 	glg.Debug("Executing ACL SAVE, ignoring params")
// 	reply, err := r.redisConn.Do(cmd)
// 	if err != nil {
// 		return "", err
// 	}
// 	return fmt.Sprint(reply), nil
// }

// aclList runs a ACL LIST on redis and responds back with the ACL's set
// func (r *RedisDB) aclList(cmd string, params ...string) (string, error) {

// }

// // PromptAction takes the command from the prompt and runs the command on redis
// func PromptAction(cmd string) {
// 	if len(strings.TrimSpace(cmd)) == 0 {
// 		return
// 	}

// 	if global.redisDB.redisConn == nil {
// 		color.Red("Redis connection has not been initialized.\n")
// 		return
// 	}

// 	defer func(start time.Time) {
// 		duration := time.Now().Sub(start)
// 		durStr := color.HiYellowString("%v", duration)
// 		printWithEmoji("stopwatch", durStr)

// 	}(time.Now())

// 	if err := action(cmd); err != nil {
// 		errOut := color.RedString("%s", err.Error())
// 		printWithEmoji("no_entry", errOut)
// 	}
// }

// func action(cmd string) error {
// 	r := global.redisDB.redisConn
// 	parts := strings.Fields(cmd)
// 	var rep interface{} = nil
// 	var err error = nil
// 	var args []interface{} = make([]interface{}, 0, len(parts))

// 	for _, a := range parts {
// 		args = append(args, a)
// 	}

// 	if len(args) == 1 {
// 		ca := inferCustomAct(args[0].(string))
// 		if ca != nil {
// 			return ca()
// 		}
// 		rep, err = r.Do(args[0].(string))
// 	} else if len(args) > 1 {
// 		ca := inferCustomAct(args[0].(string))
// 		if ca != nil {
// 			return ca(args[1:]...)
// 		}
// 		rep, err = r.Do(args[0].(string), args[1:]...)
// 	}
// 	logger.Debug("Reply: %v, Error: %v", rep, err)
// 	if err == nil {
// 		rendered := render(rep)
// 		printWithEmoji("white_check_mark", rendered)
// 	}
// 	return err
// }

// func inferCustomAct(cmd string) customAct {

// 	switch strings.ToUpper(cmd) {
// 	case "AUTH":
// 		return redisAuth
// 	case "SELECT":
// 		return redisSelect
// 	default:
// 		return nil
// 	}
// }

// type customAct func(args ...interface{}) error

// // redisAuth is an action that changes the password of the redis connection and refreshes it.
// func redisAuth(args ...interface{}) error {
// 	if len(args) != 1 {
// 		return fmt.Errorf("Unexpected number of arguments. Expected %d, found %d", 1, len(args))
// 	}

// 	password := args[0].(string)
// 	if err := global.ModifyConfig(&DBConfig{password: password}); err != nil {
// 		return err
// 	}
// 	color.Green("Redis connection refreshed with password %s\n", "****(redacted)")
// 	return nil
// }

// // redisSelect action changes the database of the redis connection and refreshes it.
// func redisSelect(args ...interface{}) error {

// 	if len(args) != 1 {
// 		return fmt.Errorf("Unexpected number of arguments. Expected %d, found %d", 1, len(args))
// 	}
// 	dbstr := args[0].(string)
// 	db, err := strconv.Atoi(dbstr)
// 	if err != nil {
// 		return fmt.Errorf("Unable convert db id to integer due to %s", err.Error())
// 	}
// 	if err := global.ModifyConfig(&DBConfig{database: db}); err != nil {
// 		return err
// 	}
// 	color.Green("Redis connection refreshed to connect to db %d\n", db)
// 	return nil
// }
