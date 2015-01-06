package edl

import (
	"bufio"
	"encoding/json"
	"math"
	"regexp"
	"strconv"
	"strings"
)

// RegExpEntry extracts information from the EDL, e.g.
//		003  IP_EP213 V     C        00:59:21:23 00:59:21:24 00:59:58:00 00:59:58:01
//		004     BL    V     C        00:00:00:00 00:00:00:00 01:00:00:00 01:00:00:00
//		004  IMG_6549 V     D    010 03:01:42:19 03:01:43:16 01:00:00:00 01:00:00:27
//		005  1P6_PANO V     C        02:30:00:06 02:30:01:01 01:00:00:27 01:00:01:22
var RegExpEntry = regexp.MustCompile(`^\s*([0-9]+)\s+(\S+)\s+(\S+)\s+(\S+)\s+(\S*)\s*` +
	`([0-9]{2}:[0-9]{2}:[0-9]{2}:[0-9]{2})\s+` + `([0-9]{2}:[0-9]{2}:[0-9]{2}:[0-9]{2})\s+` +
	`([0-9]{2}:[0-9]{2}:[0-9]{2}:[0-9]{2})\s+` + `([0-9]{2}:[0-9]{2}:[0-9]{2}:[0-9]{2})\s*\*?(.*)$`)

// scanLines also works for mac classic `\r` endings, there are probably better ways to do this ...
func scanLines(data []byte, atEOF bool) (advance int, token []byte, err error) {
	innerline, endline := regexp.MustCompile("\r([^\n])"), regexp.MustCompile("\r$")
	replaced := endline.ReplaceAll(innerline.ReplaceAll(data, []byte("\n$1")), []byte("\n"))
	return bufio.ScanLines(replaced, atEOF)
}

type Entry struct {
	Event, Reel, TrackType, EditType, Transition string
	SourceIn, SourceOut                          string
	RecordIn, RecordOut                          string
	Comment                                      []string
	TimeIn, TimeOut                              [4]string
	FramesIn, FramesOut                          int
	Elapsed, Seconds, Frames                     int
}

func NewEntry(S []string, fps int) *Entry {
	e := &Entry{Comment: make([]string, 0, 10)}
	e.Event, e.Reel, e.TrackType, e.EditType, e.Transition = S[0], S[1], S[2], S[3], S[4]
	e.SourceIn, e.SourceOut, e.RecordIn, e.RecordOut = S[5], S[6], S[7], S[8]
	if S[9] != "" {
		e.Comment = append(e.Comment, strings.TrimSpace(S[9]))
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
	for _, s := range e.Comment {
		record = append(record, s)
	}
	return record
}

func Parse(r *bufio.Reader, fps int) []*Entry {
	var entries []*Entry
	var entry *Entry
	scanner := bufio.NewScanner(r)
	scanner.Split(scanLines)
	for scanner.Scan() {
		line := scanner.Text()
		if S := RegExpEntry.FindStringSubmatch(line); S != nil {
			entry = NewEntry(S[1:], fps)
			entries = append(entries, entry)
		} else {
			if entry != nil && len(line) > 0 && line[0] == '*' {
				entry.Comment = append(entry.Comment, strings.TrimSpace(line[1:]))
			}
		}
	}
	return entries
}
