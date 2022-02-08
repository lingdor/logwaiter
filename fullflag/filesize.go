package fullflag

import (
	"fmt"
	"strconv"
	"strings"
)

type fileSizeValue int

func newFileSizeValue(val int, p *int) *fileSizeValue {
	*p = val
	return (*fileSizeValue)(p)
}

func isNum(ch int) bool {
	var low int = '0'
	var sup int = '9'
	return ch >= low && ch <= sup
}

func (f *fileSizeValue) Set(s string) (err error) {
	unit := "b"
	var num int
	for i := len(s) - 1; i > -1; i-- {
		var ch rune = rune(s[i])
		if isNum(int(ch)) {
			unit = s[i+1:]
			if num, err = strconv.Atoi(s[:i+1]); err != nil {
				return fmt.Errorf("file size %s not format", s)
			}
			break
		}
	}
	unit = strings.ToLower(unit)
	switch unit {
	case "k", "kb":
		num *= 1024
		break
	case "m", "mb":
		num *= 1024 * 1024
		break
	case "g", "gb":
		num *= 1024 * 1024 * 1024
		break
	}
	*f = fileSizeValue(num)
	return nil
}

func (i *fileSizeValue) Get() interface{} { return int(*i) }

func (i *fileSizeValue) String() string { return strconv.Itoa(int(*i)) }
