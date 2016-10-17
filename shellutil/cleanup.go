package shellutil

import (
	"log"
	"os"
	"path/filepath"
)

func CleanUp(dir string, glob []string, regex []string, recurse bool, simulate bool) error {
	rm := func(path string) error {
		log.Println(filepath.Clean(path))
		if !simulate {
			if err := os.RemoveAll(path); err != nil {
				return err
			}
		}
		return nil
	}
	return filepath.Walk(dir, GetDirWalkFunc(dir, glob, regex, recurse, rm))
}
