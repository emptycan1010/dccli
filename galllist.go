package dccli

import (
	"encoding/json"
	"errors"
	"github.com/tidwall/gjson"
	"io"
	"log"
	"net/http"
)

type Getgalldata struct {
	GallList []GallList `json:"gall_list"`
	GallInfo []GallInfo `json:"gall_info"`
}

type GallList struct {
	No             string `json:"no"`
	Hit            string `json:"hit"`
	Recommend      string `json:"recommend"`
	Img_icon       string `json:"img_icon"`
	Movie_icon     string `json:"movie_icon"`
	Recommend_icon string `json:"recommend_icon"`
	Best_chk       string `json:"best_chk"`
	Realtime_chk   string `json:"realtime_chk"`
	Realtime_l_chk string `json:"realtime_l_chk"`
	Level          string `json:"level"`
	Total_comment  string `json:"total_comment"`
	Total_voice    string `json:"total_voice"`
	User_id        string `json:"user_id"`
	Voice_icon     string `json:"voice_icon"`
	Winnerta_icon  string `json:"winnerta_icon"`
	Member_icon    string `json:"member_icon"`
	Ip             string `json:"ip"`
	Subject        string `json:"subject"`
	Name           string `json:"name"`
	Date_time      string `json:"date_time"`
	Headtext       string `json:"headtext"`
}

type GallInfo struct {
	Gall_title    string    `json:"gall_title"`
	Category      string    `json:"category"`
	File_cnt      string    `json:"file_cnt"`
	File_size     string    `json:"file_size"`
	Is_minor      bool      `json:"is_minor"`
	Head_text     []HeadTXT `json:"head_text"`
	Notify_recent string    `json:"notify_recent"`
}

type HeadTXT struct {
	No       string `json:"no"`
	Name     string `json:"name"`
	Level    string `json:"level"`
	Selected bool   `json:"selected"`
}

func (s *Session) GetGallList(gallid string) (Getgalldata, error) {
	req, err := http.NewRequest("GET", HashedURLmake(gallid, s.Appid), nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("User-Agent", "dcinside.app")
	req.Host = "app.dcinside.com"
	req.Header.Set("referer", "https://app.dcinside.com")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return Getgalldata{}, errors.New("Error Posting Request")
	}
	bod, _ := io.ReadAll(res.Body)
	if gjson.Get(string(bod), "0.result").Exists() == true {
		return Getgalldata{}, errors.New("Please refresh your Appid")
	}
	var gg []Getgalldata
	err = json.Unmarshal(bod, &gg)
	if err != nil {
		return Getgalldata{}, errors.New("Error Unmarshaling Json")
	}
	return gg[0], nil
}
