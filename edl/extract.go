package edl

import (
	"bufio"
	"encoding/csv"
	"strconv"
)

var CSVHeader = []string{"Event No", "Reel", "Track Type", "Edit Type", "Transition",
	"Source In", "Source Out", "Prog In H", "Prog In M", "Prog In S", "Prog In F",
	"Prog Out H", "Prog Out M", "Prog Out S", "Prog Out F", "Frames In", "Frames Out",
	"Elapsed Frames", "Seconds", "Frames", "Notes"}

func ExtractCSV(r *bufio.Reader, w *bufio.Writer, fps int) error {
	var err error
	writer := csv.NewWriter(w)
	if err = writer.Write(CSVHeader); err != nil {
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
	var record [21]string
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
