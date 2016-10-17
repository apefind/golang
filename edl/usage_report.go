package edl

import (
	"fmt"
	"sort"
	"strings"

	"github.com/tealeg/xlsx"
)

func NewUsageReport() (*XLSX, error) {
	xls := NewXLSX()
	_, err := xls.AddEDLSheet(XLSXMainSheetName)
	return xls, err
}

func (xls *XLSX) CreateUsageReport(fps int) error {
	sources, err := xls.getSources()
	if err != nil {
		return err
	}
	if err := xls.AddSourceSheets(sources); err != nil {
		return err
	}
	if err := xls.CreateSourceReports(); err != nil {
		return err
	}
	if err := xls.CreateSourceSummaries(sources, fps); err != nil {
		return err
	}
	return nil
}

func (xls *XLSX) AddSourceSheets(sources []string) error {
	for _, source := range sources {
		xls.RemoveSheet(source)
		_, err := xls.AddEDLSheet(source)
		if err != nil {
			return err
		}
	}
	return nil
}

func (xls *XLSX) CreateSourceReports() error {
	sheet, ok := xls.Sheet[XLSXMainSheetName]
	if !ok {
		return fmt.Errorf("no <%s> sheet found", XLSXMainSheetName)
	}
	for row := 1; row < sheet.MaxRow; row++ {
		source := xls.getSource(sheet.Cell(row, XLSXSourceColumn))
		if source != "" {
			sourceSheet, ok := xls.Sheet[source]
			if !ok {
				return fmt.Errorf("no <%s> sheet found", source)
			}
			CopyXLSXRow(sheet, sourceSheet, row, sourceSheet.MaxRow)
		}
	}
	return nil
}

func (xls *XLSX) CreateSourceSummaries(sources []string, fps int) error {
	for _, source := range sources {
		sheet, ok := xls.Sheet[source]
		if !ok {
			return fmt.Errorf("no <%s> sheet found", source)
		}
		if err := xls.CreateSourceSummary(sheet, fps); err != nil {
			return err
		}
	}
	return nil
}

func (xls *XLSX) CreateSourceSummary(sheet *xlsx.Sheet, fps int) error {
	n := sheet.MaxRow
	sheet.Cell(n, XLSXSecondsColumn).SetFormula(fmt.Sprintf(XLSXSecondsTotalFormula, n))
	sheet.Cell(n, XLSXFramesColumn).SetFormula(fmt.Sprintf(XLSXFramesTotalFormula, n))
	sheet.Cell(n+1, XLSXSecondsColumn).SetFormula(fmt.Sprintf(XLSXSecondsRemainderFormula, n+1, n+1, fps))
	sheet.Cell(n+1, XLSXFramesColumn).SetFormula(fmt.Sprintf(XLSXFramesRemainderFormula, n+1, n+1, fps, fps))
	return nil
}

func (xls *XLSX) getSources() ([]string, error) {
	sheet, ok := xls.Sheet[XLSXMainSheetName]
	if !ok {
		return nil, fmt.Errorf("no <%s> sheet found", XLSXMainSheetName)
	}
	s := make([]string, 0, sheet.MaxRow)
	for row := 1; row < sheet.MaxRow; row++ {
		source := xls.getSource(sheet.Cell(row, XLSXSourceColumn))
		if source != "" {
			s = append(s, source)
		}
	}
	sources := UniqueStrings(s)
	sort.Strings(sources)
	return sources, nil
}

func (xls *XLSX) getSource(cell *xlsx.Cell) string {
	source := cell.Value
	for _, t := range []string{"\\", "/", "*", "[", "]", ":", "?"} { // not allowed in sheet names
		source = strings.Replace(source, t, " ", -1)
	}
	return strings.ToUpper(strings.TrimSpace(source))
}
