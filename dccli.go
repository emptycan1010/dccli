package dccli

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
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
	Appid      string
	Apptoken   string
	NowGallID  string
	NowPostNo  int
}

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
	Notify_recent string    `json:"notify_recent"`
}

type HeadTXT struct {
	No       string `json:"no"`
	Name     string `json:"name"`
	Level    string `json:"level"`
	Selected bool   `json:"selected"`
}

type Comment struct {
	Total_comment string        `json:"total_comment"`
	Total_page    string        `json:"total_page"`
	Re_page       string        `json:"re_page"`
	Comment_list  []CommentList `json:"comment_list"`
}

type CommentList struct {
	Is_delete_flag   string `json:"is_delete_flag"` // 이 댓글은 게시물 작성자가 삭제하였습니다.
	Member_icon      string `json:"member_icon"`
	IpData           string `json:"ip_data"`
	Name             string `json:"name"`
	User_id          string `json:"user_id"`
	Comment_memo     string `json:"comment_memo"`
	Dccon            string `json:"dccon"`
	Dccon_detail_idx string `json:"dccon_detail_idx"`
	Comment_no       string `json:"comment_no"`
	Date_time        string `json:"date_time"`
	Del_scope        string `json:"del_scope"`
}

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

func (s *Session) GetAppID() error {
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
	h.Write([]byte(fmt.Sprintf("dcArdchk_%s", Appc[0].Date))) // value token calculated
	res, err = http.PostForm(
		"https://msign.dcinside.com/auth/mobile_app_verification",
		url.Values{
			"value_token":  {fmt.Sprintf("%x", h.Sum(nil))},
			"signature":    {"ReOo4u96nnv8Njd7707KpYiIVYQ3FlcKHDJE046Pg6s="},
			"client_token": {"hangus"},
		},
	)
	if err != nil {
		return errors.New("Error GetAppID function")
	}
	bod, _ = io.ReadAll(res.Body)
	s.Appid = gjson.Get(string(bod), "app_id").String()
	return nil
}

func (s *Session) AddComment(gallid string, gno int, datgeul string, writer string, pw string) (bool, error) {
	rr := url.Values{}
	rr.Add("id", gallid)
	rr.Add("no", strconv.Itoa(gno))
	rr.Add("comment_nick", writer)
	rr.Add("comment_pw", pw)
	rr.Add("client_token", "eGTqnqzsSzSKYCSWs7LJ8j:APA91bGCO-S2Y5IRfBlK9rWqYGBMcWc15ynPo6nDz7RczKnfURdbkYldx1-7F-sXcrFCdBD86kWqNFTGfnH2-rWmPnnBD3nU6SAtRoVSu3bZ_DwJgG4nmvHc824BGAiB49U-Aq8XXnlx")
	rr.Add("app_id", s.Appid)
	rr.Add("mode", "com_write")
	rr.Add("comment_memo", datgeul)
	req, err := http.NewRequest(
		"POST",
		"https://app.dcinside.com/api/comment_ok.php",
		strings.NewReader(rr.Encode()),
	)
	//fmt.Println(rr.Encode())

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
	return gjson.Get(string(bod), "0.result").Bool(), nil
}

func (s *Session) GetComment(gallid string, gno int, commentpage int) (Comment, error) {
	urld := Base64EncodeLink(fmt.Sprintf("https://app.dcinside.com/api/comment_new.php?id=%s&no=%d&app_id=%s&re_page=%d", gallid, gno, s.Appid, commentpage))
	req, err := http.NewRequest("GET", urld, nil)
	if err != nil {
		return Comment{}, errors.New("Error Posting Request")
	}
	req.Header.Set("User-Agent", "dcinside.app")
	req.Host = "app.dcinside.com"
	req.Header.Set("referer", "https://app.dcinside.com")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return Comment{}, errors.New("Error Posting Request")
	}
	var commentlist []Comment
	bod, _ := io.ReadAll(res.Body)
	// fmt.Println(string(bod))
	if gjson.Get(string(bod), "0.result").Exists() == true {
		return Comment{}, errors.New("Please refresh your appid")
	}
	err = json.Unmarshal(bod, &commentlist)
	if err != nil {
		return Comment{}, errors.New("Error while parsing json")
	}
	return commentlist[0], nil
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

