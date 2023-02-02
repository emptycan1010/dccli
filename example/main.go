package main

import (
	"github.com/emptycan1010/dccli"
	"github.com/tidwall/gjson"
)

func main() {
	appid := gjson.Get(dccli.GetAppID(), "app_id").String()
	print(dccli.AddComment("tsmanga", appid, 1, "aaa", "ㅇㅇ", "1111"))
}
