package main

import (
	"apefind/edl"
	"flag"
	"log"
	"os"
)

func main() {
	log.SetFlags(0)
	flag.Usage = func() { edl.Usage("[convert]", flag.CommandLine) }
	flag.Parse()
	switch flag.Arg(0) {
	case "convert":
		os.Exit(edl.ConvertCmd(os.Args[2:]))
	default:
		flag.Usage()
		os.Exit(1)
	}
}
