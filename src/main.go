package main

import (
	"fmt"
	"kubevue/app"
)

var version = "devel"
var commit string
var builddate string

func main() {
	fmt.Printf("Starting kubevue, version:%s, commit:%s, builddate:%s\n", version, commit, builddate)
	app.New().Wait()
}
