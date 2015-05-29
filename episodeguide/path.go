package episodeguide

import (
	"fmt"
	"mime"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// Check mime type for video or subtitle
func isVideoFile(filename string) bool {
	extension := filepath.Ext(filename)
	mime_type := mime.TypeByExtension(extension)
	return strings.HasPrefix(mime_type, "video") || strings.Contains(mime_type, "subtitle") ||
		strings.Contains(mime_type, "subrip") || extension == ".idx"
}

// Return video file in a directory
func GetVideoFiles(path string) []string {
	filenames, _ := filepath.Glob(path + string(filepath.Separator) + "*.*")
	video_files := make([]string, 0, len(filenames))
	for _, filename := range filenames {
		if isVideoFile(filename) {
			video_files = append(video_files, filename)
		}
	}
	return video_files
}

// Return series name from the current directory
func GetSeriesTitleFromWorkingDirectory(path string) (string, int) {
	path, err := os.Getwd()
	if err != nil {
		return "", 0
	}
	return GetSeriesTitleFromPath(path)
}

// Return series title and season id from a path
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

// Return "S01E01" from something like "4sj-dw-s05e01-dl-bluray-x264.mkv"
func GetEpisodeCodeFromFilename(filename string) string {
	_, basename := filepath.Split(filename)
	S := RE_EPISODE1.FindAllString(basename, -1)
	if len(S) > 0 {
		return strings.ToUpper(strings.Replace(S[0], ".", "", 1))
	}
	T := RE_EPISODE2.FindAllString(basename, -1)
	if len(T) > 0 {
		s := strings.Replace(T[0], "x", "E", 1)
		if len(s) == 4 {
			return "S0" + s
		} else {
			return "S" + s
		}
	}
	return ""
}

// Return a map with the renamed episodes
func GetRenamedEpisodes(title string, filenames []string, method string, noTitle bool) map[string]string {
	var series *Series
	var err error
	renamedEpisodes := make(map[string]string)
	if method == "tvrage" {
		series, err = GetTVRageSeries(title)
	}
	if err != nil {
		return renamedEpisodes
	}
	episodes := series.EpisodeMap()
	for _, filename := range filenames {
		if isVideoFile(filename) {
			code := GetEpisodeCodeFromFilename(filename)
			if noTitle {
				renamedEpisodes[filename] = code + filepath.Ext(filename)
			} else {
				if episode, ok := episodes[code]; ok {
					renamedEpisodes[filename] = code + " " + episode.Title + filepath.Ext(filename)
				} else {
					renamedEpisodes[filename] = ""
				}
			}
		}
	}
	return renamedEpisodes
}

// Rename epsiodes in a directory
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
		if basename == episode {
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
	title, seasonId := GetSeriesTitleFromPath(path)
	var series *Series
	var err error
	if method == "tvrage" {
		series, err = GetTVRageSeries(title)
	}
	if err != nil {
		return
	}
	printHeader(series.Title, "=")
	for _, season := range series.SortedSeasons() {
		if seasonId == 0 || season.Id == seasonId {
			printHeader(fmt.Sprintf("Season %d", season.Id), "-")
			for _, episode := range season.SortedEpisodes() {
				fmt.Println(episode)
			}
		}
	}
}