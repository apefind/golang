package shellutil

import (
	"bufio"
	"io"
	"os/exec"
	"time"
)

//go:generate stringer -type=CommandStatus
type CommandStatus int

const (
	Ok CommandStatus = iota
	UnspecifiedError
	ProcessKilled
	CommandNotFound
	PipeError
)

// Timeout executes a system command given a maximum lifetime, cmd and args are the arguments
// passed to exec.Command(), stdout will be written to w, flushing after each eol
func TimeOut(w *bufio.Writer, timeout time.Duration, cmd string, args ...string) (CommandStatus, error) {
	command := exec.Command(cmd, args...)
	command.Stderr = command.Stdout
	stdout, err := command.StdoutPipe()
	if err != nil {
		return PipeError, err
	}
	if err := command.Start(); err != nil {
		return CommandNotFound, err
	}
	go func() {
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
	}()
	done := make(chan error, 1)
	go func() {
		done <- command.Wait()
	}()
	select {
	case <-time.After(timeout):
		if err := command.Process.Kill(); err != nil {
			return UnspecifiedError, err
		}
		close(done)
		return ProcessKilled, nil
	case err := <-done:
		if err != nil {
			return UnspecifiedError, err
		}
	}
	return Ok, nil
}
