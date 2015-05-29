package shellutil

import (
	"bufio"
	"io"
	"os/exec"
	"time"
)

func TimeIt(w *bufio.Writer, n int, cmd string, args ...string) (time.Duration, error) {
	start := time.Now()
	for i := 0; i < n; i++ {
		command := exec.Command(cmd, args...)
		command.Stderr = command.Stdout
		stdout, err := command.StdoutPipe()
		if err != nil {
			return time.Since(start), err
		}
		if err := command.Start(); err != nil {
			return time.Since(start), err
		}
		r := bufio.NewReader(stdout)
		for {
			line, err := r.ReadString('\n')
			if err == io.EOF {
				break
			} else if err != nil {
				panic(err)
			}
			w.WriteString(line)
			w.Flush()
		}
		command.Wait()
	}
	return time.Since(start), nil
}
