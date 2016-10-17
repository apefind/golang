package edl

import (
	"fmt"
	"strings"

	"github.com/tealeg/xlsx"
)

func DiffUsageReports(filename0, filename1 string, output string, headers string) error {
	xls0, err := OpenXLSX(filename0)
	if err != nil {
		return err
	}
	xls1, err := OpenXLSX(filename1)
	if err != nil {
		return err
	}
	xls := NewXLSX()
	if err := Compare(xls0, xls1, xls, headers); err != nil {
		return err
	}
	if err := xls.Save(output); err != nil {
		return err
	}
	return nil
}

func Compare(xls0, xls1 *XLSX, result *XLSX, headers string) error {
	sheet0, ok := xls0.Sheet[XLSXMainSheetName]
	if !ok {
		return fmt.Errorf("no <%s> sheet found", XLSXMainSheetName)
	}
	sheet1, ok := xls1.Sheet[XLSXMainSheetName]
	if !ok {
		return fmt.Errorf("no <%s> sheet found", XLSXMainSheetName)
	}
	resultSheet0, err := result.AddEDLSheet("EDL1")
	if err != nil {
		return err
	}
	CopySheet(sheet0, resultSheet0)
	resultSheet1, err := result.AddEDLSheet("EDL2")
	if err != nil {
		return err
	}
	CopySheet(sheet1, resultSheet1)
	return CompareSheets(sheet0, sheet1, resultSheet0, resultSheet1, headers)
}

func CompareSheets(sheet0, sheet1 *xlsx.Sheet, resultSheet0, resultSheet1 *xlsx.Sheet, headers string) error {
	columns, err := getHashColumns(headers)
	if err != nil {
		return err
	}
	m0, m1 := getSheetMap(sheet0, columns), getSheetMap(sheet1, columns)
	for _, row := range resultSheet0.Rows {
		_, ok := m1[getRowHash(row, columns)]
		if !ok {
			for _, cell := range row.Cells {
				cell.SetStyle(XLSXDiffStyle)
			}
		}
	}
	for _, row := range resultSheet1.Rows {
		_, ok := m0[getRowHash(row, columns)]
		if !ok {
			for _, cell := range row.Cells {
				cell.SetStyle(XLSXDiffStyle)
			}
		}
	}
	return nil
}

func getSheetMap(sheet *xlsx.Sheet, columns []int) map[string]struct{} {
	m := make(map[string]struct{})
	for _, row := range sheet.Rows {
		m[getRowHash(row, columns)] = struct{}{}
	}
	return m
}

func getRowHash(row *xlsx.Row, columns []int) string {
	values := make([]string, 0, len(columns))
	for _, col := range columns {
		values = append(values, row.Cells[col].Value)
	}
	return strings.Join(values, "/")
}

func getHashColumns(headers string) ([]int, error) {

	findColumn := func(header string) (int, error) {
		for i, column := range XLSXColLayout {
			if column.Header == header {
				return i, nil
			}
		}
		return -1, fmt.Errorf("column %s not found", header)
	}

	columns := make([]int, 0)
	for _, header := range strings.Split(headers, ",") {
		col, err := findColumn(header)
		if err != nil {
			return columns, err
		}
		columns = append(columns, col)
	}
	return columns, nil
}
