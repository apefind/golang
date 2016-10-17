package edl

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/apefind/golang/shellutil"
)

type TestFile struct {
	filename string
	fps      int
}

var EDLFiles = []TestFile{
	{"testdata" + string(filepath.Separator) + "test01.edl", 24},
	{"testdata" + string(filepath.Separator) + "test01.edl", 30},
	{"testdata" + string(filepath.Separator) + "test03.txt", 24},
	{"testdata" + string(filepath.Separator) + "test04.txt", 30},
	{"testdata" + string(filepath.Separator) + "test05.edl", 30},
	{"testdata" + string(filepath.Separator) + "test06.edl", 30},
	{"testdata" + string(filepath.Separator) + "test07.txt", 30},
}

func getTestFiles(files []TestFile) []TestFile {
	testfiles := make([]TestFile, 0, len(files))
	for _, f := range files {
		if shellutil.IsFile(f.filename) {
			testfiles = append(testfiles, f)
		}
	}
	return testfiles
}

func getOutputFilename(input, output, format string, fps int) string {
	filename := strings.TrimSuffix(input, filepath.Ext(input)) + "_" + strconv.Itoa(fps) + "fps"
	return shellutil.GetOutputFilename(filename, output, "."+format)
}

func updateEDLReferenceData(output, format string) {
	for _, f := range getTestFiles(EDLFiles) {
		output := getOutputFilename(f.filename, output, format, f.fps)
		if !shellutil.IsFile(output) {
			convertEDL(f.filename, output, format, f.fps)
		}
	}
}

func compareCSVFiles(csv0, csv1 string) error {
	buf0, err := ioutil.ReadFile(csv0)
	if err != nil {
		return err
	}
	buf1, err := ioutil.ReadFile(csv1)
	if err != nil {
		return err
	}
	S, T := strings.Split(string(buf0), "\n"), strings.Split(string(buf1), "\n")
	if len(S) != len(T) {
		return fmt.Errorf("expected %d number of lines, got %d", len(S), len(T))
	}
	for i := 0; i < len(T); i++ {
		if S[i] != T[i] {
			return fmt.Errorf("expected %s, got %s", S[i], T[i])
		}
	}
	return nil
}

func compareXLSXFiles(xlsx0, xlsx1 string) error {
	xls0, err := OpenXLSX(xlsx0)
	if err != nil {
		return err
	}
	xls1, err := OpenXLSX(xlsx1)
	if err != nil {
		return err
	}
	if len(xls0.Sheet) != len(xls1.Sheet) {
		return fmt.Errorf("expected %d sheets, got %d", len(xls0.Sheet), len(xls1.Sheet))
	}
	for _, sheet0 := range xls0.Sheets {
		sheet1, ok := xls1.Sheet[sheet0.Name]
		if !ok {
			return fmt.Errorf("sheet %s not found", sheet0.Name)
		}
		if sheet0.MaxRow != sheet1.MaxRow {
			return fmt.Errorf("expected %d rows, got %d", sheet0.MaxRow, sheet1.MaxRow)
		}
		if sheet0.MaxCol != sheet1.MaxCol {
			return fmt.Errorf("expected %d columns, got %d", sheet0.MaxCol, sheet1.MaxCol)
		}
		for row := 0; row < sheet0.MaxRow; row++ {
			for col := 0; col < sheet0.MaxCol; col++ {
				cell0, cell1 := sheet0.Cell(row, col), sheet1.Cell(row, col)
				if cell0.Value != cell1.Value {
					return fmt.Errorf("expected %s at %d/%d, got %s", cell0.Value, row, col, cell1.Value)
				}
			}
		}
	}
	return nil
}

