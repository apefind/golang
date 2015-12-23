package shellutil

import (
	"log"
	"os"
	"path/filepath"
	"regexp"
)

// GetCleanUpWalkFunc returns a file/directory removal function based on glob style matching and
// regular expressions, suitable as input for filepath.Walk
func GetCleanUpWalkFunc(dir string, glob []string, re []string, recurse bool, rm func(string) error) filepath.WalkFunc {
	var regExps []*regexp.Regexp
	for _, expr := range re {
		if regExp, err := regexp.Compile(expr); err != nil {
			regExps = append(regExps, regExp)
		}
	}
	walk := func(path string, f os.FileInfo, err error) error {
		if !recurse && f.IsDir() && path != dir {
			return filepath.SkipDir
		}
		for _, pattern := range glob {
			match, err := filepath.Match(pattern, f.Name())
			if err != nil {
				return err
			}
			if match {
				if err := rm(path); err != nil {
					return err
				}
				return nil
			}
		}
		for _, regExp := range regExps {
			if regExp.MatchString(f.Name()) {
				if err := rm(path); err != nil {
					return err
				}
				return nil
			}
		}
		return nil
	}
	return walk
}

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
	return filepath.Walk(dir, GetCleanUpWalkFunc(dir, glob, regex, recurse, rm))
}
