package shellutil

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"golang.org/x/text/unicode/norm"
)

const PathSeparator = string(filepath.Separator)

type ValidFilenameFunc func(string) bool

func IsFile(path string) bool {
	stat, err := os.Stat(path)
	return err == nil && stat.Mode().IsRegular()
}

func IsDir(path string) bool {
	stat, err := os.Stat(path)
	return err == nil && stat.IsDir()
}

func IsSymlink(path string) bool {
	stat, err := os.Lstat(path)
	return err == nil && stat.Mode()&os.ModeSymlink != 0
}

// IdenticalFilenames uses unicode normalization
func IdenticalFilenames(filename0, filename1 string) bool {
	return bytes.Equal(norm.NFC.Bytes([]byte(filename0)), norm.NFC.Bytes([]byte(filename1)))
}

func GetFileBasename(path string) string {
	return strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
}

func GetDirContent(path string) ([]string, []string) {
	F, _ := ioutil.ReadDir(path)
	subdirs, filenames := make([]string, 0, len(F)), make([]string, 0, len(F))
	for _, f := range F {
		if f.IsDir() {
			subdirs = append(subdirs, f.Name())
		} else if f.Mode().IsRegular() {
			filenames = append(filenames, f.Name())
		}
	}
	return subdirs, filenames
}

func GetDirFilenames(path string, isValid ValidFilenameFunc) []string {
	F, _ := ioutil.ReadDir(path)
	filenames := make([]string, 0, len(F))
	for _, f := range F {
		if f.Mode().IsRegular() && isValid(f.Name()) {
			filenames = append(filenames, f.Name())
		}
	}
	return filenames
}

func GetCurDirFilenames(isValid ValidFilenameFunc) []string {
	path, _ := os.Getwd()
	return GetDirFilenames(path, isValid)
}

func GetOutputFilename(input, output, ext string) string {
	if output == "" {
		return strings.TrimSuffix(input, filepath.Ext(input)) + ext
	} else if IsDir(output) {
		return output + PathSeparator + GetFileBasename(input) + ext
	}
	return output + ext
}

func GetInputOutput(input, output string, isValid ValidFilenameFunc, ext string) ([]string, []string) {
	var F, G []string
	if input == "" {
		F = GetCurDirFilenames(isValid)
	} else if IsDir(input) {
		F = GetDirFilenames(input, isValid)
	} else {
		F = []string{input}
	}
	for _, f := range F {
		G = append(G, GetOutputFilename(f, output, ext))
	}
	return F, G
}

// GetDirsFromFlagSetArgs returns directories from left over arguments of flag set
func GetDirsFromFlagSetArgs(flags *flag.FlagSet) ([]string, error) {
	if flags.NArg() == 0 {
		wd, err := os.Getwd()
		return []string{wd}, err
	}
	dirs := make([]string, 0, flags.NArg())
	for _, arg := range flags.Args() {
		if !IsDir(arg) {
			return dirs, fmt.Errorf("%s is not a directory", arg)
		}
		arg, err := filepath.EvalSymlinks(arg)
		if err != nil {
			return dirs, err
		}
		dirs = append(dirs, filepath.Clean(arg))
	}
	return dirs, nil
}

// GetDirWalkFunc returns a file/directory walk function based on glob style matching and
// regular expressions, suitable as input for filepath.Walk
func GetDirWalkFunc(dir string, glob []string, re []string, recurse bool, walk func(string) error) filepath.WalkFunc {
	var regExps []*regexp.Regexp
	for _, expr := range re {
		if regExp, err := regexp.Compile(expr); err != nil {
			regExps = append(regExps, regExp)
		}
	}
	walkFunc := func(path string, f os.FileInfo, err error) error {
		if !recurse && f.IsDir() && path != dir {
			return filepath.SkipDir
		}
		for _, pattern := range glob {
			match, err := filepath.Match(pattern, f.Name())
			if err != nil {
				return err
			}
			if match {
				if err := walk(path); err != nil {
					return err
				}
				return nil
			}
		}
		for _, regExp := range regExps {
			if regExp.MatchString(f.Name()) {
				if err := walk(path); err != nil {
					return err
				}
				return nil
			}
		}
		return nil
	}
	return walkFunc
}

func FindFiles(dir string, glob []string, regex []string, recurse bool) ([]string, error) {
	filenames := []string{}
	collect := func(path string) error {
		filenames = append(filenames, path)
		return nil
	}
	return filenames, filepath.Walk(dir, GetDirWalkFunc(dir, glob, regex, recurse, collect))
}
