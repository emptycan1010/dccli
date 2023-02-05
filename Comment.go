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
