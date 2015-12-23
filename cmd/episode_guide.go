package main

import (
	"apefind/episodeguide"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

func usage(flags *flag.FlagSet) {
	fmt.Println("usage:", filepath.Base(os.Args[0]), "[rename|info]")
	flags.PrintDefaults()
}

func main() {
	log.SetFlags(0)
	flag.Usage = func() { usage(flag.CommandLine) }
	flag.Parse()
	path, _ := os.Getwd()
	switch flag.Arg(0) {
	case "rename":
		var dryRun, noTitle bool
		var method string
		var timeout time.Duration
		flags := flag.NewFlagSet("rename", flag.ExitOnError)
		flags.BoolVar(&dryRun, "dry-run", false, "just print, do not actually rename")
		flags.BoolVar(&dryRun, "d", false, "")
		flags.BoolVar(&noTitle, "no-title", false, "ignore the title, just use S01E01, S01E02, ...")
		flags.BoolVar(&noTitle, "n", false, "")
		flags.StringVar(&episodeguide.Method, "method", "tvmaze|tvrage", "tvmaze or/and tvrage")
		flags.StringVar(&episodeguide.Method, "m", "tvmaze|tvrage", "")
		flags.DurationVar(&timeout, "timeout", 10*time.Second, "stop search after given duration")
		flags.DurationVar(&timeout, "t", 10*time.Second, "")
		flags.Usage = func() { usage(flags) }
		flags.Parse(os.Args[2:])
		episodeguide.RenameEpisodes(filepath.Clean(path), method, dryRun, noTitle)
	case "info":
		var method string
		var timeout time.Duration
		flags := flag.NewFlagSet("info", flag.ExitOnError)
		flags.StringVar(&episodeguide.Method, "method", "tvmaze|tvrage", "tvmaze or/and tvrage")
		flags.StringVar(&episodeguide.Method, "m", "tvmaze|tvrage", "")
		flags.DurationVar(&timeout, "timeout", 10*time.Second, "stop search after given duration")
		flags.DurationVar(&timeout, "t", 10*time.Second, "")
		flags.Usage = func() { usage(flags) }
		flags.Parse(os.Args[2:])
		episodeguide.ListEpisodes(filepath.Clean(path), method)
	default:
		flag.Usage()
		os.Exit(1)
	}
}
