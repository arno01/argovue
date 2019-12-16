package main

import (
	"fmt"
	"kubevue/app"
)

// Default is `-s -w -X main.version={{.Version}} -X main.commit={{.ShortCommit}} -X main.date={{.Date}} -X main.builtBy=goreleaser`.
var version string
var commit string
var date string
var builtBy string

func main() {
	fmt.Printf("Starting kubevue, version:%s", version)
	a := app.New()
	a.Connect()
	a.Serve()
}
