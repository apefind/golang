package episodeguide

import (
	"fmt"
	"sort"
	"strings"
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

type SeriesReader interface {
	GetSeries() (*Series, error)
}

func getSeries(title string, readers []SeriesReader, timeout time.Duration) (*Series, error) {
	var series *Series
	var err error
	done := make(chan error, 1)
	for _, r := range readers {
		go func(r SeriesReader) {
			series, err = r.GetSeries()
			if err == nil {
				done <- err
			}
		}(r)
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

func GetSeries(title string, method string) (*Series, error) {
	readers := []SeriesReader{}
	if strings.Contains(method, "tvmaze") {
		readers = append(readers, NewTVMazeSeries(title))
	}
	if strings.Contains(method, "tvrage") {
		readers = append(readers, NewTVRageSeries(title))
	}
	timeout, err := time.ParseDuration("5.0s")
	if err != nil {
		return nil, err
	}
	return getSeries(title, readers, timeout)
}
