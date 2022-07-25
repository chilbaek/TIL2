package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Info struct {
	Type     string
	Name     string
	Crs      string
	Features []string
}

func main() {
	data, err := os.Open("moct_node.geojson")
	if err != nil {

	}
	byteValue, _ := ioutil.ReadAll(data)

	var db_info Info

	json.Unmarshal(byteValue, &db_info)

	fmt.Println()
	fmt.Println(db_info.Crs)
}
