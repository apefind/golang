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
	GetSeries(title string) (*Series, error)
}

func GetSeries(title string, readers []SeriesReader, timeout time.Duration) (*Series, error) {
	var series *Series
	var err error
	done := make(chan error, 1)
	for _, r := range readers {
		go func(r SeriesReader) {
			series, err = r.GetSeries(title)
			if err == nil {
				done <- err
			}
		}(r)
	}
	select {
	case <-time.After(timeout):
		close(done)
		return nil, fmt.Errorf("timeout after %s", Timeout)
	case err = <-done:
		close(done)
		return series, err
	}
}

func GetSeries2(title string) (*Series, error) {
	readers := []SeriesReader{}
	if strings.Contains(Method, "tvmaze") {
		readers = append(readers, &TVMazeSeriesReader{})
	}
	if strings.Contains(Method, "tvrage") {
		readers = append(readers, &TVRageSeriesReader{})
	}
	return getSeries(title, readers, Timeout)
}

func GetSeriesReaders(method string) []SeriesReader {
	readers := []SeriesReader{}
	if strings.Contains(Method, "tvmaze") {
		readers = append(readers, &TVMazeSeriesReader{})
	}
	if strings.Contains(Method, "tvrage") {
		readers = append(readers, &TVRageSeriesReader{})
	}
	return readers
}
