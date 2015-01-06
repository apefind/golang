package shellutil

import (
	"bufio"
	"io/ioutil"
	"testing"
	"time"
)

func TestTimeout(t *testing.T) {
	var cmd string = "sleep"
	var status CommandStatus
	var err error
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}
	w := bufio.NewWriter(ioutil.Discard)
	status, err = TimeOut(w, 1*time.Second, "does not exist")
	if status != CommandNotFound || err == nil {
		t.Error("expected ", CommandNotFound, " got ", status, err)
	}
	status, err = TimeOut(w, 3*time.Second, cmd, "1")
	if status != Ok || err != nil {
		t.Error("expected ", Ok, " got ", status, err)
	}
	status, err = TimeOut(w, 2*time.Second, cmd, "3")
	if status != ProcessKilled || err != nil {
		t.Error("expected ", ProcessKilled, " got ", status, err)
	}
}
