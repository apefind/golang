package episodeguide

import (
	"fmt"
	"github.com/apefind/golang/shellutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type SeriesInfoCmd struct {
	Title           string
	SeasonID        int
	Timeout         time.Duration
	Method          string
	NormalizedTitle bool
	DryRun          bool
}

func (s *SeriesInfoCmd) GetSeriesReaders() []SeriesReader {
	readers := []SeriesReader{}
	if strings.Contains(s.Method, "tvmaze") {
		readers = append(readers, &TVMazeSeriesReader{})
	}
	if strings.Contains(s.Method, "tvrage") {
		readers = append(readers, &TVRageSeriesReader{})
	}
	return readers
}

func (s *SeriesInfoCmd) GetSeries() (*Series, error) {
	return GetSeries(s.Title, s.GetSeriesReaders(), s.Timeout)
}

// GetRenamedEpisodes returns a map with the renamed episodes
func (s *SeriesInfoCmd) GetRenamedEpisodes(filenames []string) map[string]string {
	var series *Series
	var episodes map[string]*Episode
	var err error
	renamedEpisodes := make(map[string]string)
	if s.NormalizedTitle {
		for _, filename := range filenames {
			if isVideoFile(filename) {
				renamedEpisodes[filename] = GetEpisodeCodeFromFilename(filename) + filepath.Ext(filename)
			}
		}
		return renamedEpisodes
	}
	series, err = s.GetSeries()
	if err != nil {
		log.Println(err)
		return renamedEpisodes
	}
	episodes = series.EpisodeMap()
	for _, filename := range filenames {
		if isVideoFile(filename) {
			code := GetEpisodeCodeFromFilename(filename)
			if episode, ok := episodes[code]; ok {
				renamedEpisodes[filename] = GetValidFilename(code + " " + episode.Title + filepath.Ext(filename))
			} else {
				renamedEpisodes[filename] = ""
			}
		}
	}
	return renamedEpisodes
}

func (s *SeriesInfoCmd) ListEpisodes() {
	printHeader := func(s string, c string) {
		log.Println()
		log.Println(s)
		log.Println(strings.Repeat(c, len(s)))
	}
	//_, seasonID := GetSeriesTitleFromPath(path)
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
	episodes := s.GetRenamedEpisodes(GetVideoFiles(path))
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
