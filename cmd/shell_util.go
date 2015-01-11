package main

import (
	"apefind/shellutil"
	"flag"
	"log"
	"os"
)

func main() {
	log.SetFlags(0)
	flag.Usage = func() { shellutil.Usage("[prompt|timeout]", flag.CommandLine) }
	flag.Parse()
	switch flag.Arg(0) {
	case "prompt":
		os.Exit(shellutil.PromptCmd(os.Args[2:]))
	case "timeout":
		os.Exit(shellutil.TimeoutCmd(os.Args[2:]))
	case "batch":
		os.Exit(shellutil.BatchCmd(os.Args[2:]))
	default:
		flag.Usage()
		os.Exit(1)
	}
}
