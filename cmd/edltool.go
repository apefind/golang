package main

import (
	"apefind/edl"
	"flag"
	"log"
	"os"
)

func main() {
	log.SetFlags(0)
	flag.Usage = func() { edl.Usage("[extract]", flag.CommandLine) }
	flag.Parse()
	switch flag.Arg(0) {
	case "extract":
		os.Exit(edl.ExtractCmd(os.Args[2:]))
	default:
		flag.Usage()
		os.Exit(1)
	}
}
