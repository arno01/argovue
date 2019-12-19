package main

import (
	"fmt"
	"kubevue/app"
)

var version = "devel"
var builddate string

func main() {
	fmt.Printf("Starting kubevue, version:%s builddate:%s\n", version, builddate)
	app.New().Wait()
}
