package episodeguide

import (
	"fmt"
	"os"
)

func GetTVMazeSeriesFromTestData() (*Series, error) {
	var r *os.File
	var err error
	r, err = os.Open("testdata/tvmaze_search.json")
	if err != nil {
		return nil, err
	}
	defer r.Close()
	_, err = GetTVMazeSeriesID(r, "Justified")
	if err != nil {
		return nil, err
	}
	r, err = os.Open("testdata/tvmaze_episode_list.json")
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return GetTVMazeSeries(r, "Justified")
}

func ExampleTVMaze() {
	series, err := GetTVMazeSeriesFromTestData()
	if err != nil {
		fmt.Println(err)
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
	// S01E01 Fire in the Hole
	// S01E02 Riverbrook
	// S01E03 Fixer
	// S01E04 Long in the Tooth
	// S02E01 The Moonshine War
	// S02E02 The Life Inside
	// S02E03 The I of the Storm
	// S02E04 For Blood or Money
	// S03E01 The Gunfighter
	// S03E02 Cut Ties
	// S03E03 Harlan Roulette
	// S03E04 The Devil You Know
	// S04E01 Hole in the Wall
	// S04E02 Where's Waldo?
	// S04E03 Truth and Consequences
	// S04E04 The Bird Has Flown
}
