package episodeguide

import (
	"apefind/shellutil"
	"fmt"
	"mime"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

// isVideoFile checks mime type for video or subtitle
func isVideoFile(filename string) bool {
	extension := filepath.Ext(filename)
	mimeType := mime.TypeByExtension(extension)
	return strings.HasPrefix(mimeType, "video") || strings.Contains(mimeType, "subtitle") ||
		strings.Contains(mimeType, "subrip") || extension == ".idx"
}

// GetVideoFiles returns video files in a directory
func GetVideoFiles(path string) []string {
	filenames, _ := filepath.Glob(path + string(filepath.Separator) + "*.*")
	videoFiles := make([]string, 0, len(filenames))
	for _, filename := range filenames {
		if isVideoFile(filename) {
			videoFiles = append(videoFiles, filename)
		}
	}
	return videoFiles
}

// GetSeriesTitleFromWorkingDirectory returns the series name from the current directory
func GetSeriesTitleFromWorkingDirectory(path string) (string, int) {
	path, err := os.Getwd()
	if err != nil {
		return "", 0
	}
	return GetSeriesTitleFromPath(path)
}

// GetSeriesTitleFromPath returns the series title and season id from a path
func GetSeriesTitleFromPath(path string) (string, int) {
	dirname, basename := filepath.Split(filepath.Clean(path))
	if strings.HasPrefix(strings.ToLower(basename), "season") {
		dirname, _ := filepath.Split(dirname)
		id, err := strconv.Atoi(basename[6:])
		if err != nil {
			return filepath.Base(dirname), 0
		}
		return filepath.Base(dirname), id
	}
	return basename, 0
}

// GetEpisodeCodeFromFilename returns "S05E01" from something like "4sj-dw-s05e01-dl-bluray-x264.mkv"
func GetEpisodeCodeFromFilename(filename string) string {
	var S []string
	var reEpisode *regexp.Regexp
	_, basename := filepath.Split(filename)
	reEpisode = regexp.MustCompile(`[sS][0-2][0-9].?[eE][0-3][0-9]`)
	S = reEpisode.FindAllString(basename, -1)
	if len(S) > 0 {
		return strings.ToUpper(strings.Replace(S[0], ".", "", 1))
	}
	reEpisode = regexp.MustCompile(`[0-2]?[0-9]x[0-3][0-9]`)
	S = reEpisode.FindAllString(basename, -1)
	if len(S) > 0 {
		s := strings.Replace(S[0], "x", "E", 1)
		if len(s) == 4 {
			return "S0" + s
		}
		return "S" + s
	}
	reEpisode = regexp.MustCompile(`[.-][0-9][0-3][0-9][.-]`)
	S = reEpisode.FindAllString(basename, -1)
	if len(S) > 0 {
		return "S0" + S[0][1:2] + "E" + S[0][2:4]
	}
	return ""
}

func GetValidFilename(filename string) string {
	return strings.Replace(filename, "/", "-", -1)
}

// GetRenamedEpisodes returns a map with the renamed episodes
func GetRenamedEpisodes(title string, filenames []string, method string, noTitle bool) map[string]string {
	var series *Series
	var episodes map[string]*Episode
	var err error
	renamedEpisodes := make(map[string]string)
	if noTitle {
		for _, filename := range filenames {
			if isVideoFile(filename) {
				renamedEpisodes[filename] = GetEpisodeCodeFromFilename(filename) + filepath.Ext(filename)
			}
		}
		return renamedEpisodes
	}
	series, err = GetSeries(title)
	if err != nil {
		fmt.Println(err)
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

// RenameEpisodes renames epsiodes in a directory
func RenameEpisodes(path string, method string, dryRun bool, noTitle bool) {
	title, _ := GetSeriesTitleFromPath(path)
	episodes := GetRenamedEpisodes(title, GetVideoFiles(path), method, noTitle)
	filenames := make([]string, 0, len(episodes))
	for k := range episodes {
		filenames = append(filenames, k)
	}
	sort.Strings(filenames)
	for _, filename := range filenames {
		episode := episodes[filename]
		dirname, basename := filepath.Split(filename)
		if shellutil.IdenticalFilenames(basename, episode) {
			fmt.Println(basename, "-> ok")
		} else if episode == "" {
			fmt.Println(basename, "-> title not found")
		} else {
			fmt.Println(basename, "->", episode)
			if !dryRun {
				os.Rename(filename, dirname+string(filepath.Separator)+episode)
			}
		}
	}
}

func ListEpisodes(path string, method string) {
	printHeader := func(s string, c string) {
		fmt.Println()
		fmt.Println(s)
		fmt.Println(strings.Repeat(c, len(s)))
	}
	title, seasonID := GetSeriesTitleFromPath(path)
	var series *Series
	var err error
	series, err = GetSeries(title)
	if err != nil {
		return
	}
	printHeader(series.Title, "=")
	for _, season := range series.SortedSeasons() {
		if seasonID == 0 || season.ID == seasonID {
			printHeader(fmt.Sprintf("Season %d", season.ID), "-")
			for _, episode := range season.SortedEpisodes() {
				fmt.Println(episode)
			}
		}
	}
}

type SeriesInfo struct {
	Title           string
	Timeout         time.Duration
	Method          string
	NormatizedTitle bool
}

// GetRenamedEpisodes returns a map with the renamed episodes
func (s *SeriesInfo) GetRenamedEpisodes(filenames []string) map[string]string {
	var series *Series
	var episodes map[string]*Episode
	var err error
	renamedEpisodes := make(map[string]string)
	if s.NormatizedTitle {
		for _, filename := range filenames {
			if isVideoFile(filename) {
				renamedEpisodes[filename] = GetEpisodeCodeFromFilename(filename) + filepath.Ext(filename)
			}
		}
		return renamedEpisodes
	}
	series, err = GetSeries(s.Title, GetSeriesReaders(s.Method), s.Timeout)
	if err != nil {
		fmt.Println(err)
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

func (s *SeriesInfo) ListEpisodes(path string) {
	printHeader := func(s string, c string) {
		fmt.Println()
		fmt.Println(s)
		fmt.Println(strings.Repeat(c, len(s)))
	}
	title, seasonID := GetSeriesTitleFromPath(path)
	var series *Series
	var err error
	series, err = GetSeries(title, GetSeriesReaders(s.Method), s.Timeout)
	if err != nil {
		return
	}
	printHeader(series.Title, "=")
	for _, season := range series.SortedSeasons() {
		if seasonID == 0 || season.ID == seasonID {
			printHeader(fmt.Sprintf("Season %d", season.ID), "-")
			for _, episode := range season.SortedEpisodes() {
				fmt.Println(episode)
			}
		}
	}
}

// Rename epsiodes in a directory
func (s *SeriesInfo) RenameEpisodes(path string) {
	title, _ := GetSeriesTitleFromPath(path)
	episodes := s.GetRenamedEpisodes(title, GetVideoFiles(path), s.Method, s.NormalizedTitle)
	filenames := make([]string, 0, len(episodes))
	for k := range episodes {
		filenames = append(filenames, k)
	}
	sort.Strings(filenames)
	for _, filename := range filenames {
		episode := episodes[filename]
		dirname, basename := filepath.Split(filename)
		if shellutil.IdenticalFilenames(basename, episode) {
			fmt.Println(basename, "-> ok")
		} else if episode == "" {
			fmt.Println(basename, "-> title not found")
		} else {
			fmt.Println(basename, "->", episode)
			if !dryRun {
				os.Rename(filename, dirname+string(filepath.Separator)+episode)
			}
		}
	}
}
