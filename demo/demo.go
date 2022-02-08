package main

import (
	"fmt"
	"time"
)

func main() {

	i := 0
	for {
		i++
		fmt.Println("number:", i)
		time.Sleep(time.Second)
	}

}
