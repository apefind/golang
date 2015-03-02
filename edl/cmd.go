package edl

import (
	"apefind/shellutil"
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func Usage(cmd string, flags *flag.FlagSet) {
	if cmd == "extract" {
		fmt.Fprintf(os.Stderr, "\n%s %s -i <edl file> -o <csv file> -fps [24|30] [-auto]\n\n",
			filepath.Base(os.Args[0]), cmd)
		fmt.Fprintf(os.Stderr, "\tExtract information from edit decision list into csv\n\n")
	} else {
		fmt.Fprintf(os.Stderr, "%s %s\n", filepath.Base(os.Args[0]), cmd)
	}
	flags.PrintDefaults()
}

func ExtractCmd(args []string) int {

	extractCSV := func(input, output string, fps int) int {
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
			g, err := os.Create(output)
			if err != nil {
				log.Println(err)
				return 1
			}
			defer g.Close()
			w = bufio.NewWriter(g)
		}
		if err := ExtractCSV(r, w, fps); err != nil {
			log.Println(err)
			return 1
		}
		return 0
	}

	isEDLFile := func(path string) bool {
		ext := strings.ToLower(filepath.Ext(path))
		return ext == ".edl" || ext == ".txt"
	}

	var input, output string
	var fps int
	var auto bool
	flags := flag.NewFlagSet("extract", flag.ExitOnError)
	flags.Usage = func() { Usage("extract", flags) }
	flags.StringVar(&input, "input", "", "edit decision list or standard input")
	flags.StringVar(&input, "i", "", "")
	flags.StringVar(&output, "output", "", "csv output file or standard output")
	flags.StringVar(&output, "o", "", "")
	flags.IntVar(&fps, "frames-per-second", 30, "frames per second, usually 24 or 30")
	flags.IntVar(&fps, "fps", 30, "")
	flags.BoolVar(&auto, "auto", false, "automatically choose input files and output names")
	flags.BoolVar(&auto, "a", false, "")
	if err := flags.Parse(args); err != nil {
		flags.Usage()
		return 1
	}
	var F, G []string
	if auto {
		t0 := time.Now()
		log.Println("This is edtool,", t0.Format(time.ANSIC))
		log.Println("    * extracting information from edl files")
		F, G = shellutil.GetInputOutput(input, output, isEDLFile, ".csv")
		status := 0
		for i, f := range F {
			log.Println("        ", G[i])
			status |= extractCSV(f, G[i], fps)
		}
		t1 := time.Now()
		log.Printf("    * terminated at %s, total duration %s\n", t1.Format(time.ANSIC), t1.Sub(t0))
		log.Println("done")
		return status
	} else {
		return extractCSV(input, output, fps)
	}
}
