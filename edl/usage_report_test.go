package edl

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/apefind/golang/shellutil"
)

var UsageReportFiles = []TestFile{
	{"testdata" + string(filepath.Separator) + "test_usage_report_01.xlsx", 24},
	{"testdata" + string(filepath.Separator) + "test_usage_report_03.xlsx", 24},
	{"testdata" + string(filepath.Separator) + "test_usage_report_04.xlsx", 30},
}

func createUsageReport(input, output string, fps int) error {
	edl, err := OpenXLSX(input)
	if err != nil {
		return err
	}
	report := NewXLSX()
	_, err = report.AddEDLSheet(XLSXMainSheetName)
	if err != nil {
		return err
	}
	CopySheet(edl.Sheet[XLSXMainSheetName], report.Sheet[XLSXMainSheetName])
	report.Filename = output
	if err := report.CreateUsageReport(fps); err != nil {
		return err
	}
	if err := report.Save(report.Filename); err != nil {
		return err
	}
	return nil
}

func updateUsageReportReferenceData(output string) {
	for _, f := range getTestFiles(UsageReportFiles) {
		output := getOutputFilename(f.filename, output, "xlsx", f.fps)
		if !shellutil.IsFile(output) {
			createUsageReport(f.filename, output, f.fps)
		}
	}
}

func TestUsageReport(t *testing.T) {
	updateUsageReportReferenceData("testdata")
	tmpDir, _ := ioutil.TempDir("/tmp", "test_edl_usage_report_")
	for _, f := range getTestFiles(UsageReportFiles) {
		t.Logf("checking %s\n", f.filename)
		tmpXLSX := getOutputFilename(f.filename, tmpDir, "xlsx", f.fps)
		if err := createUsageReport(f.filename, tmpXLSX, f.fps); err != nil {
			t.Error(err)
		}
		refXLSX := getOutputFilename(f.filename, "testdata", "xlsx", f.fps)
		if err := compareXLSXFiles(refXLSX, tmpXLSX); err != nil {
			t.Errorf("%s: %s", tmpXLSX, err)
		}
	}
}
