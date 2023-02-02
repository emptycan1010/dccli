package main

import (
	"fmt"
	"github.com/emptycan1010/dccli"
	"github.com/tidwall/gjson"
)

func main() {
	appid := gjson.Get(dccli.GetAppID(), "app_id").String()
	r := dccli.GetGallList("weatherbaby", appid)
	for i := 0; i < len(r.GallList); i++ {
		fmt.Println(r.GallList[i].Subject)
	}
	print(dccli.AddComment("tsmanga", appid, 1, "aaa", "ㅇㅇ", "0000"))
}
