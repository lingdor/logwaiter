package fullflag

import (
	"flag"
	"fmt"
	"time"
)

var params map[string]string = make(map[string]string, 0)
var maxlen int = 0

func StringVar(p *string, name string, short, value, usage string) {
	flag.StringVar(p, name, value, "")
	if short != "" {
		flag.StringVar(p, short, value, "")
	}
	AddUsage(name, short, usage)
}
func BoolVar(p *bool, name string, short string, value bool, usage string) {
	flag.BoolVar(p, name, value, "")
	if short != "" {
		flag.BoolVar(p, short, value, "")
	}
	AddUsage(name, short, usage)
}
func DurationVar(p *time.Duration, name string, short string, value time.Duration, usage string) {
	flag.DurationVar(p, name, value, "")
	if short != "" {
		flag.DurationVar(p, short, value, "")
	}
	AddUsage(name, short, usage)
}
func IntVar(p *int, name string, short string, value int, usage string) {
	flag.IntVar(p, name, value, "")
	if short != "" {
		flag.IntVar(p, short, value, "")
	}
	AddUsage(name, short, usage)
}
func FileSizeVar(p *int, name string, short string, value int, usage string) {
	flag.CommandLine.Var(newFileSizeValue(value, p), name, usage)
	AddUsage(name, short, usage)
}
func AddUsage(name string, short, usage string) {
	var key string
	if len(short) > 0 {
		key = fmt.Sprintf("-%s, --%s", short, name)
	} else {
		key = fmt.Sprintf("--%s", name)
	}
	if len(key) > maxlen {
		maxlen = len(key)
	}
	params[key] = usage
}
func Parse() {
	flag.Parse()
}
func Usage() {
	fullLen := maxlen
	fmt.Println("Usage of LogWaiter")
	for k, v := range params {
		bsKey := []byte(k)
		fullBS := make([]byte, fullLen)
		copy(fullBS, bsKey)
		for i := len(bsKey); i < fullLen; i++ {
			fullBS[i] = byte(' ')
		}
		fmt.Println(string(fullBS), v)
	}
}
