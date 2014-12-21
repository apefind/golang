package main

import (
	"apefind/edl"
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func usage(flags *flag.FlagSet) {
	fmt.Println("usage:", filepath.Base(os.Args[0]), "[convert]")
	flags.PrintDefaults()
}

func convert(input string, output string, fps int) error {
	var r *bufio.Reader
	if input == "stdin" {
		r = bufio.NewReader(os.Stdin)
	} else {
		f, err := os.Open(input)
		if err != nil {
			return err
		}
		defer f.Close()
		r = bufio.NewReader(f)
	}
	var w *bufio.Writer
	if output == "stdout" {
		w = bufio.NewWriter(os.Stdout)
	} else {
		f, err := os.Create(output)
		if err != nil {
			return err
		}
		defer f.Close()
		w = bufio.NewWriter(f)
	}
	return edl.ConvertToCSV(r, w, fps)
}

func main() {
	flag.Usage = func() { usage(flag.CommandLine) }
	flag.Parse()
	switch flag.Arg(0) {
	case "convert":
		var input, output string
		var fps int
		flags := flag.NewFlagSet("convert", flag.ExitOnError)
		flags.StringVar(&input, "input", "stdin", "standard input, file or directory")
		flags.StringVar(&output, "output", "stdout", "standard output, file or directory")
		flags.IntVar(&fps, "fps", 30, "frames per second, usually 24 or 30")
		flags.Parse(os.Args[2:])
		if err := convert(input, output, fps); err != nil {
			fmt.Println(err)
		}
	default:
		flag.Usage()
		os.Exit(1)
	}
}
