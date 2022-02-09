package main

import (
	"flag"
	"fmt"
	"time"
)

func main() {

	var t time.Duration
	flag.DurationVar(&t, "t", time.Second, "split time duration")
	flag.Parse()

	i := 0
	for {
		i++
		fmt.Println("number:", i)
		time.Sleep(t)
	}

}
