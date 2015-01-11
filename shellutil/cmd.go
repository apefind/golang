package shellutil

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

func Usage(cmd string, flags *flag.FlagSet) {
	if cmd == "prompt" {
		fmt.Fprintf(os.Stderr, "\n%s %s -l <length> -r <ratio>\n\n", filepath.Base(os.Args[0]), cmd)
		fmt.Fprintf(os.Stderr, "\tPrint a length limited prompt\n\n")
	} else if cmd == "timeout" {
		fmt.Fprintf(os.Stderr, "\n%s %s -d <duration> command [args ...]\n\n",
			filepath.Base(os.Args[0]), cmd)
		fmt.Fprintf(os.Stderr, "\tRun a command under time limitation\n\n")
	} else if cmd == "batch" {
		fmt.Fprintf(os.Stderr, "\n%s %s -n <jobs> -d <duration> command [args ...]\n\n",
			filepath.Base(os.Args[0]), cmd)
		fmt.Fprintf(os.Stderr, "\tRun a command under time limitation\n\n")
	} else {
		fmt.Fprintf(os.Stderr, "\n%s %s\n\n", filepath.Base(os.Args[0]), cmd)
		fmt.Fprintf(os.Stderr, "\tRun some shell utility\n\n")
	}
	flags.PrintDefaults()
}

func PromptCmd(args []string) int {
	var length int
	var ratio float64
	flags := flag.NewFlagSet("prompt", flag.ExitOnError)
	flags.IntVar(&length, "length", 32, "maximum length of the prompt")
	flags.IntVar(&length, "l", 32, "")
	flags.Float64Var(&ratio, "ratio", 0.75, "ratio between head/tail of the path")
	flags.Float64Var(&ratio, "r", 0.75, "")
	flags.Usage = func() { Usage("prompt", flags) }
	flags.Parse(args)
	fmt.Print(GetShellPrompt(length, ratio))
	return 0
}

func TimeoutCmd(args []string) int {
	var duration time.Duration
	flags := flag.NewFlagSet("timeout", flag.ExitOnError)
	flags.DurationVar(&duration, "duration", 0*time.Second, "kill the command after given duration")
	flags.DurationVar(&duration, "d", 0*time.Second, "")
	flags.Usage = func() { Usage("timeout", flags) }
	flags.Parse(args)
	cmd := flags.Arg(0)
	if cmd == "" {
		flags.Usage()
		fmt.Fprintf(os.Stderr, "\nno command specified\n\n")
		return int(CommandNotFound)
	}
	var status CommandStatus
	var err error
	w := bufio.NewWriter(os.Stdout)
	if flags.NArg() == 0 {
		status, err = TimeOut(w, duration, cmd)
	} else {
		status, err = TimeOut(w, duration, cmd, flags.Args()[1:]...)
	}
	if err != nil {
		log.Println(err)
	}
	return int(status)
}

func BatchCmd(args []string) int {
	//Batch()
	return 0
}
