package episodeguide

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

const TVMazeShowQueryURL string = `http://api.tvmaze.com/singlesearch/shows?q=%s`
const TVMazeEpisodeQueryURL string = `http://api.tvmaze.com/shows/%d/episodes`

type TVMazeEpisode struct {
	Id      int    `json:"number"`
	Season  int    `json:"season"`
	Title   string `json:"name"`
	Summary string `json:"summary"`
}

type TVMazeSeries struct {
	Name string `json:"name"`
	Id   int    `json:"id"`
}

func getTVMazeSeriesId(r io.Reader, title string) (int, error) {
	content, err := ioutil.ReadAll(r)
	if err != nil {
		return 0, err
	}
	var result TVMazeSeries
	if err := json.Unmarshal(content, &result); err != nil {
		return 0, err
	}
	return result.Id, nil
}

func getTVMazeSeries(r io.Reader, title string) (*Series, error) {
	content, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	var epsiodesTVMaze []TVMazeEpisode
	if err := json.Unmarshal(content, &epsiodesTVMaze); err != nil {
		return nil, err
	}
	series := NewSeries(title, "")
	for _, epsiodeTVMaze := range epsiodesTVMaze {
		if _, ok := series.Seasons[epsiodeTVMaze.Season]; !ok {
			series.AddSeason(epsiodeTVMaze.Season, "", "")
		}
		series.Seasons[epsiodeTVMaze.Season].AddEpisode(epsiodeTVMaze.Id, epsiodeTVMaze.Title, "")
	}
	return series, nil
}

func GetTVMazeSeriesId(title string) (int, error) {
	response, err := http.Get(fmt.Sprintf(TVMazeShowQueryURL, url.QueryEscape(title)))
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()
	return getTVMazeSeriesId(response.Body, title)
}

func GetTVMazeSeries(title string) (*Series, error) {
	id, err := GetTVMazeSeriesId(title)
	if err != nil {
		return nil, err
	}
	response, err := http.Get(fmt.Sprintf(TVMazeEpisodeQueryURL, id))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	return getTVMazeSeries(response.Body, title)
}
