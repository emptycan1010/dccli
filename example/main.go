package main

import (
	"fmt"
	"github.com/emptycan1010/dccli"
	"github.com/tidwall/gjson"
)

func main() {
	r := dccli.GetGallList("weatherbaby", gjson.Get(dccli.GetAppID(), "app_id").String())
	for i := 0; i < len(r.GallList); i++ {
		fmt.Println(r.GallList[i].Subject)
	}
}
