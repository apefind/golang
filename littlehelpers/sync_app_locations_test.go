package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/apefind/golang/shellutil"
)

func ExampleSyncAppLocations() {
	applications := []string{
		"Atom.app",
		"LibreOffice.app",
		"Firefox.app",
		"FileZilla.app",
		"iTerm.app",
		"Unknown.app",
		"Network/Firefox.app",
		"Network/FileZilla.app",
		"Programming/Atom.app",
		"Tools/iTerm.app",
	}
	tmpdir, _ := ioutil.TempDir("/tmp", "test_sync_app_locations_")
	for _, path := range applications {
		os.MkdirAll(tmpdir+string(filepath.Separator)+path, 0755)
	}
	err := SyncAppLocations(tmpdir, "")
	if err != nil {
		fmt.Println(err)
	}
	walk := func(path string) error {
		fmt.Println(strings.TrimSuffix(strings.TrimPrefix(path, tmpdir+"/"), "/"))
		return nil
	}
	filepath.Walk(tmpdir, shellutil.GetDirWalkFunc(tmpdir, strings.Fields("*.app"), strings.Fields(""), true, walk))
	// Output:
	// Network/FileZilla.app
	// Network/Firefox.app
	// Office/LibreOffice.app
	// Programming/Atom.app
	// Tools/iTerm.app
	// Unknown.app
}
