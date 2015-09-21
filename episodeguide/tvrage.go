package episodeguide

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const urlTVRageSearch string = `http://services.tvrage.com/feeds/search.php?show=%s`
const urlTVRageEpisodeList string = `http://services.tvrage.com/feeds/episode_list.php?sid=%s`

type xmlTVRageShow struct {
	_    xml.Name `xml:"show"`
	Id   string   `xml:"showid"`
	Name string   `xml:"name"`
}

type xmlTVRageQueryResult struct {
	_     xml.Name        `xml:"Results"`
	Shows []xmlTVRageShow `xml:"show"`
}

type xmlTVRageEpisode struct {
	_     xml.Name `xml:"episode"`
	Id    string   `xml:"seasonnum"`
	Title string   `xml:"title"`
}

type xmlTVRageSeason struct {
	_        xml.Name           `xml:"Season"`
	Id       string             `xml:"no,attr"`
	Episodes []xmlTVRageEpisode `xml:"episode"`
}

type xmlTVRageSeries struct {
	_       xml.Name          `xml:"Show"`
	Name    string            `xml:"name"`
	Seasons []xmlTVRageSeason `xml:"Episodelist>Season"`
}

func GetTVRageSeriesId(r io.Reader, title string) (string, error) {
	content, err := ioutil.ReadAll(r)
	if err != nil {
		return "", err
	}
	var result xmlTVRageQueryResult
	if err := xml.Unmarshal(content, &result); err != nil {
		return "", err
	}
	for _, show := range result.Shows {
		if strings.ToLower(title) == strings.ToLower(show.Name) {
			return show.Id, nil
		}
	}
	return "", fmt.Errorf("show id for <%s> not found", title)
}

func GetTVRageSeries(r io.Reader) (*Series, error) {
	content, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	var seriesTVRage xmlTVRageSeries
	if err := xml.Unmarshal(content, &seriesTVRage); err != nil {
		return nil, err
	}
	series := NewSeries(seriesTVRage.Name, "")
	for _, seasonTVRage := range seriesTVRage.Seasons {
		seasonId, err := strconv.Atoi(seasonTVRage.Id)
		if err != nil {
			return nil, err
		}
		if _, ok := series.Seasons[seasonId]; !ok {
			series.AddSeason(seasonId, "", "")
		}
		for _, episodeTVRage := range seasonTVRage.Episodes {
			episodeId, err := strconv.Atoi(episodeTVRage.Id)
			if err != nil {
				return nil, err
			}
			series.Seasons[seasonId].AddEpisode(episodeId, episodeTVRage.Title, "")
		}
	}
	return series, nil
}

type TVRageSeries struct {
	Title  string
	Id     string
	Series *Series
}

func NewTVRageSeries(title string) *TVRageSeries {
	return &TVRageSeries{
		Title:  title,
		Id:     "",
		Series: nil,
	}
}

func (s *TVRageSeries) String() string {
	return "tvrage series: " + s.Title
}

func (s *TVRageSeries) GetId() error {
	var response *http.Response
	var err error
	response, err = http.Get(fmt.Sprintf(urlTVRageSearch, url.QueryEscape(s.Title)))
	if err != nil {
		return err
	}
	defer response.Body.Close()
	s.Id, err = GetTVRageSeriesId(response.Body, s.Title)
	return err
}

func (s *TVRageSeries) Get() error {
	var response *http.Response
	var err error
	err = s.GetId()
	if err != nil {
		return err
	}
	response, err = http.Get(fmt.Sprintf(urlTVRageEpisodeList, url.QueryEscape(s.Id)))
	if err != nil {
		return err
	}
	defer response.Body.Close()
	s.Series, err = GetTVRageSeries(response.Body)
	return err
}

func (s *TVRageSeries) GetSeries() (*Series, error) {
	if err := s.Get(); err != nil {
		return nil, err
	}
	return s.Series, nil
}
