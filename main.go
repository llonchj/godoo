package main

//go:generate esc -o generator/static.go -pkg generator tmpl types api

import (
	"github.com/llonchj/godoo/cmd"
)

func main() {
	cmd.Execute()
}
