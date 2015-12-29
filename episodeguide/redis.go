package episodeguide

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"gopkg.in/redis.v3"
)

var RedisServer = "localhost"
var RedisPort = 6379
var RedisPassword = ""
var RedisDB int64 = 0

type RedisEpisode struct {
	Season      int    `json:"season"`
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func GetRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", RedisServer, RedisPort),
		Password: RedisPassword,
		DB:       RedisDB,
	})
}

func GetRedisKey(title, code string) string {
	return fmt.Sprintf("EpisodeGuide|%s|%s", title, code)
}

func GetRedisValue(episode *Episode) (string, error) {
	result, err := json.Marshal(&RedisEpisode{
		Season:      episode.Season.ID,
		ID:          episode.ID,
		Title:       episode.Title,
		Description: episode.Description,
	})
	return string(result), err
}

func UpdateRedisSeries(series *Series) error {
	client := GetRedisClient()
	for _, season := range series.SortedSeasons() {
		for _, episode := range season.SortedEpisodes() {
			value, err := GetRedisValue(episode)
			if err != nil {
				return err
			}
			val := client.Set(GetRedisKey(episode.Season.Series.Title, episode.Code()), value, 0)
			if err := val.Err(); err != nil {
				return err
			}
		}
	}
	return nil
}

func GetRenamedRedisEpisodes(title string, filenames []string) (map[string]string, error) {
	client := GetRedisClient()
	renamedEpisodes := make(map[string]string)
	for _, filename := range filenames {
		if isVideoFile(filename) {
			code := GetEpisodeCodeFromFilename(filename)
			value, err := client.Get(GetRedisKey(title, code)).Result()
			if err != nil {
				return renamedEpisodes, err
			}
			episode := RedisEpisode{}
			json.Unmarshal([]byte(value), &episode)
			renamedEpisodes[filename] = GetEpisodeFilename(code, episode.Title, filepath.Ext(filename))
		}
	}
	return renamedEpisodes, nil
}
