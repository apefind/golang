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
	Id      int    `json:"number"`
	Season  int    `json:"season"`
	Title   string `json:"name"`
	Summary string `json:"summary"`
}

type jsonTVMazeSeries struct {
	Name string `json:"name"`
	Id   int    `json:"id"`
}

func GetTVMazeSeriesId(r io.Reader, title string) (int, error) {
	content, err := ioutil.ReadAll(r)
	if err != nil {
		return 0, err
	}
	var result jsonTVMazeSeries
	if err := json.Unmarshal(content, &result); err != nil {
		return 0, err
	}
	return result.Id, nil
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
		series.Seasons[epsiodeTVMaze.Season].AddEpisode(epsiodeTVMaze.Id, epsiodeTVMaze.Title, "")
	}
	return series, nil
}

type TVMazeSeries struct {
	Title  string
	Id     int
	Series *Series
}

func NewTVMazeSeries(title string) *TVMazeSeries {
	return &TVMazeSeries{
		Title:  title,
		Id:     0,
		Series: nil,
	}
}

func (s *TVMazeSeries) String() string {
	return "tvmaze series: " + s.Title
}

func (s *TVMazeSeries) GetId() error {
	var response *http.Response
	var err error
	response, err = http.Get(fmt.Sprintf(urlTVMazeShowQuery, url.QueryEscape(s.Title)))
	if err != nil {
		return err
	}
	defer response.Body.Close()
	s.Id, err = GetTVMazeSeriesId(response.Body, s.Title)
	return err
}

func (s *TVMazeSeries) Get() error {
	var response *http.Response
	var err error
	err = s.GetId()
	if err != nil {
		return err
	}
	response, err = http.Get(fmt.Sprintf(urlTVMazeEpisodeQuery, s.Id))
	if err != nil {
		return err
	}
	defer response.Body.Close()
	s.Series, err = GetTVMazeSeries(response.Body, s.Title)
	return err
}

func (s *TVMazeSeries) GetSeries() (*Series, error) {
	if err := s.Get(); err != nil {
		return nil, err
	}
	return s.Series, nil
}
