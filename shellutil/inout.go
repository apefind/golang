package shellutil

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

type InOut struct {
	bufio.ReadWriter
	In, Out string
	in, out *os.File
}

func NewInOut(input, output string) *InOut {
	return &InOut{In: input, Out: output}
}

func (inout *InOut) Open() error {
	var err error
	if inout.In == "" {
		inout.Reader = bufio.NewReader(os.Stdin)
	} else {
		inout.in, err = os.Open(inout.In)
		if err != nil {
			return err
		}
		inout.Reader = bufio.NewReader(inout.in)
	}
	if inout.Out == "" {
		inout.Writer = bufio.NewWriter(os.Stdout)
	} else {
		inout.out, err = os.Create(inout.Out)
		if err != nil {
			return err
		}
		inout.Writer = bufio.NewWriter(inout.out)
	}
	return nil
}

func (inout *InOut) Flush() {
	inout.Writer.Flush()
}

func (inout *InOut) Close() {
	if inout.in != nil {
		inout.in.Close()
	}
	inout.Flush()
	if inout.out != nil {
		inout.out.Close()
	}
}

func (inout *InOut) SetAutoOutput(ext string) {
	if ext == "" {
		ext = filepath.Ext(inout.In)
	}
	if inout.Out == "" {
		inout.Out = strings.TrimSuffix(inout.In, filepath.Ext(inout.In)) + ext
	} else if IsDir(inout.Out) {
		inout.Out = inout.Out + PathSeparator + GetFileBasename(inout.In) + ext
	} else {
		inout.Out = inout.Out + ext
	}
}

func (inout *InOut) CreateOutputDirectories(ext string) {
}

func GetMultiInOut(input, output string, isValid ValidFilenameFunc, ext string) []*InOut {
	var F []string
	if input == "" {
		F = GetCurDirFilenames(isValid)
	} else if IsDir(input) {
		F = GetDirFilenames(input, isValid)
	} else {
		F = []string{input}
	}
	var I []*InOut
	for _, f := range F {
		I = append(I, NewInOut(f, GetOutputFilename(f, output, ext)))
	}
	return I
}
