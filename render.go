package main

/*
// render will render the reply from redis into a human friendly format with RT latency
func render(reply interface{}) string {
	defer func() {
		if r := recover(); r != nil {
			glg.Error("Faced panic while render %v - Panic %v", reply, r)
			color.Red("Panic while rendering: %v", r)
		}
	}()

	if reply == nil {
		return color.HiGreenString("<nil>")
	}

	switch val := reply.(type) {
	case byte, int8, int16, int32, int64, int:
		return color.GreenString("(integer) %d", val)
	case float32, float64:
		return color.GreenString("(float) %f", val)
	case bool:
		return color.GreenString("(boolean) %t", val)
	case string:
		return color.GreenString("(string) %s", val)
	case []byte:
		if s := string(val); utf8.ValidString(s) {
			return color.GreenString("(string) %s", s)
		}
		b64 := base64.RawStdEncoding.EncodeToString(val)
		return color.GreenString("(binary) %s", b64)
	}
	return renderComplex(reply)
}

// renderComplex renders complex data type like Slice and Map
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
}*/
