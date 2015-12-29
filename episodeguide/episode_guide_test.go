package episodeguide

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestPath(t *testing.T) {
	t.Log("testing path functionality")
	episodes := make(map[string]string)
	episodes["4sj-dw-s05e01-dl-bluray-x264.mkv"] = "S05E01"
	episodes["4sj-dw-s05e02-dl-bluray-x264.mkv"] = "S05E02"
	episodes["4sj-dw-5x02-dl-bluray-x264.mkv"] = "S05E02"
	episodes["4sj-dw-05x02-dl-bluray-x264.mkv"] = "S05E02"
	episodes["4sj-dw-f02-dl-bluray-x264.mkv"] = ""
	episodes["blabla S01.E02 blabla"] = "S01E02"
	episodes["tvs-ded-dl-ithd-x264-304.mkv"] = "S03E04"
	episodes["tvs-ded-dl-ithd-x264.304.mkv"] = "S03E04"
	episodes["tvs-ded-dl-ithd-x264+304.mkv"] = ""
	episodes["tvs-ded-dl-ithd-x264+304+xyz.mkv"] = ""
	for path := range episodes {
		episode := GetEpisodeCodeFromFilename(path)
		if episode != episodes[path] {
			t.Error("expected", episodes[path], ", but got", episode)
		}
	}
}

func TestFilename(t *testing.T) {
	t.Log("testing filename functionality")
	filenames := make(map[string]bool)
	filenames["4sj-dw-s05e01-dl-bluray-x264.mkv"] = true
	filenames["4sj-dw-s05e01-dl-bluray-x264.avi"] = true
	filenames["4sj-dw-s05e01-dl-bluray-x264.wmv"] = true
	filenames["4sj-dw-s05e01-dl-bluray-x264.jpeg"] = false
	filenames["4sj-dw-s05e01-dl-bluray-x264.txt"] = false
	for filename := range filenames {
		if isVideoFile(filename) != filenames[filename] {
			if filenames[filename] {
				t.Error("expected", filename, "to be a video file")
			} else {
				t.Error("expected", filename, "not to be a video file")
			}
		}
	}
}

func TestEpisode(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping episode funktionality in short mode")
	}
	t.Log("testing episode functionality")
	tmpdir, err := ioutil.TempDir("/tmp", "test_episode_guide")
	if err != nil {
		t.Error("cannot create temporary directory")
	}
	path := tmpdir + string(filepath.Separator) + "The Simpsons" + string(filepath.Separator) + "Season5"
	if os.MkdirAll(path, 0700) != nil {
		t.Error("cannot create example directory")
	}
	t.Log("example directory", path)
	filenames := []string{
		"4sj-dw-s05e01-dl-bluray-x264.mkv",
		"4sj-dw-5x03-dl-bluray-x264.mkv",
		"S99E12-dl-bluray-x264.mkv",
		"S05E02 Cape Feare.mkv",
	}
	for _, filename := range filenames {
		f, err := os.Create(path + string(filepath.Separator) + filename)
		if err != nil {
			t.Error("cannot create file", filename)
		}
		defer f.Close()
	}
	var info SeriesInfoCmd
	info.Title, info.SeasonID = GetSeriesTitleFromPath(path)
	info.Method = "tvmaze|tvrage|redis"
	info.Timeout = 5 * time.Second
	info.NormalizedTitle = false
	info.DryRun = false
	episodes := info.GetRenamedEpisodes(filenames)
	renamedFiles := []string{
		"S05E01 Homer's Barbershop Quartet.mkv",
		"S05E03 Homer Goes to College.mkv",
		"",
		"S05E02 Cape Feare.mkv",
	}
	for i, filename := range filenames {
		if episodes[filename] != renamedFiles[i] {
			t.Error("expected", renamedFiles[i], "instead of", episodes[filename])
		}
	}
	info.RenameEpisodes(filepath.Clean(path))
	info.NormalizedTitle = true
	episodes = info.GetRenamedEpisodes(filenames)
	normalizedFiles := []string{
		"S05E01.mkv",
		"S05E03.mkv",
		"",
		"S05E02.mkv",
	}
	for i, filename := range filenames {
		if episodes[filename] != normalizedFiles[i] {
			t.Error("expected", normalizedFiles[i], "instead of", episodes[filename])
		}
	}
	info.RenameEpisodes(filepath.Clean(path))
}

func ExampleGetEpisodeFromFilename() {
	episodes := []string{
		"4sj-dw-s05e01-dl-bluray-x264.mkv",
		"4sj-dw-s05e02-dl-bluray-x264.mkv",
		"4sj-dw-5x02-dl-bluray-x264.mkv",
		"4sj-dw-05x02-dl-bluray-x264.mkv",
		"4sj-dw-f02-dl-bluray-x264.mkv",
	}
	for _, path := range episodes {
		fmt.Println(path, "->", GetEpisodeCodeFromFilename(path))
	}
	// Output:
	// 4sj-dw-s05e01-dl-bluray-x264.mkv -> S05E01
	// 4sj-dw-s05e02-dl-bluray-x264.mkv -> S05E02
	// 4sj-dw-5x02-dl-bluray-x264.mkv -> S05E02
	// 4sj-dw-05x02-dl-bluray-x264.mkv -> S05E02
	// 4sj-dw-f02-dl-bluray-x264.mkv ->
}

func ExampleGetSeriesFromPath() {
	paths := []string{
		"/Users/test/Movies/Millenium",
		"/Users/test/Movies/Millenium/Season3",
		"/Users/test/Movies/X-Files/season3",
		"Movies/X-Files",
	}
	for _, path := range paths {
		title, seasonID := GetSeriesTitleFromPath(path)
		fmt.Println(path, "->", title, seasonID)
	}
	// Output:
	// /Users/test/Movies/Millenium -> Millenium 0
	// /Users/test/Movies/Millenium/Season3 -> Millenium 3
	// /Users/test/Movies/X-Files/season3 -> X-Files 3
	// Movies/X-Files -> X-Files 0
}
