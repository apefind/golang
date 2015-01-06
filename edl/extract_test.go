package edl

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type TestEDL struct {
	input, output string
	fps           int
}

var testEDL []TestEDL = []TestEDL{
	{"test01.edl", "test01_24.csv", 24},
	{"test01.edl", "test01_30.csv", 30},
}

func TestCSVExtract(t *testing.T) {
	for _, data := range testEDL {
		var buffer bytes.Buffer
		t.Logf("checking %s\n", data.input)
		f, err := os.Open("testdata" + string(filepath.Separator) + data.input)
		if err != nil {
			t.Error("cannot open input", data.input, ":", err)
			continue
		}
		defer f.Close()
		ExtractCSV(bufio.NewReader(f), bufio.NewWriter(&buffer), data.fps)
		output, err := ioutil.ReadFile("testdata" + string(filepath.Separator) + data.output)
		if err != nil {
			t.Error("cannot open output", data.output, ":", err)
			continue
		}
		S, T := strings.Split(string(buffer.Bytes()), "\n"), strings.Split(string(output), "\n")
		if len(S) != len(T) {
			t.Error("expected", len(S), "number of lines, got", len(T))
			continue
		}
		for i := 0; i < len(T); i++ {
			if S[i] != T[i] {
				t.Error("expected", S[i], "got", T[i])
				break
			}
		}
	}
}

var exampleEDL = `TITLE: TEST DATA GOW
001  BL       V     C        00:00:00:00 00:00:02:00 00:00:00:00 00:00:02:00
002  SF003    V     C        00:59:59:28 01:00:00:07 00:00:02:00 00:00:02:07
* FROM CLIP NAME: GOW01.jpg
003  BL       V     C        00:00:00:00 00:00:01:15 00:00:02:07 00:00:03:22
004  SF003    V     C        00:59:59:28 01:00:00:07 00:00:03:22 00:00:04:04
* FROM CLIP NAME: GOW02.jpg
005  BL       V     C        00:00:00:00 00:00:02:19 00:00:04:04 00:00:06:23
006  SF003    V     C        00:59:59:28 01:00:00:15 00:00:06:23 00:00:07:12
* FROM CLIP NAME: GOW03.jpg
007  BL       V     C        00:00:00:00 00:00:01:15 00:00:07:12 00:00:09:02
008  SF003    V     C        00:59:59:28 01:00:01:25 00:00:09:02 00:00:10:24
* FROM CLIP NAME: GOW04.jpg`

func ExampleExtractCSV() {
	ExtractCSV(bufio.NewReader(strings.NewReader(exampleEDL)), bufio.NewWriter(os.Stdout), 24)
	// Output:
	// Event No,Reel,Track Type,Edit Type,Transition,Source In,Source Out,Prog In H,Prog In M,Prog In S,Prog In F,Prog Out H,Prog Out M,Prog Out S,Prog Out F,Frames In,Frames Out,Elapsed Frames,Seconds,Frames,Comments
	// 001,BL,V,C,,00:00:00:00,00:00:02:00,00,00,00,00,00,00,02,00,0,48,48,2,0
	// 002,SF003,V,C,,00:59:59:28,01:00:00:07,00,00,02,00,00,00,02,07,48,55,7,0,7,FROM CLIP NAME: GOW01.jpg
	// 003,BL,V,C,,00:00:00:00,00:00:01:15,00,00,02,07,00,00,03,22,55,94,39,1,15
	// 004,SF003,V,C,,00:59:59:28,01:00:00:07,00,00,03,22,00,00,04,04,94,100,6,0,6,FROM CLIP NAME: GOW02.jpg
	// 005,BL,V,C,,00:00:00:00,00:00:02:19,00,00,04,04,00,00,06,23,100,167,67,2,19
	// 006,SF003,V,C,,00:59:59:28,01:00:00:15,00,00,06,23,00,00,07,12,167,180,13,0,13,FROM CLIP NAME: GOW03.jpg
	// 007,BL,V,C,,00:00:00:00,00:00:01:15,00,00,07,12,00,00,09,02,180,218,38,1,14
	// 008,SF003,V,C,,00:59:59:28,01:00:01:25,00,00,09,02,00,00,10,24,218,264,46,1,22,FROM CLIP NAME: GOW04.jpg
	// ,,,,,,,,,,,,,,,,,,7,96
	// ,,,,,,,,,,,,,,,,,,11,0
}
