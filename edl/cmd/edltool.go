package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/apefind/golang/edl"
	"github.com/apefind/golang/shellutil"
)

func usage(cmd string, flags *flag.FlagSet) {
	if cmd == "convert" {
		fmt.Fprintf(os.Stderr, "\n%s %s -fps [24|30] -format [xlsx|csv] -output <file or directory> ...\n\n",
			filepath.Base(os.Args[0]), cmd)
		fmt.Fprintf(os.Stderr, "\tConvert edit decision into xlsx or csv\n\n")
	} else if cmd == "extract" {
		fmt.Fprintf(os.Stderr, "\n%s %s -output <file or directory> ...\n\n", filepath.Base(os.Args[0]), cmd)
		fmt.Fprintf(os.Stderr, "\tExtract usage information into xlsx or csv\n\n")
	} else {
		fmt.Fprintf(os.Stderr, "%s %s\n", filepath.Base(os.Args[0]), cmd)
	}
	flags.PrintDefaults()
}

func ConvertCmd(args []string) int {
	var output, format string
	var fps int
	flags := flag.NewFlagSet("convert", flag.ExitOnError)
	flags.Usage = func() { usage("convert", flags) }
	flags.StringVar(&output, "output", "", "csv output file or standard output")
	flags.StringVar(&output, "o", "", "short for -output")
	flags.StringVar(&format, "format", "xlsx", "csv od xlsx")
	flags.IntVar(&fps, "fps", 30, "frames per second")
	if err := flags.Parse(args); err != nil {
		flags.Usage()
		return 1
	}

	getInOut := func() *shellutil.InOut {
		inout := shellutil.NewInOut(flags.Args()...)
		if output != "" {
			inout.SetOutput(output, "."+format)
		}
		return inout
	}

	extractCSV := func(inout *shellutil.InOut) error {
		if err := inout.Open(); err != nil {
			return err
		}
		defer inout.Close()
		if err := edl.ConvertToCSV(inout.Reader, inout.Writer, fps); err != nil {
			return err
		}
		return nil
	}

	extractXLSX := func(inout *shellutil.InOut) error {
		if err := inout.Open(); err != nil {
			return err
		}
		defer inout.Close()
		xls := edl.NewXLSX()
		_, err := xls.AddEDLSheet(edl.XLSXMainSheetName)
		if err != nil {
			return err
		}
		if err := xls.ConvertToXLSX(inout.Reader, fps); err != nil {
			return err
		}
		if err := xls.Write(inout.Writer); err != nil {
			return err
		}
		return nil
	}

	if output == "" {
		log.SetOutput(ioutil.Discard)
	}
	log.Println("    * converting edl(s) to " + format)
	inout := getInOut()
	for _, input := range inout.Input {
		log.Printf("        %s", input)
	}
	log.Println("    * creating result file")
	if format == "xlsx" {
		if err := extractXLSX(inout); err != nil {
			log.Println("        ", err)
			return 1
		}
	} else {
		if err := extractCSV(inout); err != nil {
			log.Println("        ", err)
			return 1
		}
	}
	log.Printf("        %s", inout.Output)
	return 0
}

func UsageReportCmd(args []string) int {
	var input, output string
	var fps int
	flags := flag.NewFlagSet("report", flag.ExitOnError)
	flags.Usage = func() { usage("report", flags) }
	flags.StringVar(&input, "input", "", "converted edl")
	flags.StringVar(&input, "i", "", "short for -input")
	flags.StringVar(&output, "output", "", "usage report")
	flags.StringVar(&output, "o", "", "short for -output")
	flags.IntVar(&fps, "fps", 30, "frames per second")

	if err := flags.Parse(args); err != nil {
		flags.Usage()
		return 1
	}

	getUsageReport := func(input, output string) (*edl.XLSX, error) {
		xls, err := edl.OpenXLSX(input)
		if err != nil || output == "" {
			return xls, err
		}
		report, err := edl.NewUsageReport()
		if err != nil {
			return report, err
		}
		edl.CopySheet(xls.Sheet[edl.XLSXMainSheetName], report.Sheet[edl.XLSXMainSheetName])
		report.Filename = shellutil.GetOutputFilename(strings.TrimSuffix(input, filepath.Ext(input))+" Usage Report",
			output, ".xlsx")
		return report, nil
	}

	log.Println("    * reading converted edl")
	xls, err := getUsageReport(input, output)
	if err != nil {
		log.Println("        ", err)
		return 1
	}
	log.Println("        ", input)
	log.Println("    * creating usage report")
	if err := xls.CreateUsageReport(fps); err != nil {
		log.Println("        ", err)
		return 1
	}
	log.Println("        ", xls.Filename)
	if err := xls.Save(xls.Filename); err != nil {
		log.Println("        ", err)
		return 1
	}
	return 0
}

func DiffCmd(args []string) int {
	var input, output, columns string
	flags := flag.NewFlagSet("report", flag.ExitOnError)
	flags.Usage = func() { usage("report", flags) }
	flags.StringVar(&input, "input", "", "usage reports")
	flags.StringVar(&input, "i", "", "short for -input")
	flags.StringVar(&output, "output", "", "usage report")
	flags.StringVar(&output, "o", "", "short for -output")
	flags.StringVar(&columns, "columns", "Reel,Source In,Source Out,Notes", "columns to define identical rows")
	if err := flags.Parse(args); err != nil {
		flags.Usage()
		return 1
	}
	t := strings.Split(input, ":")
	if len(t) != 2 {
		flags.Usage()
		return 1
	}
	log.Println("    * comparing usage reports")
	log.Println("        ", t[0])
	log.Println("        ", t[1])
	log.Println("    * creating diff")
	if err := edl.DiffUsageReports(t[0], t[1], output, columns); err != nil {
		log.Println("        ", err)
		return 1
	}
	log.Println("        ", output)
	return 0
}

func main() {
	log.SetFlags(0)
	flag.Usage = func() { usage("[convert|usage_report|diff]", flag.CommandLine) }
	flag.Parse()
	t0 := time.Now()
	log.Println("This is edtool,", t0.Format(time.ANSIC))
	status := 1
	switch flag.Arg(0) {
	case "convert":
		status = ConvertCmd(os.Args[2:])
	case "usage_report":
		status = UsageReportCmd(os.Args[2:])
	case "diff":
		status = DiffCmd(os.Args[2:])
	default:
		flag.PrintDefaults()
	}
	t1 := time.Now()
	log.Printf("    * terminated at %s, total duration %s\n", t1.Format(time.ANSIC), t1.Sub(t0))
	log.Println("done")
	os.Exit(status)
}
