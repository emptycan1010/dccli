package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
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
	fmt.Println(GetAppID())
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
	//fmt.Println(string(bod))
	return string(bod)
}
