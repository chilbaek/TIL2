package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

type GeojsonFormat struct {
	Type string `json:"type"`
	Name string `json:"name"`
	Crs  struct {
		Type       string `json:"type"`
		Properties struct {
			Name string `json:"name"`
		} `json:"properties"`
	} `json:"crs"`
	Features []Feature `json:"features"`
}

type Feature struct {
	Type       string `json:"type"`
	Properties struct {
		NodeID   string `json:"NODE_ID"`
		NodeType string `json:"NODE_TYPE"`
		NodeName string `json:"NODE_NAME"`
		TurnP    string `json:"TURN_P"`
	} `json:"properties"`
	Geometry struct {
		Type        string    `json:"type"`
		Coordinates []float64 `json:"coordinates"`
	} `json:"geometry"`
}

var output = ""

func main() {
	data, _ := os.Open("moct_node.geojson")
	byteValue, _ := ioutil.ReadAll(data)

	var geojsonFormat GeojsonFormat

	json.Unmarshal(byteValue, &geojsonFormat)

	fmt.Println("start")

	for _, f := range geojsonFormat.Features {
		b, _ := json.Marshal(f)
		key := "node:" + f.Properties.NodeID
		value := b
		genRedisProto("set", key, string(value))
	}

	fmt.Println("end")
	fmt.Println("write file")

	ioutil.WriteFile("mmap-test.db", []byte(output), 0644)
}

func genRedisProto(cmd string, key string, value string) {
	str := "" +
		"*3\r\n" +
		"$" + strconv.Itoa(len(cmd)) + "\r\n" +
		cmd + "\r\n" +
		"$" + strconv.Itoa(len([]byte(key))) + "\r\n" +
		key + "\r\n" +
		"$" + strconv.Itoa(len([]byte(value))) + "\r\n" + // +4 ??
		value + "\r\n"

	output = output + str
}
