package main

import (
	"www.github.com/apefind/golang/shellutil"
	"flag"
	//"io/ioutil"
	"log"
	"os"
)

func main() {
	log.SetFlags(0)
	//log.SetOutput(ioutil.Discard)
	flag.Usage = func() { shellutil.Usage("[prompt|timeout|timeit|cleanup]", flag.CommandLine) }
	flag.Parse()
	switch flag.Arg(0) {
	case "prompt":
		os.Exit(shellutil.PromptCmd(os.Args[2:]))
	case "timeout":
		os.Exit(shellutil.TimeOutCmd(os.Args[2:]))
	case "timeit":
		os.Exit(shellutil.TimeItCmd(os.Args[2:]))
	case "cleanup":
		os.Exit(shellutil.CleanUpCmd(os.Args[2:]))
	default:
		flag.Usage()
		os.Exit(1)
	}
}
