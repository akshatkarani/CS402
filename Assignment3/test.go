package main

import (
	"fmt"
	"time"
)

func main() {
	dur := "100011754180000" + "ns"
	fmt.Println(time.ParseDuration(dur))
}
