package MUSIC_SERVICES

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	// "os"
	"errors"
	// "strings"
)

type MusicData struct {
	Author string `json:"Author"`
	Name   string `json:"Name"`
	URL    string `json:"URL"`
}

type AjaxPayload struct {
	MusicData []MusicData `json:"MusicData"`
}

type ResponseAjax struct {
	ResponseType string
	Status       string
	Message      string
	Payload      AjaxPayload
}

func (r *ResponseAjax) UnmarshalJSON(data []byte) error {
	// 1. Декодируем массив в срез сырых байт json.RawMessage
	var raw []json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// 2. Валидируем длину массива
	if len(raw) < 4 {
		return errors.New("invalid response: array must have at least 4 elements")
	}

	// 3. Извлекаем строковые поля
	if err := json.Unmarshal(raw[0], &r.ResponseType); err != nil {
		return fmt.Errorf("failed to parse ResponseType: %w", err)
	}
	if err := json.Unmarshal(raw[1], &r.Status); err != nil {
		return fmt.Errorf("failed to parse Status: %w", err)
	}
	if err := json.Unmarshal(raw[2], &r.Message); err != nil {
		return fmt.Errorf("failed to parse Message: %w", err)
	}

	// 4. Извлекаем объект с плейлистами и треками
	if err := json.Unmarshal(raw[3], &r.Payload); err != nil {
		return fmt.Errorf("failed to parse Payload: %w", err)
	}

	return nil
}

func SearchMusic(query string) ([]MusicData, error) {
	json, err := GetJson(query)
	if err != nil {
		return nil, err
	}
	return ParseJson(json)
}

func ParseJson(json_byte []byte) ([]MusicData, error) {
	var response ResponseAjax
	if err := json.Unmarshal(json_byte, &response); err != nil {
		return nil, err
	}
	log.Printf("Status: %s\nMessage: %s\n", response.Status, response.Message)
	return response.Payload.MusicData, nil
}

func GetJson(query string) ([]byte, error) {
	baseURL := "https://my.mail.ru/cgi-bin/my/ajax"
	params := url.Values{}
	params.Add("ajax_call", "1")
	params.Add("func_name", "music.search")
	params.Add("arg_extended", "1")
	params.Add("arg_search_params", `{"music":{"limit":100},"playlist":{"limit":50},"album":{"limit":10},"artist":{"limit":10}}`)
	params.Add("arg_offset", "0")
	params.Add("arg_limit", "100")

	params.Add("arg_query", query)
	finalURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	req, err := http.NewRequest(http.MethodGet, finalURL, nil)
	req.Header.Add("X-Requested-With", "XMLHttpRequest")
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)

	return data, nil
}

// func main() {
// 	// var query string
// 	// fmt.Println("Введите название композиции: ")
// 	// fmt.Scanln(&query)
// 	// err := GetJson(query)
// 	// if err != nil {
// 	// 	log.Fatal(err)
// 	// }
// 	tracks, err := ParseJSONFile("response.json")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	for i := range 10 {
// 		fmt.Printf("https:%s\n", tracks[i].URL)
// 		fmt.Printf("\x1b]8;;https:%s/\x1b\\%s - %s\x1b]8;;\x1b\\\n", tracks[i].URL, tracks[i].Author, tracks[i].Name)
// 	}
// 	// data, err := os.ReadFile("test.json")
// 	// if err != nil {
// 	// 	log.Fatal(err)
// 	// }
// 	// jsonFile := make(map[string]any)
// 	// if err = json.Unmarshal(data, &jsonFile); err != nil {
// 	// 	log.Fatal(err)
// 	// }
// 	// rawdata, ok := jsonFile["MusicData"]
// 	// if !ok {
// 	// 	log.Fatal("Key dosn't exist")
// 	// }
// 	// musicList, ok := rawdata.([]any)
// 	// if !ok {
// 	// 	log.Fatal("not an array")
// 	// }
// 	// for i := range 10 {
// 	// 	item := musicList[i]
// 	// 	musicData, ok := item.(map[string]any)
// 	// 	if !ok {
// 	// 		log.Printf("This %d not object", i)
// 	// 	}
// 	// 	fmt.Printf("%s - %s\n", musicData["Author"], musicData["Name"])
// 	// 	url, _ := musicData["URL"].(string)
// 	// 	fmt.Printf("URL: https:%s\n", url)
// 	// }
// }
