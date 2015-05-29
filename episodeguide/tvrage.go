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

const TVRageSearchURL string = `http://services.tvrage.com/feeds/search.php?show=%s`
const TVRageEpisodeListURL string = `http://services.tvrage.com/feeds/episode_list.php?sid=%s`

type TVRageShow struct {
	_    xml.Name `xml:"show"`
	Id   string   `xml:"showid"`
	Name string   `xml:"name"`
}

type TVRageQueryResult struct {
	_     xml.Name     `xml:"Results"`
	Shows []TVRageShow `xml:"show"`
}

type TVRageEpisode struct {
	_     xml.Name `xml:"episode"`
	Id    string   `xml:"seasonnum"`
	Title string   `xml:"title"`
}

type TVRageSeason struct {
	_        xml.Name        `xml:"Season"`
	Id       string          `xml:"no,attr"`
	Episodes []TVRageEpisode `xml:"episode"`
}

type TVRageSeries struct {
	_       xml.Name       `xml:"Show"`
	Name    string         `xml:"name"`
	Seasons []TVRageSeason `xml:"Episodelist>Season"`
}

func getTVRageSeriesId(r io.Reader, title string) (string, error) {
	content, err := ioutil.ReadAll(r)
	if err != nil {
		return "", err
	}
	var result TVRageQueryResult
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

func getTVRageSeries(r io.Reader) (*Series, error) {
	content, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	var seriesTVRage TVRageSeries
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

func GetTVRageSeriesId(title string) (string, error) {
	response, err := http.Get(fmt.Sprintf(TVRageSearchURL, url.QueryEscape(title)))
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	return getTVRageSeriesId(response.Body, title)
}

func GetTVRageSeries(title string) (*Series, error) {
	id, err := GetTVRageSeriesId(title)
	if err != nil {
		return nil, err
	}
	response, err := http.Get(fmt.Sprintf(TVRageEpisodeListURL, url.QueryEscape(id)))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	return getTVRageSeries(response.Body)
}
