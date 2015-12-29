package episodeguide

import (
	"fmt"
)

type Episode struct {
	Season      *Season
	ID          int
	Title       string
	Description string
}

func (episode *Episode) String() string {
	return episode.Code() + " " + episode.Title
}

func (episode *Episode) Code() string {
	return fmt.Sprintf("S%02dE%02d", episode.Season.ID, episode.ID)
}
