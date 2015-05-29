package episodeguide

import (
	//	"fmt"
	"os"
)

func GetTVRageSeriesFake(title string) (*Series, error) {
	var r *os.File
	var err error
	r, err = os.Open("testdata/tvrage_search.xml")
	if err != nil {
		return nil, err
	}
	defer r.Close()
	_, err = getTVRageSeriesId(r, title)
	if err != nil {
		return nil, err
	}
	r, err = os.Open("testdata/tvrage_episode_list.xml")
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return getTVRageSeries(r)
}
