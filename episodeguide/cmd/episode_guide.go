package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/apefind/golang/episodeguide"
	"github.com/apefind/golang/shellutil"
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
		flags.StringVar(&info.Source, "source", "tvmaze|tvrage|redis", "tvmaze or/and tvrage")
		flags.DurationVar(&info.Timeout, "timeout", 10*time.Second, "stop search after given duration")
		flags.StringVar(&info.Title, "series", "", "series title")
		flags.Usage = func() { usage(flags) }
		flags.Parse(os.Args[2:])
		dirs, err := shellutil.GetDirsFromFlagSetArgs(flags)
		if err != nil {
			log.Println(err)
			os.Exit(-1)
		}
		for _, dir := range dirs {
			if info.Title == "" {
				info.Title, _ = episodeguide.GetSeriesTitleFromPath(dir)
			}
			info.RenameEpisodes(filepath.Clean(dir))
		}
	case "info":
		flags := flag.NewFlagSet("info", flag.ExitOnError)
		flags.StringVar(&info.Source, "source", "tvmaze|tvrage", "tvmaze or/and tvrage")
		flags.DurationVar(&info.Timeout, "timeout", 10*time.Second, "stop search after given duration")
		flags.StringVar(&info.Title, "series", "", "series title")
		flags.IntVar(&info.SeasonID, "season", 0, "series season")
		flags.Usage = func() { usage(flags) }
		flags.Parse(os.Args[2:])
		if info.Title == "" {
			dirs, err := shellutil.GetDirsFromFlagSetArgs(flags)
			if err != nil {
				log.Println(err)
				os.Exit(-1)
			}
			for _, dir := range dirs {
				info.Title, info.SeasonID = episodeguide.GetSeriesTitleFromPath(dir)
				info.ListEpisodes()
			}
		} else {
			info.ListEpisodes()
		}
	default:
		flag.Usage()
		os.Exit(1)
	}
}
