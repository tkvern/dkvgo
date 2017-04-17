package main

import (
	"os"

	"dkvgo/worker"
)

func main() {
	var opts = worker.NewOptions()
	opts.TryFromCmdArgs(os.Args[1:])
	var workerd = worker.NewDkvWorker(opts)
	workerd.Main()
}
