package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"io"
	"log"
	"net/http"
	"net/url"
)

type AppCheckstruct struct {
	Result        bool   `json:"result"`
	Ver           string `json:"ver"`
	Notice        bool   `json:"notice"`
	Notice_update bool   `json:"notice_update"`
	Date          string `json:"date"`
}

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
	notify_recent string    `json:"notify_recent"`
}

type HeadTXT struct {
	No       string `json:"no"`
	Name     string `json:"name"`
	Level    string `json:"level"`
	Selected bool   `json:"selected"`
}

func main() {
	//fmt.Println(HashedURLmake("weatherbaby", gjson.Get(GetAppID(), "app_id").String()))
	//http.Get(HashedURLmake("weatherbaby", gjson.Get(GetAppID(), "app_id").String()))
	//req, err := http.NewRequest("GET", HashedURLmake("weatherbaby", gjson.Get(GetAppID(), "app_id").String()), nil)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//req.Header.Set("User-Agent", "dcinside.app")
	//req.Host = "app.dcinside.com"
	//req.Header.Set("referer", "https://app.dcinside.com")
	//client := &http.Client{}
	//res, err := client.Do(req)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//bod, _ := io.ReadAll(res.Body)
	//fmt.Println(string(bod))

	//fmt.Printf(GetGallList("onii", gjson.Get(GetAppID(), "app_id").String()))
	r := GetGallList("weatherbaby", gjson.Get(GetAppID(), "app_id").String())
	for _, v := range r {
		fmt.Println(v.GallList)
	}
}

func GetGallList(gallid string, appid string) []Getgalldata {
	req, err := http.NewRequest("GET", HashedURLmake(gallid, appid), nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("User-Agent", "dcinside.app")
	req.Host = "app.dcinside.com"
	req.Header.Set("referer", "https://app.dcinside.com")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	bod, _ := io.ReadAll(res.Body)
	var gg []Getgalldata
	err = json.Unmarshal(bod, &gg)
	if err != nil {
		log.Fatal(err)
	}
	return gg
}

func HashedURLmake(gallid string, appid string) string {
	input := []byte(
		fmt.Sprintf("https://app.dcinside.com/api/gall_list_new.php?id=%s&page=1&app_id=%s",
			gallid,
			appid,
		),
	)
	return fmt.Sprintf("https://app.dcinside.com/api/redirect.php?hash=%s", base64.StdEncoding.EncodeToString(input))
}

func GetAppID() string {
	res, err := http.Get("http://json2.dcinside.com/json0/app_check_A_rina.php")
	if err != nil {
		log.Fatal(err)
	}

	bod, _ := io.ReadAll(res.Body)
	//fmt.Println(string(bod))
	var Appc []AppCheckstruct
	err = json.Unmarshal(bod, &Appc)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(Appc[0].Date)
	//fmt.Sprintf("dcArdchk_%s", Appc[0].Date)
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("dcArdchk_%s", Appc[0].Date))) // value token calculated
	res, err = http.PostForm(
		"https://msign.dcinside.com/auth/mobile_app_verification",
		url.Values{
			"value_token":  {fmt.Sprintf("%x", h.Sum(nil))},
			"signature":    {"5rJxRKJ2YLHgBgj6RdMZBl2X0KcftUuMoXVug0bsKd0="},
			"client_token": {"hangus"},
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	bod, _ = io.ReadAll(res.Body)
	return string(bod)
}
