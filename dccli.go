package dccli

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
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
	Notify_recent string    `json:"notify_recent"`
}

type HeadTXT struct {
	No       string `json:"no"`
	Name     string `json:"name"`
	Level    string `json:"level"`
	Selected bool   `json:"selected"`
}

//func main() {
//	r := GetGallList("weatherbaby", gjson.Get(GetAppID(), "app_id").String())
//
//	for i := 0; i < len(r.GallList); i++ {
//		fmt.Println(r.GallList[i].Subject)
//	}
//}

func GetGallList(gallid string, appid string) (Getgalldata, error) {
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
		return Getgalldata{}, errors.New("Error Posting Request")
	}
	bod, _ := io.ReadAll(res.Body)
	var gg []Getgalldata
	err = json.Unmarshal(bod, &gg)
	if err != nil {
		return Getgalldata{}, errors.New("Error Unmarshaling Json")
	}
	return gg[0], nil
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

func GetAppID() (string, error) {
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
		return "", errors.New("Error GetAppID function")
	}
	bod, _ = io.ReadAll(res.Body)
	return string(bod), nil
}

func AddComment(gallid string, appid string, gno int, datgeul string, writer string, pw string) (bool, error) {
	rr := url.Values{}
	rr.Add("id", gallid)
	rr.Add("no", strconv.Itoa(gno))
	rr.Add("comment_nick", writer)
	rr.Add("comment_pw", pw)
	rr.Add("client_token", "eGTqnqzsSzSKYCSWs7LJ8j:APA91bGCO-S2Y5IRfBlK9rWqYGBMcWc15ynPo6nDz7RczKnfURdbkYldx1-7F-sXcrFCdBD86kWqNFTGfnH2-rWmPnnBD3nU6SAtRoVSu3bZ_DwJgG4nmvHc824BGAiB49U-Aq8XXnlx")
	rr.Add("app_id", appid)
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
	return gjson.Get(string(bod), "0.result").Bool(), nil
}

func GetComment(gallid string, appid string, gno int, commentpage int) (Comment, error) {
	urld := Base64EncodeLink(fmt.Sprintf("https://app.dcinside.com/api/comment_new.php?id=%s&no=%d&app_id=%s&re_page=%d", gallid, gno, appid, commentpage))
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
	err = json.Unmarshal(bod, &commentlist)
	if err != nil {
		return Comment{}, errors.New("Error while parsing json")
	}
	return commentlist[0], nil
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

type Comment struct {
	Total_comment string        `json:"total_comment"`
	Total_page    string        `json:"total_page"`
	Re_page       string        `json:"re_page"`
	Comment_list  []CommentList `json:"comment_list"`
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

func GetPost(gallid string, appid string, gno int) (Post, error) {
	// url is https://app.dcinside.com/api/gall_view_new.php?id=tsmanga&no=1&app_id=T0RtOWkzbFRhVEJndnExU3hmMC80QTV1WVgzQ21SNHdxRS9jRjRocDJUVT0%3D&client_id=eGTqnqzsSzSKYCSWs7LJ8j%3AAPA91bGCO-S2Y5IRfBlK9rWqYGBMcWc15ynPo6nDz7RczKnfURdbkYldx1-7F-sXcrFCdBD86kWqNFTGfnH2-rWmPnnBD3nU6SAtRoVSu3bZ_DwJgG4nmvHc824BGAiB49U-Aq8XXnlx7
	urld := Base64EncodeLink(fmt.Sprintf("https://app.dcinside.com/api/gall_view_new.php?id=%s&no=%d&app_id=%s&client_id=eGTqnqzsSzSKYCSWs7LJ8j:APA91bGCO-S2Y5IRfBlK9rWqYGBMcWc15ynPo6nDz7RczKnfURdbkYldx1-7F-sXcrFCdBD86kWqNFTGfnH2-rWmPnnBD3nU6SAtRoVSu3bZ_DwJgG4nmvHc824BGAiB49U-Aq8XXnlx", gallid, gno, appid))
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
	var post []Post
	err = json.Unmarshal(bod, &post)
	if err != nil {
		return Post{}, errors.New("Error while parsing json")
	}
	return post[0], nil
}
