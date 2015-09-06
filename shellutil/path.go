package shellutil

import (
	"bytes"
	"golang.org/x/text/unicode/norm"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const PathSeparator = string(filepath.Separator)

type ValidFilenameFunc func(string) bool

func IsFile(path string) bool {
	stat, err := os.Stat(path)
	return err == nil && stat.Mode().IsRegular()
}

func IsDirectory(path string) bool {
	stat, err := os.Stat(path)
	return err == nil && stat.IsDir()
}

// EqualFilenames uses unicode normmalization
func EqualFilenames(filename0, filename1 string) bool {
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
	} else if IsDirectory(output) {
		return output + PathSeparator + GetFileBasename(input) + ext
	}
	return output + ext
}

func GetInputOutput(input, output string, isValid ValidFilenameFunc, ext string) ([]string, []string) {
	var F, G []string
	if input == "" {
		F = GetCurDirFilenames(isValid)
	} else if IsDirectory(input) {
		F = GetDirFilenames(input, isValid)
	} else {
		F = []string{input}
	}
	for _, f := range F {
		G = append(G, GetOutputFilename(f, output, ext))
	}
	return F, G
}
