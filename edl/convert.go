package edl

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"strconv"

	"github.com/tealeg/xlsx"
)

func ConvertToCSV(r *bufio.Reader, w *bufio.Writer, fps int) error {
	var header []string
	for _, col := range XLSXColLayout {
		header = append(header, col.Header)
	}
	var err error
	writer := csv.NewWriter(w)
	if err = writer.Write(header); err != nil {
		return err
	}
	for _, entry := range Parse(r, fps) {
		if err = writer.Write(entry.CSV()); err != nil {
			return err
		}
		w.Flush()
	}
	return nil
}

func (xls *XLSX) ConvertToXLSX(r *bufio.Reader, fps int) error {
	sheet, ok := xls.Sheet[XLSXMainSheetName]
	if !ok {
		return fmt.Errorf("no <%s> sheet found", XLSXMainSheetName)
	}
	for _, entry := range Parse(r, fps) {
		row := sheet.AddRow()
		for i, value := range entry.CSV() {
			cell := row.AddCell()
			if XLSXColLayout[i].CellType == xlsx.CellTypeString {
				cell.SetString(value)
			} else if XLSXColLayout[i].CellType == xlsx.CellTypeNumeric {
				if v, err := strconv.Atoi(value); err == nil {
					if XLSXColLayout[i].NumericFormat != "" {
						cell.SetFloatWithFormat(float64(v), XLSXColLayout[i].NumericFormat)
					} else {
						cell.SetInt(v)
					}
				} else {
					cell.SetString(value)
				}
			}
		}
	}
	return nil
}
