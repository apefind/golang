package episodeguide

import (
	"mime"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
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
func GetSeriesTitleFromWorkingDirectory() (string, int) {
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
