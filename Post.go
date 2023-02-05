package dccli

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Post struct {
	View_info PostViewInfo `json:"view_info"`
	View_Main PostViewMain `json:"view_main"`
}

type PostViewInfo struct {
	Galltitle      string `json:"galltitle"`
	Category       string `json:"category"`
	Subject        string `json:"subject"`
	No             string `json:"no"`
	Name           string `json:"name"`
	Level          string `json:"level"`
	Member_icon    string `json:"member_icon"`
	Total_comment  string `json:"total_comment"`
	Ip             string `json:"ip"`
	Img_chk        string `json:"img_chk"`
	Recommend_chk  string `json:"recommend_chk"`
	Winnerta_chk   string `json:"winnerta_chk"`
	Voice_chk      string `json:"voice_chk"`
	Hit            string `json:"hit"`
	Write_type     string `json:"write_type"`
	User_id        string `json:"user_id"`
	Prev_link      string `json:"prev_link"`
	Prev_subject   string `json:"prev_subject"`
	Headtitle      string `json:"headtitle"`
	Next_link      string `json:"next_link"`
	Next_subject   string `json:"next_subject"`
	Best_chk       string `json:"best_chk"`
	Realtime_l_chk string `json:"realtime_l_chk"`
	IsNotice       string `json:"isNotice"`
	Date_time      string `json:"date_time"`
	Alarm_flag     int    `json:"alarm_flag"`
	Is_minor       bool   `json:"is_minor"`
}

type PostViewMain struct {
	Memo             string `json:"memo"`
	Recommend        string `json:"recommend"`
	Recommend_member string `json:"recommend_member"`
	Nonrecommend     string `json:"nonrecommend"`
	Nonrecomm_user   bool   `json:"nonrecomm_user"`
}

func (s *Session) GetPost(gallid string, gno int) (Post, error) {
	// url is https://app.dcinside.com/api/gall_view_new.php?id=tsmanga&no=1&app_id=T0RtOWkzbFRhVEJndnExU3hmMC80QTV1WVgzQ21SNHdxRS9jRjRocDJUVT0%3D&client_id=eGTqnqzsSzSKYCSWs7LJ8j%3AAPA91bGCO-S2Y5IRfBlK9rWqYGBMcWc15ynPo6nDz7RczKnfURdbkYldx1-7F-sXcrFCdBD86kWqNFTGfnH2-rWmPnnBD3nU6SAtRoVSu3bZ_DwJgG4nmvHc824BGAiB49U-Aq8XXnlx7
	urld := Base64EncodeLink(fmt.Sprintf("https://app.dcinside.com/api/gall_view_new.php?id=%s&no=%d&app_id=%s&client_id=eGTqnqzsSzSKYCSWs7LJ8j:APA91bGCO-S2Y5IRfBlK9rWqYGBMcWc15ynPo6nDz7RczKnfURdbkYldx1-7F-sXcrFCdBD86kWqNFTGfnH2-rWmPnnBD3nU6SAtRoVSu3bZ_DwJgG4nmvHc824BGAiB49U-Aq8XXnlx", gallid, gno, s.Appid))
	req, err := http.NewRequest("GET", urld, nil)
	if err != nil {
		return Post{}, errors.New("Error Making Request")
	}
	req.Header.Set("User-Agent", "dcinside.app")
	req.Host = "app.dcinside.com"
	req.Header.Set("referer", "https://app.dcinside.com")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return Post{}, errors.New("Error Posting Request")
	}
	bod, _ := io.ReadAll(res.Body)
	// fmt.Println(string(bod))
	if gjson.Get(string(bod), "0.result").Exists() == true {
		return Post{}, errors.New("Please refresh your appid")
	}
	var post []Post
	err = json.Unmarshal(bod, &post)
	if err != nil {
		return Post{}, errors.New("Error while parsing json")
	}
	return post[0], nil
}

func (s *Session) DelPost(gallid string, gno int, pw string) (bool, error) {
	rr := url.Values{}
	rr.Add("id", gallid)
	rr.Add("no", strconv.Itoa(gno))
	rr.Add("write_pw", pw)
	rr.Add("app_id", s.Appid)
	rr.Add("mode", "board_del")
	rr.Add("client_token", "eGTqnqzsSzSKYCSWs7LJ8j:APA91bGCO-S2Y5IRfBlK9rWqYGBMcWc15ynPo6nDz7RczKnfURdbkYldx1-7F-sXcrFCdBD86kWqNFTGfnH2-rWmPnnBD3nU6SAtRoVSu3bZ_DwJgG4nmvHc824BGAiB49U-Aq8XXnlx")
	req, err := http.NewRequest(
		"POST",
		"https://app.dcinside.com/api/gall_del.php",
		strings.NewReader(rr.Encode()),
	)

	if err != nil {
		return false, errors.New("Error Posting Request")
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("user-agent", "dcinside.app")
	req.Header.Set("Host", "app.dcinside.com")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return false, errors.New("Error Posting Request")
	}
	bod, _ := io.ReadAll(res.Body)
	if gjson.Get(string(bod), "0.result").Exists() == true {
		return false, errors.New("Please refresh your appid")
	}
	return gjson.Get(string(bod), "result").Bool(), nil
}
