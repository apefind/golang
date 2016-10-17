package edl

import (
	"fmt"

	"github.com/tealeg/xlsx"
)

var XLSXHeaderStyle *xlsx.Style
var XLSXDiffStyle *xlsx.Style
var XLSXColLayout []struct {
	Header        string
	CellType      xlsx.CellType
	Width         float64
	NumericFormat string
}
var XLSXMainSheetName = "EDL"
var XLSXSecondsColumn = 18
var XLSXFramesColumn = 19
var XLSXSourceColumn = 20
var XLSXSecondsTotalFormula = "=SUM(S2:S%d)"           // e.g.	=SUM(S2:S8)
var XLSXFramesTotalFormula = "=SUM(T2:T%d)"            // 		=SUM(T2:T8)
var XLSXSecondsRemainderFormula = "=S%d+INT(T%d/%d)"   // 		=S9+INT(T9/30)
var XLSXFramesRemainderFormula = "=T%d-INT(T%d/%d)*%d" // 		=T9-INT(T9/30)*30

type XLSX struct {
	xlsx.File
	Filename string
}

func init() {
	XLSXHeaderStyle = xlsx.NewStyle()
	XLSXHeaderStyle.Font.Bold = true
	XLSXDiffStyle = xlsx.NewStyle()
	XLSXDiffStyle.ApplyFill = true
	XLSXDiffStyle.Fill = *xlsx.NewFill("solid", "CCF2FF", "00000000")
	XLSXColLayout = []struct {
		Header        string
		CellType      xlsx.CellType
		Width         float64
		NumericFormat string
	}{
		{"Event No", xlsx.CellTypeString, xlsx.ColWidth, ""},
		{"Reel", xlsx.CellTypeString, 30.0, ""},
		{"Track Type", xlsx.CellTypeString, 3.5, ""},
		{"Edit Type", xlsx.CellTypeString, 3.5, ""},
		{"Transition", xlsx.CellTypeString, xlsx.ColWidth, ""},
		{"Source In", xlsx.CellTypeString, 15.0, "00"},
		{"Source Out", xlsx.CellTypeString, 15.0, "00"},
		{"Prog In H", xlsx.CellTypeNumeric, 3.5, "00"},
		{"Prog In M", xlsx.CellTypeNumeric, 3.5, "00"},
		{"Prog In S", xlsx.CellTypeNumeric, 3.5, "00"},
		{"Prog In F", xlsx.CellTypeNumeric, 3.5, "00"},
		{"Prog Out H", xlsx.CellTypeNumeric, 3.5, "00"},
		{"Prog Out M", xlsx.CellTypeNumeric, 3.5, "00"},
		{"Prog Out S", xlsx.CellTypeNumeric, 3.5, "00"},
		{"Prog Out F", xlsx.CellTypeNumeric, 3.5, "00"},
		{"Frames In", xlsx.CellTypeNumeric, xlsx.ColWidth, ""},
		{"Frames Out", xlsx.CellTypeNumeric, xlsx.ColWidth, ""},
		{"Elapsed Frames", xlsx.CellTypeNumeric, xlsx.ColWidth, ""},
		{"Seconds", xlsx.CellTypeNumeric, xlsx.ColWidth, ""},
		{"Frames", xlsx.CellTypeNumeric, xlsx.ColWidth, ""},
		{"Source", xlsx.CellTypeNumeric, xlsx.ColWidth, ""},
		{"Notes", xlsx.CellTypeString, 75.0, ""},
	}
}

func NewXLSX() *XLSX {
	return &XLSX{*xlsx.NewFile(), ""}
}

func OpenXLSX(filename string) (*XLSX, error) {
	xls_, err := xlsx.OpenFile(filename)
	if err != nil {
		return nil, err
	}
	xls := &XLSX{*xls_, filename}
	_, ok := xls.Sheet[XLSXMainSheetName]
	if !ok {
		return xls, fmt.Errorf("no <%s> sheet found", XLSXMainSheetName)
	}
	return xls, nil
}

func (xls *XLSX) AddEDLSheet(name string) (*xlsx.Sheet, error) {
	sheet, err := xls.AddSheet(name)
	if err != nil {
		return sheet, err
	}
	for i, layout := range XLSXColLayout {
		sheet.SetColWidth(i, i, layout.Width)
	}
	row := sheet.AddRow()
	for _, layout := range XLSXColLayout {
		cell := row.AddCell()
		cell.Value = layout.Header
		cell.SetStyle(XLSXHeaderStyle)
	}
	return sheet, nil
}
