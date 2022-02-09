package common

import (
	"fmt"
	"os"
)

func AppRecover() {
	err := recover()
	if err != nil {
		fmt.Println("logwaiter running error:")
		fmt.Println(err)
		os.Exit(1)
	}
}

func CheckPanic(err error) {

	if err != nil {
		panic(err)
	}
}