func convertEDL(input, output, format string, fps int) error {
	r, err := os.Open(input)
	if err != nil {
		return err
	}
	defer r.Close()
	w, err := os.Create(output)
	if err != nil {
		return err
	}
	defer w.Close()
	if format == "csv" {
		return ConvertToCSV(bufio.NewReader(r), bufio.NewWriter(w), fps)
	} else if format == "xlsx" {
		xls := NewXLSX()
		_, err := xls.AddEDLSheet(XLSXMainSheetName)
		if err != nil {
			return err
		}
		if err := xls.ConvertToXLSX(bufio.NewReader(r), fps); err != nil {
			return err
		}
		return xls.Write(bufio.NewWriter(w))
	}
	return nil
}

func TestConvertEDLCSV(t *testing.T) {
	updateEDLReferenceData("testdata", "csv")
	tmpDir, _ := ioutil.TempDir("/tmp", "test_edl_csv_")
	for _, f := range getTestFiles(EDLFiles) {
		t.Logf("checking %s\n", f.filename)
		tmpCSV := getOutputFilename(f.filename, tmpDir, "csv", f.fps)
		if err := convertEDL(f.filename, tmpCSV, "csv", f.fps); err != nil {
			t.Error(err)
		}
		refCSV := getOutputFilename(f.filename, "testdata", "csv", f.fps)
		if err := compareCSVFiles(refCSV, tmpCSV); err != nil {
			t.Errorf("%s: %s", tmpCSV, err)
		}
	}
}

func TestConvertEDLXLSX(t *testing.T) {
	updateEDLReferenceData("testdata", "xlsx")
	tmpDir, _ := ioutil.TempDir("/tmp", "test_edl_xlsx_")
	for _, f := range getTestFiles(EDLFiles) {
		t.Logf("checking %s\n", f.filename)
		tmpXLSX := getOutputFilename(f.filename, tmpDir, "xlsx", f.fps)
		if err := convertEDL(f.filename, tmpXLSX, "xlsx", f.fps); err != nil {
			t.Error(err)
		}
		refXLSX := getOutputFilename(f.filename, "testdata", "xlsx", f.fps)
		if err := compareXLSXFiles(refXLSX, tmpXLSX); err != nil {
			t.Errorf("%s: %s", tmpXLSX, err)
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

func ExampleConvertToCSV() {
	ConvertToCSV(bufio.NewReader(strings.NewReader(exampleEDL)), bufio.NewWriter(os.Stdout), 24)
	// Output:
	// Event No,Reel,Track Type,Edit Type,Transition,Source In,Source Out,Prog In H,Prog In M,Prog In S,Prog In F,Prog Out H,Prog Out M,Prog Out S,Prog Out F,Frames In,Frames Out,Elapsed Frames,Seconds,Frames,Source,Notes
	// 001,BL,V,C,,00:00:00:00,00:00:02:00,00,00,00,00,00,00,02,00,0,48,48,2,0,,
	// 002,SF003,V,C,,00:59:59:28,01:00:00:07,00,00,02,00,00,00,02,07,48,55,7,0,7,,FROM CLIP NAME: GOW01.jpg
	// 003,BL,V,C,,00:00:00:00,00:00:01:15,00,00,02,07,00,00,03,22,55,94,39,1,15,,
	// 004,SF003,V,C,,00:59:59:28,01:00:00:07,00,00,03,22,00,00,04,04,94,100,6,0,6,,FROM CLIP NAME: GOW02.jpg
	// 005,BL,V,C,,00:00:00:00,00:00:02:19,00,00,04,04,00,00,06,23,100,167,67,2,19,,
	// 006,SF003,V,C,,00:59:59:28,01:00:00:15,00,00,06,23,00,00,07,12,167,180,13,0,13,,FROM CLIP NAME: GOW03.jpg
	// 007,BL,V,C,,00:00:00:00,00:00:01:15,00,00,07,12,00,00,09,02,180,218,38,1,14,,
	// 008,SF003,V,C,,00:59:59:28,01:00:01:25,00,00,09,02,00,00,10,24,218,264,46,1,22,,FROM CLIP NAME: GOW04.jpg
}
