package episodeguide

import (
	"fmt"
	"sort"
)

type Season struct {
	Series             *Series
	ID                 int
	Title, Description string
	Episodes           map[int]*Episode
}

func NewSeason(series *Series, id int, title, description string) *Season {
	return &Season{
		Series:      series,
		ID:          id,
		Title:       title,
		Description: description,
		Episodes:    make(map[int]*Episode),
	}
}

func (season *Season) String() string {
	return fmt.Sprintf("Season %d: %s", season.ID, season.Title)
}

func (season *Season) AddEpisode(id int, title, description string) {
	season.Episodes[id] = &Episode{
		Season:      season,
		ID:          id,
		Title:       title,
		Description: description,
	}
}

func (season *Season) SortedEpisodes() []*Episode {
	J := make([]int, 0, len(season.Episodes))
	for j := range season.Episodes {
		J = append(J, j)
	}
	sort.Ints(J)
	episodes := make([]*Episode, len(J))
	for i, j := range J {
		episodes[i] = season.Episodes[j]
	}
	return episodes
}
