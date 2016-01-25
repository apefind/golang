package shellutil

import (
	"bufio"
	"io/ioutil"
	"testing"
	"time"
)

func TestTimeout(t *testing.T) {
	var cmd = "sleep"
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}
	w := bufio.NewWriter(ioutil.Discard)
	if err := TimeOut(w, 1*time.Second, "does not exist"); err == nil {
		t.Error("expected <exec: does not exist: executable file not found in $PATH>, got", err)
	}
	if err := TimeOut(w, 3*time.Second, cmd, "1"); err != nil {
		t.Error(err)
	}
	if err := TimeOut(w, 2*time.Second, cmd, "3"); err != ProcessKilled {
		t.Error("expected ", ProcessKilled, " got ", err)
	}
}
