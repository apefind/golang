package edl

import (
	"fmt"

	"github.com/tealeg/xlsx"
)

func (xls *XLSX) RemoveSheet(name string) error {
	sheet, ok := xls.Sheet[name]
	if !ok {
		return fmt.Errorf("sheet <%s> not found", name)
	}
	i := xls.GetSheetIndex(sheet)
	xls.Sheets = append(xls.Sheets[:i], xls.Sheets[i+1:]...)
	delete(xls.Sheet, name)
	return nil
}

func CopySheet(source, target *xlsx.Sheet) {
	CopyXLSXArea(source, target, 0, 0, source.MaxRow, source.MaxCol, 0, 0)
}

func CopyXLSXArea(source, target *xlsx.Sheet, row0, col0, row1, col1, rowOffset, colOffset int) {
	for row := row0; row < row1; row++ {
		for col := col0; col < col1; col++ {
			CopyXLSXCell(source.Cell(row, col), target.Cell(row+rowOffset, col+colOffset))
		}
	}
}

func CopyXLSXRow(sheet0, sheet1 *xlsx.Sheet, row0, row1 int) {
	for col := 0; col < sheet0.MaxCol; col++ {
		CopyXLSXCell(sheet0.Cell(row0, col), sheet1.Cell(row1, col))
	}
}

func CopyXLSXCell(source, target *xlsx.Cell) {
	switch source.Type() {
	case xlsx.CellTypeString:
		target.SetString(source.Value)
	case xlsx.CellTypeFormula:
		target.SetFormula(source.Formula())
	case xlsx.CellTypeNumeric:
		if v, err := source.Float(); err == nil {
			target.SetFloatWithFormat(v, source.NumFmt)
		} else {
			target.SetString(source.Value)
		}
	default:
		target.SetValue(source.Value)
	}
	target.SetStyle(source.GetStyle())
	target.NumFmt = source.NumFmt
}

func (xls *XLSX) GetSheetIndex(sheet *xlsx.Sheet) int {
	for i, s := range xls.Sheets {
		if s == sheet {
			return i
		}
	}
	return -1
}

func UniqueStrings(s []string) []string {
	m := make(map[string]struct{})
	for _, t := range s {
		m[t] = struct{}{}
	}
	s = make([]string, 0, len(s))
	for t, _ := range m {
		s = append(s, t)
	}
	return s
}
