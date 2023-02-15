package main

import (
	"fmt"

	dccli "github.com/emptycan1010/dcgo"
)

func main() {
	// r, e := dccli.GetAppID()
	// if e != nil {
	// 	log.Fatalln(e)
	// }
	// appid := gjson.Get(r, "app_id").String()
	// fmt.Println(appid)
	// print(dccli.AddComment("tsmanga", appid, 1, "aaa", "ㅇㅇ", "1111"))
	// res, e := dccli.GetComment("tsmanga", appid, 1, 1)
	// if e != nil {
	// 	log.Fatalln(e)
	// }
	// fmt.Println(res)
	// r, _ := dccli.Login("adfasdfasdf", "111@")
	// print(dccli.DelComment("tsmanga", appid, 1, 39, "1111"))
	d := dccli.New()
	// d.NoLogID = "ㅇㅇ"
	// d.NoLogPW = "1111"
	// d.GetAppID()
	// memo := []dccli.MemoBlock{}
	// memo = append(memo, dccli.MemoBlock{Content: "<div>test</div>"})
	// fmt.Println(d.RequestPost("tsmanga", "가나마", memo))
	d.GetAppID()
	d.FetchFCMToken()
	d.NoLogID = "ㅇㅇ"
	d.NoLogPW = "1111"
	// fmt.Println(d.FCM)
	// fmt.Println(d.Appid)
	// fmt.Println(d.FCM.Token)
	fmt.Println(d.AddComment("tsmanga", 1, "ㅇㅇ", "ㅇㅇ", "1111"))
	// d.RequestPost("tsmanga", "asdf", []dccli.MemoBlock{{Content: "<div>xtx</div>"}})
}
