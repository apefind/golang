// tvrage is down ...
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
	ID   string   `xml:"showid"`
	Name string   `xml:"name"`
}

type xmlTVRageQueryResult struct {
	_     xml.Name        `xml:"Results"`
	Shows []xmlTVRageShow `xml:"show"`
}

type xmlTVRageEpisode struct {
	_     xml.Name `xml:"episode"`
	ID    string   `xml:"seasonnum"`
	Title string   `xml:"title"`
}

type xmlTVRageSeason struct {
	_        xml.Name           `xml:"Season"`
	ID       string             `xml:"no,attr"`
	Episodes []xmlTVRageEpisode `xml:"episode"`
}

type xmlTVRageSeries struct {
	_       xml.Name          `xml:"Show"`
	Name    string            `xml:"name"`
	Seasons []xmlTVRageSeason `xml:"Episodelist>Season"`
}

func GetTVRageSeriesID(r io.Reader, title string) (string, error) {
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
			return show.ID, nil
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
		seasonID, err := strconv.Atoi(seasonTVRage.ID)
		if err != nil {
			return nil, err
		}
		if _, ok := series.Seasons[seasonID]; !ok {
			series.AddSeason(seasonID, "", "")
		}
		for _, episodeTVRage := range seasonTVRage.Episodes {
			episodeID, err := strconv.Atoi(episodeTVRage.ID)
			if err != nil {
				return nil, err
			}
			series.Seasons[seasonID].AddEpisode(episodeID, episodeTVRage.Title, "")
		}
	}
	return series, nil
}

type TVRageSeriesReader struct {
}

func (s *TVRageSeriesReader) GetSeries(title string) (*Series, error) {
	var response *http.Response
	var err error
	var id string
	response, err = http.Get(fmt.Sprintf(urlTVRageSearch, url.QueryEscape(title)))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	id, err = GetTVRageSeriesID(response.Body, title)
	if err != nil {
		return nil, err
	}
	response, err = http.Get(fmt.Sprintf(urlTVRageEpisodeList, url.QueryEscape(id)))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	return GetTVRageSeries(response.Body)
}
