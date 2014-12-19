package edl

import (
	"regexp"
)

// REGEXP_ENTRY extracts information from the EDL, e.g.
//		003  IP_EP213 V     C        00:59:21:23 00:59:21:24 00:59:58:00 00:59:58:01
//		004     BL    V     C        00:00:00:00 00:00:00:00 01:00:00:00 01:00:00:00
//		004  IMG_6549 V     D    010 03:01:42:19 03:01:43:16 01:00:00:00 01:00:00:27
//		005  1P6_PANO V     C        02:30:00:06 02:30:01:01 01:00:00:27 01:00:01:22
var REGEXP_ENTRY = regexp.MustCompile(`^\s*([0-9]+)\s+(\S+)\s+(\S+)\s+(\S+)\s+(\S*)\s*` +
	`([0-9]{2}:[0-9]{2}:[0-9]{2}:[0-9]{2})\s+` + `([0-9]{2}:[0-9]{2}:[0-9]{2}:[0-9]{2})\s+` +
	`([0-9]{2}:[0-9]{2}:[0-9]{2}:[0-9]{2})\s+` + `([0-9]{2}:[0-9]{2}:[0-9]{2}:[0-9]{2})\s*\*?(.*)$`)

var CSV_HEADER = []string{"Event No", "Reel", "Track Type", "Edit Type", "Transition",
	"Source In", "Source Out", "Prog In H", "Prog In M", "Prog In S", "Prog In F",
	"Prog Out H", "Prog Out M", "Prog Out S", "Prog Out F", "Frames In", "Frames Out",
	"Elapsed Frames", "Seconds", "Frames", "Comments"}
