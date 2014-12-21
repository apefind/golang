package edl

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"math"
	"os"
	"strconv"
	"strings"
)

type Entry struct {
	Event, Reel, TrackType, EditType, Transition string
	SourceIn, SourceOut             string
	RecordIn, RecordOut             string
	Comments                        []string
	TimeIn, TimeOut                 [4]string
	FramesIn, FramesOut			    int
	Elapsed, Seconds, Frames        int
}

func NewEntry(S []string, fps int) *Entry {
	e := &Entry{Comments: make([]string, 0, 10)}
	e.Event, e.Reel, e.TrackType, e.EditType, e.Transition = S[0], S[1], S[2], S[3], S[4]
	e.SourceIn, e.SourceOut, e.RecordIn, e.RecordOut = S[5], S[6], S[7], S[8]
	if S[9] != "" {
		e.Comments = append(e.Comments, S[9])
	}
	var time [4]int
	for i, s := range strings.Split(e.RecordIn, ":") {
		time[i], _ = strconv.Atoi(s)
		e.TimeIn[i] = s
	}
	e.FramesIn = time[0]*60*60*fps + time[1]*60*fps + time[2]*fps + time[3]
	for i, s := range strings.Split(e.RecordOut, ":") {
		time[i], _ = strconv.Atoi(s)
		e.TimeOut[i] = s
	}
	e.FramesOut = time[0]*60*60*fps + time[1]*60*fps + time[2]*fps + time[3]
	e.Elapsed = e.FramesOut - e.FramesIn
	e.Seconds = e.Elapsed / fps
	e.Frames = int(math.Mod(float64(e.Elapsed), float64(fps)))
	return e
}

func (e *Entry) String() string {
	buffer, err := json.MarshalIndent(e, "", "    ")
	if err != nil {
		return ""
	}
	return string(buffer[:])
}

func (e *Entry) CSV() []string {
	var record = []string{e.Event, e.Reel, e.TrackType, e.EditType, e.Transition,
		e.SourceIn, e.SourceOut}
	for _, s := range e.TimeIn {
		record = append(record, s)
	}
	for _, s := range e.TimeOut {
		record = append(record, s)
	}
	record = append(record, strconv.Itoa(e.FramesIn))
	record = append(record, strconv.Itoa(e.FramesOut))
	record = append(record, strconv.Itoa(e.Elapsed))
	record = append(record, strconv.Itoa(e.Seconds))
	record = append(record, strconv.Itoa(e.Frames))
	for _, s := range e.Comments {
		record = append(record, s)
	}
	return record
}

func Parse(r *bufio.Reader, fps int) []*Entry {
	var entries []*Entry
	var entry *Entry
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if S := REGEXP_ENTRY.FindStringSubmatch(line); S != nil {
			entry = NewEntry(S[1:], fps)
			entries = append(entries, entry)
		} else {
			if entry != nil && line[0] == '*' {
				entry.Comments = append(entry.Comments, line[1:])
			}
		}
	}
	return entries
}

func ConvertToCSV(r *bufio.Reader, w *bufio.Writer, fps int) error {
	var err error
	writer := csv.NewWriter(w)
	//writer.Comma = '\t'
	if err = writer.Write(CSV_HEADER); err != nil {
		return err
	}
	seconds, frames := 0, 0
	for _, entry := range Parse(r, fps) {
		if err = writer.Write(entry.CSV()); err != nil {
			return err
		}
		w.Flush()
		seconds, frames = seconds + entry.Seconds, frames + entry.Frames
	}
	var record [20]string
	record[18], record[19] = strconv.Itoa(seconds), strconv.Itoa(frames)
	if err = writer.Write(record[:]); err != nil {
		return err
	}
	w.Flush()
	record[18], record[19] = strconv.Itoa(seconds + int(frames/fps)),
		strconv.Itoa(frames - int(frames/fps)*fps)
	if err = writer.Write(record[:]); err != nil {
		return err
	}
	w.Flush()
	return nil
}

func ConvertFileToCSV(edlfile string, csvfile string, fps int) error {
	f, err := os.Open(edlfile)
	if err != nil {
		return err
	}
	defer f.Close()
	return ConvertToCSV(bufio.NewReader(f), bufio.NewWriter(os.Stdout), fps)
}
