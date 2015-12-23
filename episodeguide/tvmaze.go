package episodeguide

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

const urlTVMazeShowQuery string = `http://api.tvmaze.com/singlesearch/shows?q=%s`
const urlTVMazeEpisodeQuery string = `http://api.tvmaze.com/shows/%d/episodes`

type jsonTVMazeEpisode struct {
	ID      int    `json:"number"`
	Season  int    `json:"season"`
	Title   string `json:"name"`
	Summary string `json:"summary"`
}

type jsonTVMazeSeries struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

func GetTVMazeSeriesID(r io.Reader, title string) (int, error) {
	content, err := ioutil.ReadAll(r)
	if err != nil {
		return 0, err
	}
	var result jsonTVMazeSeries
	if err := json.Unmarshal(content, &result); err != nil {
		return 0, err
	}
	return result.ID, nil
}

func GetTVMazeSeries(r io.Reader, title string) (*Series, error) {
	content, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	var epsiodesTVMaze []jsonTVMazeEpisode
	if err := json.Unmarshal(content, &epsiodesTVMaze); err != nil {
		return nil, err
	}
	series := NewSeries(title, "")
	for _, epsiodeTVMaze := range epsiodesTVMaze {
		if _, ok := series.Seasons[epsiodeTVMaze.Season]; !ok {
			series.AddSeason(epsiodeTVMaze.Season, "", "")
		}
		series.Seasons[epsiodeTVMaze.Season].AddEpisode(epsiodeTVMaze.ID, epsiodeTVMaze.Title, "")
	}
	return series, nil
}

type TVMazeSeriesReader struct {
}

func (s *TVMazeSeriesReader) GetSeries(title string) (*Series, error) {
	var response *http.Response
	var err error
	var id int
	response, err = http.Get(fmt.Sprintf(urlTVMazeShowQuery, url.QueryEscape(title)))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	id, err = GetTVMazeSeriesID(response.Body, title)
	if err != nil {
		return nil, err
	}
	response, err = http.Get(fmt.Sprintf(urlTVMazeEpisodeQuery, id))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	return GetTVMazeSeries(response.Body, title)
}
