package episodeguide

import (
	"fmt"
	"os"
)

func GetTVRageSeriesFromTestData(title string) (*Series, error) {
	var r *os.File
	var err error
	r, err = os.Open("testdata/tvrage_search.xml")
	if err != nil {
		return nil, err
	}
	defer r.Close()
	_, err = GetTVRageSeriesID(r, title)
	if err != nil {
		return nil, err
	}
	r, err = os.Open("testdata/tvrage_episode_list.xml")
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return GetTVRageSeries(r)
}

func ExampleTVRage() {
	series, err := GetTVRageSeriesFromTestData("Buffy The Vampire Slayer")
	if err != nil {
		return
	}
	for i, season := range series.SortedSeasons() {
		for j, episode := range season.SortedEpisodes() {
			fmt.Println(episode)
			if j > 2 {
				break
			}
		}
		if i > 2 {
			break
		}
	}
	// Output:
	// S01E01 Welcome to the Hellmouth (1)
	// S01E02 The Harvest (2)
	// S01E03 Witch
	// S01E04 Teacher's Pet
	// S02E01 When She Was Bad
	// S02E02 Some Assembly Required
	// S02E03 School Hard
	// S02E04 Inca Mummy Girl
	// S03E01 Anne
	// S03E02 Dead Man's Party
	// S03E03 Faith, Hope & Trick
	// S03E04 Beauty and the Beasts
	// S04E01 The Freshman
	// S04E02 Living Conditions
	// S04E03 The Harsh Light of Day
	// S04E04 Fear, Itself
}
