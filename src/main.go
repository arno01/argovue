package main

import (
	"fmt"
	"argovue/app"
)

var version = "devel"
var commit string
var builddate string

func main() {
	fmt.Printf("Starting ArgoVue, version:%s, commit:%s, builddate:%s\n", version, commit, builddate)
	app.New().Wait()
}
