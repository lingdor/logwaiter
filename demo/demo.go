package main

import (
	"flag"
	"fmt"
	"time"
)

func main() {

	var t time.Duration
	var n int
	flag.DurationVar(&t, "t", time.Second, "split time duration")
	flag.IntVar(&n, "n", 0, "loop count")
	flag.Parse()

	i := 0
	for {
		i++
		if n != 0 && i > n {
			break
		}
		fmt.Println("number:", i)
		time.Sleep(t)
	}

}
