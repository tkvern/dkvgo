package main

import (
	"dkvgo/scheduler"
)

func main() {
	scheduler.ParseCmdArgs()
	scheduler.Run()
}
