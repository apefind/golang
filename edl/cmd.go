package edl

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func Usage(cmd string, flags *flag.FlagSet) {
	if cmd == "extract" {
		fmt.Fprintf(os.Stderr, "\n%s %s -i <edl file> -o <csv file> -fps [24|30]\n\n",
			filepath.Base(os.Args[0]), cmd)
		fmt.Fprintf(os.Stderr, "\tExtract information from edit decision list into csv\n\n")
	} else {
		fmt.Fprintf(os.Stderr, "%s %s\n", filepath.Base(os.Args[0]), cmd)
	}
	flags.PrintDefaults()
}

func ExtractCmd(args []string) int {
	var input, output string
	var fps int
	flags := flag.NewFlagSet("extract", flag.ExitOnError)
	flags.Usage = func() { Usage("extract", flags) }
	flags.StringVar(&input, "input", "", "edit decision list or standard input")
	flags.StringVar(&input, "i", "", "")
	flags.StringVar(&output, "output", "", "csv output file or standard output")
	flags.StringVar(&output, "o", "", "")
	flags.IntVar(&fps, "frames-per-second", 30, "frames per second, usually 24 or 30")
	flags.IntVar(&fps, "fps", 30, "")
	if err := flags.Parse(args); err != nil {
		flags.Usage()
		return 1
	}
	var r *bufio.Reader
	if input == "" {
		r = bufio.NewReader(os.Stdin)
	} else {
		f, err := os.Open(input)
		if err != nil {
			log.Println(err)
			return 1
		}
		defer f.Close()
		r = bufio.NewReader(f)
	}
	var w *bufio.Writer
	if output == "" {
		w = bufio.NewWriter(os.Stdout)
	} else {
		f, err := os.Create(output)
		if err != nil {
			log.Println(err)
			return 1
		}
		defer f.Close()
		w = bufio.NewWriter(f)
	}
	if err := ExtractCSV(r, w, fps); err != nil {
		log.Println(err)
		return 1
	}
	return 0
}
