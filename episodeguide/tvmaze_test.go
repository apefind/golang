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
		return
	}
	for _, season := range series.SortedSeasons() {
		for _, episode := range season.SortedEpisodes() {
			fmt.Println(episode)
		}
	}
	// Output:
	// S01E01 Fire in the Hole
	// S01E02 Riverbrook
	// S01E03 Fixer
	// S01E04 Long in the Tooth
	// S01E05 The Lord of War and Thunder
	// S01E06 The Collection
	// S01E07 Blind Spot
	// S01E08 Blowback
	// S01E09 Hatless
	// S01E10 The Hammer
	// S01E11 Veterans
	// S01E12 Fathers and Sons
	// S01E13 Bulletville
	// S02E01 The Moonshine War
	// S02E02 The Life Inside
	// S02E03 The I of the Storm
	// S02E04 For Blood or Money
	// S02E05 Cottonmouth
	// S02E06 Blaze of Glory
	// S02E07 Save My Love
	// S02E08 The Spoil
	// S02E09 Brother's Keeper
	// S02E10 Debts and Accounts
	// S02E11 Full Commitment
	// S02E12 Reckoning
	// S02E13 Bloody Harlan
	// S03E01 The Gunfighter
	// S03E02 Cut Ties
	// S03E03 Harlan Roulette
	// S03E04 The Devil You Know
	// S03E05 Thick as Mud
	// S03E06 When the Guns Come Out
	// S03E07 The Man Behind the Curtain
	// S03E08 Watching the Detectives
	// S03E09 Loose Ends
	// S03E10 Guy Walks Into a Bar
	// S03E11 Measures
	// S03E12 Coalition
	// S03E13 Slaughterhouse
	// S04E01 Hole in the Wall
	// S04E02 Where's Waldo?
	// S04E03 Truth and Consequences
	// S04E04 The Bird Has Flown
	// S04E05 Kin
	// S04E06 Foot Chase
	// S04E07 Money Trap
	// S04E08 Outlaw
	// S04E09 The Hatchet Tour
	// S04E10 Get Drew
	// S04E11 Decoy
	// S04E12 Peace of Mind
	// S04E13 Ghosts
	// S05E01 A Murder Of Crowes
	// S05E02 The Kids Aren't All Right
	// S05E03 Good Intentions
	// S05E04 Over the Mountain
	// S05E05 Shot All to Hell
	// S05E06 Kill The Messenger
	// S05E07 Raw Deal
	// S05E08 Whistle Past the Graveyard
	// S05E09 Wrong Roads
	// S05E10 Weight
	// S05E11 The Toll
	// S05E12 Starvation
	// S05E13 Restitution
	// S06E01 Fate's Right Hand
	// S06E02 Cash Game
	// S06E03 Noblesse Oblige
	// S06E04 The Trash And The Snake
	// S06E05 Sounding
	// S06E06 Alive Day
	// S06E07 The Hunt
	// S06E08 Dark As a Dungeon
	// S06E09 Burned
	// S06E10 Trust
	// S06E11 Fugitive Number One
	// S06E12 Collateral
	// S06E13 The Promise
}
