package episodeguide

import (
	"regexp"
)

// Regular expression for season and episode
var RE_EPISODE1 = regexp.MustCompile(`[sS][0-2][0-9].?[eE][0-3][0-9]`)
var RE_EPISODE2 = regexp.MustCompile(`[0-2]?[0-9]x[0-3][0-9]`)
var RE_EPISODE3 = regexp.MustCompile(`[.-][0-9][0-3][0-9][.-]`)
