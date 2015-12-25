package main

import (
	"flag"
	"fmt"
	"github.com/apefind/golang/episodeguide"
	"github.com/apefind/golang/shellutil"
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
	var info episodeguide.SeriesInfoCmd
	switch flag.Arg(0) {
	case "rename":
		info.Title, info.SeasonID = episodeguide.GetSeriesTitleFromWorkingDirectory()
		flags := flag.NewFlagSet("rename", flag.ExitOnError)
		flags.BoolVar(&info.DryRun, "dry-run", false, "just print, do not actually rename")
		flags.BoolVar(&info.DryRun, "d", false, "short for -dry-run")
		flags.BoolVar(&info.NormalizedTitle, "no-title", false, "ignore the title, just use S01E01, S01E02, ...")
		flags.BoolVar(&info.NormalizedTitle, "n", false, "short for -no-title")
		flags.StringVar(&info.Method, "method", "tvmaze|tvrage", "tvmaze or/and tvrage")
		flags.StringVar(&info.Method, "m", "tvmaze|tvrage", "short for -method")
		flags.DurationVar(&info.Timeout, "timeout", 10*time.Second, "stop search after given duration")
		flags.DurationVar(&info.Timeout, "t", 10*time.Second, "")
		flags.Usage = func() { usage(flags) }
		flags.Parse(os.Args[2:])
		path, _ := os.Getwd()
		info.RenameEpisodes(filepath.Clean(path))
	case "info":
		flags := flag.NewFlagSet("info", flag.ExitOnError)
		flags.StringVar(&info.Method, "method", "tvmaze|tvrage", "tvmaze or/and tvrage")
		flags.StringVar(&info.Method, "m", "tvmaze|tvrage", "")
		flags.DurationVar(&info.Timeout, "timeout", 10*time.Second, "stop search after given duration")
		flags.DurationVar(&info.Timeout, "t", 10*time.Second, "")
		flags.Usage = func() { usage(flags) }
		flags.Parse(os.Args[2:])
		dirs, err := shellutil.GetDirsFromFlagSetArgs(flags)
		if err != nil {
			log.Println(err)
			os.Exit(-1)
		}
		for _, dir := range dirs {
			info.Title, info.SeasonID = episodeguide.GetSeriesTitleFromPath(dir)
			info.ListEpisodes()
		}
	default:
		flag.Usage()
		os.Exit(1)
	}
}
