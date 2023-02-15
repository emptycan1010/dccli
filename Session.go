package dccli

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/emptycan1010/dcgo/checkin"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
)

type Session struct {
	Account    Account
	isLoggedin bool
	NoLogID    string
	NoLogPW    string
	Appid      string
	Apptoken   string
	NowGallID  string
	NowPostNo  int
	FCM        AccountFCM
}

type AppCheckstruct struct {
	Result        bool   `json:"result"`
	Ver           string `json:"ver"`
	Notice        bool   `json:"notice"`
	Notice_update bool   `json:"notice_update"`
	Date          string `json:"date"`
}

type Account struct {
	Result           bool   `json:"result"`
	User_id          string `json:"user_id"`
	User_no          string `json:"user_no"`
	Name             string `json:"name"`
	Is_adult         string `json:"is_adult"`
	Is_dormancy      int    `json:"is_dormancy"`
	Otp_token        string `json:"otp_token"`
	Is_gonick        int    `json:"is_gonick"`
	Is_security_code string `json:"is_security_code"`
	Auth_change      int    `json:"auth_change"`
	Stype            string `json:"stype"`
	Pw_campaign      int    `json:"pw_campaign"`
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

func Base64EncodeLink(input string) string {
	return fmt.Sprintf("https://app.dcinside.com/api/redirect.php?hash=%s", base64.StdEncoding.EncodeToString([]byte(input)))
}

func (s *Session) GetAppID() error {
	//{
	//	"fid": "fT-9GN8ASwOa9ihWpuokdn",
	//	"appId": "1:477369754343:android:d2ffdd960120a207727842",
	//	"authVersion": "FIS_v2",
	//	"sdkVersion": "a:17.0.2"
	//}
	//{
	//	"name": "projects/477369754343/installations/fT-9GN8ASwOa9ihWpuokdn",
	//	"fid": "fT-9GN8ASwOa9ihWpuokdn",
	//	"refreshToken": "3_AS3qfwKZ1zsz4C0dvSZdg9CBYSKG4MBEoYrNuKiGg-908_yTBGRkxTD1qeI_vuzCOGb5GSj8O8cxdCwWFXT0fEBlPNEmAkbPV5ZFOVRg-yQKojU",
	//	"authToken": {
	//	"token": "eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJhcHBJZCI6IjE6NDc3MzY5NzU0MzQzOmFuZHJvaWQ6ZDJmZmRkOTYwMTIwYTIwNzcyNzg0MiIsImV4cCI6MTY3NjIwOTYxOCwiZmlkIjoiZlQtOUdOOEFTd09hOWloV3B1b2tkbiIsInByb2plY3ROdW1iZXIiOjQ3NzM2OTc1NDM0M30.AB2LPV8wRgIhALHo8OYiKb41UxwuCyjPLJ21qQQM2Ofme63jdbQc0YzHAiEAvbRYIf13I0NqMmHBe5iRz7-Hglcx0-RfCf0sOi8XWnw",
	//		"expiresIn": "604800s"
	//}
	//token=fT-9GN8ASwOa9ihWpuokdn:APA91bHW2DbvpDTeJxUA_ACwoLzPkCfJpWqj5N2Eb9H7gYz9D28e1jJH_RRXZoDDMKClZSlXXVosI10BlHGcFgOg1dkkJRm8qCaU9Fci7V2q9ZSRSefw0tA7xW1A_3jl8UU5GG3_uLNL
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
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("dcArdchk_%s", Appc[0].Date)))
	req, err := http.NewRequest(
		"POST",
		"https://msign.dcinside.com/auth/mobile_app_verification",
		strings.NewReader(url.Values{
			"vName":        {"4.7.5"},
			"vCode":        {"100028"},
			"pkg":          {"com.dcinside.app"},
			"value_token":  {fmt.Sprintf("%x", h.Sum(nil))},
			"signature":    {"ReOo4u96nnv8Njd7707KpYiIVYQ3FlcKHDJE046Pg6s="},
			"client_token": {s.FCM.Token},
		}.Encode()),
	)
	if err != nil {
		return errors.New("Error Making Request")
	}
	// req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "dcinside.app")
	req.Header.Set("accept-encoding", "gzip")
	req.Header.Set("referer", "http://www.dcinside.com")
	req.Header.Set("user-agent", "dcinside.app")
	client := &http.Client{}
	res, err = client.Do(req)
	if err != nil {
		return errors.New("Error Posting Request")
	}
	gre, err := gzip.NewReader(res.Body)
	if err != nil {
		return errors.New("Error gzip")
	}
	bod, _ = io.ReadAll(gre)
	fmt.Println(string(bod))
	s.Appid = gjson.Get(string(bod), "app_id").String()
	if gjson.Get(string(bod), "result").Bool() == false {
		return errors.New("Error GetAppID function")
	}
	return nil
}

