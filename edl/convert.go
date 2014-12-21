package edl

import (
	"bufio"
	"encoding/csv"
	"strconv"
)

func ConvertToCSV(r *bufio.Reader, w *bufio.Writer, fps int) error {
	var err error
	writer := csv.NewWriter(w)
	if err = writer.Write(CSV_HEADER); err != nil {
		return err
	}
	seconds, frames := 0, 0
	for _, entry := range Parse(r, fps) {
		if err = writer.Write(entry.CSV()); err != nil {
			return err
		}
		w.Flush()
		seconds, frames = seconds+entry.Seconds, frames+entry.Frames
	}
	var record [20]string
	record[18], record[19] = strconv.Itoa(seconds), strconv.Itoa(frames)
	if err = writer.Write(record[:]); err != nil {
		return err
	}
	w.Flush()
	record[18], record[19] = strconv.Itoa(seconds+int(frames/fps)),
		strconv.Itoa(frames-int(frames/fps)*fps)
	if err = writer.Write(record[:]); err != nil {
		return err
	}
	w.Flush()
	return nil
}
