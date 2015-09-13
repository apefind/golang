package episodeguide

import (
	"fmt"
	"sort"
	"time"
)

type Series struct {
	Title, Description string
	Seasons            map[int]*Season
}

func NewSeries(title, description string) *Series {
	return &Series{
		Title:       title,
		Description: description,
		Seasons:     make(map[int]*Season),
	}
}

func (series *Series) String() string {
	return fmt.Sprintf("%s", series.Title)
}

func (series *Series) AddSeason(id int, title, description string) {
	series.Seasons[id] = NewSeason(series, id, title, description)
}

func (series *Series) SortedSeasons() []*Season {
	J := make([]int, 0, len(series.Seasons))
	for j := range series.Seasons {
		J = append(J, j)
	}
	sort.Ints(J)
	seasons := make([]*Season, len(J))
	for i, j := range J {
		seasons[i] = series.Seasons[j]
	}
	return seasons
}

// Return a map to directly access the episodes via episode code, e.g. "S02E03"
func (series *Series) EpisodeMap() map[string]*Episode {
	episodes := make(map[string]*Episode)
	for _, season := range series.SortedSeasons() {
		for _, episode := range season.SortedEpisodes() {
			episodes[episode.Code()] = episode
		}
	}
	return episodes
}

func GetSeries(title string, method string) (*Series, error) {
	var series *Series
	var err error
	done := make(chan error, 1)
	go func() {
		series, err = GetTVRageSeries(title)
		done <- err
	}()
	go func() {
		series, err = GetTVMazeSeries(title)
		done <- err
	}()
	var timeout time.Duration
	timeout, err = time.ParseDuration("5s")
	if err != nil {
		return series, err
	}
	select {
	case <-time.After(timeout):
		close(done)
		return series, fmt.Errorf("timeout after %s", timeout)
	case err = <-done:
		close(done)
		return series, err
	}
}
