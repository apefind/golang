package episodeguide

import (
	//"encoding/json"

	"github.com/boltdb/bolt"
)

//func UpdateBoltSeries(series *Series) error {
//}

func GetRenamedBoltEpisodes(title string, filenames []string) (map[string]string, error) {
	renamedEpisodes := make(map[string]string)
	db, err := bolt.Open("bolt.db", 0600, nil)
	if err != nil {
		return renamedEpisodes, err
	}
	defer db.Close()
	return renamedEpisodes, nil
}
