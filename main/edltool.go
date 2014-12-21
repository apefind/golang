package main

import (
	"apefind/edl"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	var input, output string
	var fps int
	flag.StringVar(&input, "input", "", "input file or directory")
	flag.StringVar(&output, "output", "", "input file or directory")
	flag.IntVar(&fps, "fps", 30, "frames per second")
	flag.Usage = func() {
		fmt.Fprintf(os.Stdout, "\nUsage: %s -input [input] -output [output] -fps [30|24]\n\n",
			filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}
	flag.Parse()
	if err := edl.ConvertFileToCSV(input, output, fps); err != nil {
		fmt.Println(err)
	}
}
