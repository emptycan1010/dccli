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
	fmt.Println(GetGallList("onii", gjson.Get(GetAppID(), "app_id").String()))

}

func GetGallList(gallid string, appid string) string {
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
	return string(bod)
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
