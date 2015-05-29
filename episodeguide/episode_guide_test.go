package episodeguide

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
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
	RenameEpisodes(filepath.Clean(path), "tvrage", true, false)
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
		title, _ := GetSeriesTitleFromPath(path)
		fmt.Println(path, "->", title)
	}
	// Output:
	// /Users/test/Movies/Millenium -> Millenium
	// /Users/test/Movies/Millenium/Season3 -> Millenium
	// /Users/test/Movies/X-Files/season3 -> X-Files
	// Movies/X-Files -> X-Files
}
