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
	d.NoLogID = "ㅇㅇ"
	d.NoLogPW = "1111"
	d.GetAppID()
	memo := []dccli.MemoBlock{}
	memo = append(memo, dccli.MemoBlock{Content: "<div>test</div>"})
	fmt.Println(d.RequestPost("tsmanga", memo, "테스트"))
}
