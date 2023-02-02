package main

import (
	"log"

	dccli "github.com/emptycan1010/dcgo"
	"github.com/tidwall/gjson"
)

func main() {
	r, e := dccli.GetAppID()
	if e != nil {
		log.Fatalln(e)
	}
	appid := gjson.Get(r, "app_id").String()
	print(dccli.AddComment("tsmanga", appid, 1, "aaa", "ㅇㅇ", "1111"))
}
