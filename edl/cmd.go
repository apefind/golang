package edl

import (
	"www.github.com/apefind/golang/shellutil"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
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
		inout := shellutil.NewInOut(input, output)
		if err := inout.Open(); err != nil {
			log.Println(err)
			return 1
		}
		defer inout.Close()
		if err := ExtractCSV(inout.Reader, inout.Writer, fps); err != nil {
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
	flags.StringVar(&input, "i", "", "short for -input")
	flags.StringVar(&output, "output", "", "csv output file or standard output")
	flags.StringVar(&output, "o", "", "short for -output")
	flags.IntVar(&fps, "frames-per-second", 30, "frames per second, usually 24 or 30")
	flags.IntVar(&fps, "fps", 30, "short for -frames-per-second")
	flags.BoolVar(&auto, "auto", false, "automatically choose input files and output names")
	if err := flags.Parse(args); err != nil {
		flags.Usage()
		return 1
	}
	var F, G []string
	if auto {
		var wg sync.WaitGroup
		t0 := time.Now()
		log.Println("This is edtool,", t0.Format(time.ANSIC))
		log.Println("    * extracting information from edl files")
		F, G = shellutil.GetInputOutput(input, output, isEDLFile, ".csv")
		status := 0
		for i, f := range F {
			wg.Add(1)
			go func(f, g string) {
				defer wg.Done()
				log.Println("        ", g)
				status |= extractCSV(f, g, fps)
			}(f, G[i])
		}
		wg.Wait()
		t1 := time.Now()
		log.Printf("    * terminated at %s, total duration %s\n", t1.Format(time.ANSIC), t1.Sub(t0))
		log.Println("done")
		return status
	} else {
		return extractCSV(input, output, fps)
	}
}