func (s *Session) Login(id string, pw string) error {
	if s.isLoggedin == true {
		return errors.New("Already logged in")
	}
	rr := url.Values{}
	rr.Add("user_id", id)
	rr.Add("user_pw", pw)
	rr.Add("client_token", s.FCM.Token)
	rr.Add("mode", "login_normal")

	req, err := http.NewRequest(
		"POST",
		"https://msign.dcinside.com/api/login",
		strings.NewReader(rr.Encode()),
	)
	if err != nil {
		return errors.New("Error Making Request")
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("user-agent", "dcinside.app")
	req.Header.Set("Host", "msign.dcinside.com")
	req.Header.Set("referer", "http://www.dcinside.com")
	client := &http.Client{}

	res, err := client.Do(req)

	if err != nil {
		return errors.New("Error Posting Request")
	}
	bod, _ := io.ReadAll(res.Body)
	// fmt.Println(string(bod))
	var account Account
	e := json.Unmarshal(bod, &account)
	if e != nil {
		return errors.New("Error while parsing json")
	}
	// fmt.Println(account)
	if account.Result == true {
		s.isLoggedin = true
		s.Account = account
		return nil
	} else {
		return errors.New("Failed to Log in")
	}
} // 객체지향 추가 완?료

func New() *Session {
	p := &Session{}
	p.isLoggedin = false
	return p
}

func (s *Session) FetchFCMToken() {
	// rr := url.Values{}
	// rr.Add("fid", "")
	// rr.Add("refreshToken", "")
	// rr.Add("appId", "1:477369754343:android:1f4e2da7c458e2a7")
	// rr.Add("authVersion", "FIS_v2")
	// rr.Add("sdkVersion", "a:17.0.2")
	// // I must encode rr into Gzip

	r, e := http.NewRequest("POST", "https://firebaseinstallations.googleapis.com/v1/projects/dcinside-b3f40/installations", nil)
	if e != nil {
		panic(e)
	}
	r.Header.Set("accept", "application/json")
	r.Header.Set("content-type", "application/json")
	//r.Header.Set("content-encoding", "gzip")
	// r.Header.Set("accept-encoding", "gzip")
	r.Header.Set("host", "firebaseinstallations.googleapis.com")
	r.Header.Set("user-agent", "Dalvik/2.1.0 (Linux; U; Android 13; Pixel 5 Build/TP1A.221105.002)")
	r.Header.Set("x-android-cert", "E6DA04787492CDBD34C77F31B890A3FAA3682D44")
	r.Header.Set("x-android-package", "com.dcinside.app")
	r.Header.Set("x-firebase-client", "H4sIAAAAAAAAAKtWykhNLCpJSk0sKVayio7VUSpLLSrOzM9TslIyUqoFAFyivEQfAAAA")
	r.Header.Set("x-goog-api-key", "AIzaSyDcbVof_4Bi2GwJ1H8NjSwSTaMPPZeCE38")
	b := bytes.NewBuffer(
		[]byte(
			`{
		"fid": "",
		"appId": "1:477369754343:android:d2ffdd960120a207727842",
		"authVersion": "FIS_v2",
		"sdkVersion": "a:17.0.2"}`,
		),
	)
	r.Body = io.NopCloser(b)

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		panic(err)
	}
	bod, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	// fmt.Println(string(bod))
	var accountFCM AccountFCM
	e = json.Unmarshal(bod, &accountFCM)
	if e != nil {
		panic(e)
	}
	s.FCM = accountFCM
	// fmt.Println(s.FCM)
	// 이제 토큰가져오면됨 ㅋㅋ

	var andid int64 = 0
	var fingerprint = "google/razor/flo:7.1.1/NMF26Q/1602158:user/release-keys"
	var hdw = "flo"
	radio := "FLO-04.04"
	clid := "android-google"
	sdkver := int32(25)
	loc := "ko"
	tz := "KST"
	var lcms int64 = 0
	var zero int32 = 0
	var three int32 = 0
	checkinReq := checkin.CheckinRequest{
		TimeZone:         &tz,
		AndroidId:        &andid,
		Locale:           &loc,
		Version:          &three,
		OtaCert:          []string{"--no-output--"},
		MacAddress:       []string{"02", "00", "00", "00", "00", "00"},
		Fragment:         &zero,
		UserSerialNumber: &zero, // 0

		//CheckinRequest_Checkin
		Checkin: &checkin.CheckinRequest_Checkin{
			Build: &checkin.CheckinRequest_Checkin_Build{
				Fingerprint: &fingerprint,
				Hardware:    &hdw,
				Radio:       &radio,
				ClientId:    &clid,
				SdkVersion:  &sdkver,
			},
			LastCheckinMs: &lcms,
		},
	}

	//checkinReq.String()

	r, e = http.NewRequest("POST", "https://android.clients.google.com/checkin", bytes.NewBufferString(checkinReq.String()))
	r.Header.Set("Content-Type", "application/x-protobuf")
	r.Header.Set("User-Agent", "Android-Checkin/3.0")

	res, err = client.Do(r)
	if err != nil {
		panic(err)
	}
	bod, err = io.ReadAll(res.Body)
	fmt.Println(string(bod))
	rr := url.Values{}
	rr.Add("X-subtype", "477369754343")
	rr.Add("sender", "477369754343")
	rr.Add("X-app_ver", "4.7.5")
	rr.Add("X-appid", gjson.Get(string(bod), "fid").String())
	rr.Add("X-scope", "*")
	rr.Add("X-Goog-Firebase-Installations-Auth", gjson.Get(string(bod), "authToken.token").String())
	rr.Add("X-gmp_app_id", "1:477369754343:android:1f4e2da7c458e2a7")
	rr.Add("X-firebase-app-name-hash", "R1dAH9Ui7M-ynoznwBdw01tLxhI")
	rr.Add("X-app_ver_name", "100028")
	rr.Add("app", "com.dcinside.app")
	rr.Add("device", strconv.FormatInt(andid, 10))
	rr.Add("app_ver", "4.7.5")
	rr.Add("gcm_ver", "221215022")
	rr.Add("cert", "E6DA04787492CDBD34C77F31B890A3FAA3682D44")
	r, e = http.NewRequest("POST", "https://android.apis.google.com/c2dm/register3", strings.NewReader(rr.Encode()))
	if e != nil {
		panic(e)
	}
	r.Header.Set("authorization", "AidLogin 3966377448498170683:2982263657081238075")

	fmt.Sprintf("AidLogin %s:%s", "", "")

	r.Header.Set("host", "android.apis.google.com")
	r.Header.Set("app", "com.dcinside.app")
	client = &http.Client{}
	res, err = client.Do(r)
	if err != nil {
		panic(err)
	}
	bod, err = io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	// fmt.Println(string(bod))
	// fmt.Println(string(bod))
	s.FCM.Token = string(bod)[6:]
}

type AccountFCM struct {
	Name         string `json:"name"`
	Fid          string `json:"fid"`
	RefreshToken string `json:"refreshToken"`
	AuthToken    struct {
		Token     string `json:"token"`
		ExpiresIn string `json:"expiresIn"`
	} `json:"authToken"`
	Token string
}
