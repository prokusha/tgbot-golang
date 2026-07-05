package MUSIC_SERVICES

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	// "os"
	"encoding/json"
	"strings"

	"github.com/lrstanley/go-ytdlp"
	MP3_HELPER "github.com/prokusha/tgbot-golang/internal/helper/mp3"
	// "github.com/rs/zerolog/log"
)

type VideoResult struct {
	Title      string  `json:"title"`
	Uploader   string  `json:"uploader"`
	ID         string  `json:"id"`
	Duration   float64 `json:"duration"`
	WebpageURL string  `json:"webpage_url"`
}

type AudioResult struct {
	Track   string   `json:"track"`
	Artists []string `json:"artists"`

	isValid bool
}

func GetMusicURL(url string) (string, error) {
	log.Println("Start YTDLP. GET AUDIOURL...")
	dl := ytdlp.New().
		Format("ba").
		ExtractAudio().
		AudioFormat("mp3").
		AudioQuality("128K").
		EmbedMetadata().
		Paths("home:/tmp/").
		Print("after_move:filepath")

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	result, err := dl.Run(ctx, url)
	if err != nil {
		return "", fmt.Errorf("YTDLP run command error: %s", err)
	}
	log.Println("DONE!")

	return result.Stdout, nil
}

func SetTagsFromJSON(path_media string) error {
	index := strings.LastIndex(path_media, ".")
	path_json := path_media[:index]
	audio_result, err := parseInfoJson(path_json + ".info.json")
	if err != nil {
		return fmt.Errorf("can't get metadata from json: %s", err)
	}
	var tags MP3_HELPER.Tags
	tags.FilePath = path_media
	tags.Artist = strings.Join(audio_result.Artists, ",")
	tags.Title = audio_result.Track
	err = MP3_HELPER.SetTags(tags)
	return err
}

func parseInfoJson(path_json string) (AudioResult, error) {
	file, err := os.ReadFile(path_json)
	if err != nil {
		return AudioResult{}, fmt.Errorf("YT read json file error: %s", err)
	}
	info := AudioResult{isValid: false}
	if err = json.Unmarshal(file, &info); err == nil {
		info.isValid = true
	}
	return info, nil
}

func yt_search(query string) {
	searchQuery := fmt.Sprintf("ytsearch5:%s official audio", query)
	dl := ytdlp.New().SkipDownload().DumpJSON()
	result, err := dl.Run(context.TODO(), searchQuery)
	if err != nil {
		log.Panic(err)
	}
	// yt-dlp выводит по одной JSON-строке на каждый результат
	lines := strings.Split(strings.TrimSpace(result.Stdout), "\n")

	fmt.Printf("Найдено результатов: %d\n\n", len(lines))

	for _, line := range lines {
		var video VideoResult
		if err := json.Unmarshal([]byte(line), &video); err != nil {
			continue
		}
		fmt.Printf("• %s\n\n  ID: %s | Автор: %s | Ссылка: %s\n",
			video.Title, video.ID, video.Uploader, video.WebpageURL)
	}
}
