package shellutil

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type InOut struct {
	bufio.ReadWriter
	Input  []string
	Output string
	input  []*os.File
	output *os.File
}

func NewInOut(input ...string) *InOut {
	return &InOut{
		Input:  input,
		Output: "",
	}
}

func (inout *InOut) SetOutput(output, ext string) {
	if ext == "" {
		ext = filepath.Ext(inout.Input[0])
	}
	if output == "" {
		inout.Output = strings.TrimSuffix(inout.Input[0], filepath.Ext(inout.Input[0])) + ext
	} else if IsDir(output) {
		inout.Output = output + PathSeparator + GetFileBasename(inout.Input[0]) + ext
	} else {
		inout.Output = output
	}
}

func (inout *InOut) Open() error {
	if len(inout.Input) == 0 {
		inout.Reader = bufio.NewReader(os.Stdin)
	} else {
		inout.input = make([]*os.File, len(inout.Input))
		var readers []io.Reader
		for _, input := range inout.Input {
			r, err := os.Open(input)
			if err != nil {
				return err
			}
			inout.input = append(inout.input, r)
			readers = append(readers, bufio.NewReader(r))
		}
		inout.Reader = bufio.NewReader(io.MultiReader(readers...))
	}
	if inout.Output == "" {
		inout.Writer = bufio.NewWriter(os.Stdout)
	} else {
		var err error
		inout.output, err = os.Create(inout.Output)
		if err != nil {
			return err
		}
		inout.Writer = bufio.NewWriter(inout.output)
	}
	return nil
}

func (inout *InOut) Flush() {
	inout.Writer.Flush()
}

func (inout *InOut) Close() {
	if inout.input != nil {
		for _, input := range inout.input {
			input.Close()
		}
	}
	inout.Flush()
	if inout.output != nil {
		inout.output.Close()
	}
}

func (inout *InOut) SetAutoOutput(ext string) {
	if ext == "" {
		ext = filepath.Ext(inout.Input[0])
	}
	if inout.Output == "" {
		inout.Output = strings.TrimSuffix(inout.Input[0], filepath.Ext(inout.Input[0])) + ext
	} else if IsDir(inout.Output) {
		inout.Output = inout.Output + PathSeparator + GetFileBasename(inout.Input[0]) + ext
	} else {
		inout.Output = inout.Output + ext
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
		inout := NewInOut(f)
		inout.SetOutput(GetOutputFilename(f, output, ext), "")
		I = append(I, inout)
	}
	return I
}
