package episodeguide

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/apefind/golang/shellutil"
)

type SeriesInfoCmd struct {
	Title           string
	SeasonID        int
	Timeout         time.Duration
	Source          string
	NormalizedTitle bool
	DryRun          bool
	UpdateRedis     bool
}

func (s *SeriesInfoCmd) GetSeriesReaders() []SeriesReader {
	readers := []SeriesReader{}
	if strings.Contains(s.Source, "tvmaze") {
		readers = append(readers, &TVMazeSeriesReader{})
	}
	if strings.Contains(s.Source, "tvrage") {
		readers = append(readers, &TVRageSeriesReader{})
	}
	return readers
}

func (s *SeriesInfoCmd) GetSeries() (*Series, error) {
	series, err := GetSeries(s.Title, s.GetSeriesReaders(), s.Timeout)
	if err != nil {
		return series, err
	}
	if strings.Contains(s.Source, "redis") {
		if err := UpdateRedisSeries(series); err != nil {
			return series, err
		}
	}
	return series, err
}

func (s *SeriesInfoCmd) GetRenamedNormalizedEpisodes(filenames []string) map[string]string {
	renamedEpisodes := make(map[string]string)
	for _, filename := range filenames {
		code := GetEpisodeCodeFromFilename(filename)
		if code != "" {
			renamedEpisodes[filename] = code + filepath.Ext(filename)
		} else {
			renamedEpisodes[filename] = ""
		}
	}
	return renamedEpisodes
}

// GetRenamedEpisodes returns a map with the renamed episodes
func (s *SeriesInfoCmd) GetRenamedEpisodes(filenames []string) map[string]string {
	var series *Series
	var episodes map[string]*Episode
	var renamedEpisodes map[string]string
	var err error
	done := make(chan error, 1)
	go func() {
		renamedEpisodes, err = GetRenamedRedisEpisodes(s.Title, filenames)
		if err != nil {
			return
		}
		done <- err
	}()
	go func() {
		renamedEpisodes = make(map[string]string)
		series, err = s.GetSeries()
		if err != nil {
			return
		}
		episodes = series.EpisodeMap()
		for _, filename := range filenames {
			code := GetEpisodeCodeFromFilename(filename)
			if episode, ok := episodes[code]; ok {
				renamedEpisodes[filename] = GetEpisodeFilename(code, episode.Title, filepath.Ext(filename))
			} else {
				renamedEpisodes[filename] = ""
			}
		}
		done <- err
	}()
	select {
	case <-time.After(s.Timeout):
		close(done)
		return renamedEpisodes
	case err = <-done:
		close(done)
		return renamedEpisodes
	}
}

func (s *SeriesInfoCmd) ListEpisodes() {
	printHeader := func(s string, c string) {
		log.Println()
		log.Println(s)
		log.Println(strings.Repeat(c, len(s)))
	}
	var series *Series
	var err error
	series, err = s.GetSeries()
	if err != nil {
		return
	}
	printHeader(series.Title, "=")
	for _, season := range series.SortedSeasons() {
		if s.SeasonID == 0 || season.ID == s.SeasonID {
			printHeader(fmt.Sprintf("Season %d", season.ID), "-")
			for _, episode := range season.SortedEpisodes() {
				log.Println(episode)
			}
		}
	}
}

// Rename epsiodes in a directory
func (s *SeriesInfoCmd) RenameEpisodes(path string) {
	var episodes map[string]string
	if s.NormalizedTitle {
		episodes = s.GetRenamedNormalizedEpisodes(GetVideoFiles(path))
	} else {
		episodes = s.GetRenamedEpisodes(GetVideoFiles(path))
	}
	filenames := make([]string, 0, len(episodes))
	for k := range episodes {
		filenames = append(filenames, k)
	}
	sort.Strings(filenames)
	for _, filename := range filenames {
		episode := episodes[filename]
		dirname, basename := filepath.Split(filename)
		if shellutil.IdenticalFilenames(basename, episode) {
			log.Println(basename, "-> ok")
		} else if episode == "" {
			log.Println(basename, "-> title not found")
		} else {
			log.Println(basename, "->", episode)
			if !s.DryRun {
				os.Rename(filename, dirname+string(filepath.Separator)+episode)
			}
		}
	}
}