func (s *Session) Login(id string, pw string) error {
	if s.isLoggedin == true {
		return errors.New("Already logged in")
	}
	rr := url.Values{}
	rr.Add("user_id", id)
	rr.Add("user_pw", pw)
	rr.Add("client_token", "eGTqnqzsSzSKYCSWs7LJ8j:APA91bGCO-S2Y5IRfBlK9rWqYGBMcWc15ynPo6nDz7RczKnfURdbkYldx1-7F-sXcrFCdBD86kWqNFTGfnH2-rWmPnnBD3nU6SAtRoVSu3bZ_DwJgG4nmvHc824BGAiB49U-Aq8XXnlx")
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
	json.Unmarshal(bod, &account)
	// fmt.Println(account)
	if account.Result == true {
		s.isLoggedin = true
		s.Account = account
		return nil
	} else {
		return errors.New("Failed to Log in")
	}
} // 객체지향 추가 완?료

func (s *Session) DelComment(gallid string, gno int, commentno int, pw string) (bool, error) {
	rr := url.Values{}
	rr.Add("id", gallid)
	rr.Add("no", strconv.Itoa(gno))
	rr.Add("comment_no", strconv.Itoa(commentno))
	rr.Add("comment_pw", pw)
	rr.Add("app_id", s.Appid)
	rr.Add("mode", "comment_del")
	rr.Add("client_token", "eGTqnqzsSzSKYCSWs7LJ8j:APA91bGCO-S2Y5IRfBlK9rWqYGBMcWc15ynPo6nDz7RczKnfURdbkYldx1-7F-sXcrFCdBD86kWqNFTGfnH2-rWmPnnBD3nU6SAtRoVSu3bZ_DwJgG4nmvHc824BGAiB49U-Aq8XXnlx")
	req, err := http.NewRequest(
		"POST",
		"https://app.dcinside.com/api/comment_del.php",
		strings.NewReader(rr.Encode()),
	)
	if err != nil {
		return false, errors.New("Error making post Request")
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("user-agent", "dcinside.app")
	req.Header.Set("Host", "app.dcinside.com")
	req.Header.Set("referer", "https://app.dcinside.com")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return false, errors.New("Error Posting Request")
	}
	bod, _ := io.ReadAll(res.Body)

	if gjson.Get(string(bod), "0.cause").String() == "certification" {
		return false, errors.New("Please refresh your appid")
	}
	return gjson.Get(string(bod), "0.result").Bool(), nil
}

func New() *Session {
	p := &Session{}
	p.isLoggedin = false
	return p
}

//func (s *Session) FetchFCMToken() {
//	r, e := http.NewRequest("POST", "https://firebaseinstallations.googleapis.com/v1/projects/dcinside-b3f40/installations", nil)
//	if e != nil {
//		panic(e)
//	}
//
//	r.Header.Set("accept", "application/json")
//	r.Header.Set("accept-encoding", "gzip")
//	r.Header.Set("cache-control", "no-cache")
//	r.Header.Set("connection", "Keep-Alive")
//	r.Header.Set("content-encoding", "gzip")
//	r.Header.Set("host", "firebaseinstallations.googleapis.com")
//	r.Header.Set("user-agent", "Dalvik/2.1.0 (Linux; U; Android 13; Pixel 5 Build/TP1A.221105.002)")
//	r.Header.Set("x-android-cert", "43BD70DFC365EC1749F0424D28174DA44EE7659D")
//	r.Header.Set("x-android-package", "com.dcinside.app.android")
//	r.Header.Set("x-firebase-client", "H4sIAAAAAAAAAKtWykhNLCpJSk0sKVayio7VUSpLLSrOzM9TslIyUqoFAFyivEQfAAAA")
//	r.Header.Set("x-goog-api-key", "AIzaSyDcbVof_4Bi2GwJ1H8NjSwSTaMPPZeCE38")
//	b := bytes.NewBuffer([]byte(`{
//  "fid": "f7RXAqYIR6iACLGVP06qb4",
//  "appId": "1:477369754343:android:d2ffdd960120a207727842",
//  "authVersion": "FIS_v2",
//  "sdkVersion": "a:17.0.2"}`))
//	r.Body = io.NopCloser(b)
//	client := &http.Client{}
//	res, err := client.Do(r)
//	if err != nil {
//		panic(err)
//	}
//	bod, _ := io.ReadAll(res.Body)
//	fmt.Println(string(bod)) // Must get fid, appid,
//}

// 위에꺼 FCM토큰관련한건데 아직 작동도안하고 수정할거많아서 일단 주석처리함
