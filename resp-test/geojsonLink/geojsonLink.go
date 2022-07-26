package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"
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
		NodeID    string  `json:"NODE_ID,omitempty"`
		NodeType  string  `json:"NODE_TYPE,omitempty"`
		NodeName  string  `json:"NODE_NAME,omitempty"`
		TurnP     string  `json:"TURN_P,omitempty"`
		Connect   string  `json:"CONNECT,omitempty"`
		FNode     string  `json:"F_NODE,omitempty"`
		Lanes     int     `json:"LANES,omitempty"`
		Length    float64 `json:"LENGTH,omitempty"`
		LinkID    string  `json:"LINK_ID,omitempty"`
		MaxSpd    int     `json:"MAX_SPD,omitempty"`
		MultiLink string  `json:"MULTI_LINK,omitempty"`
		RoadName  string  `json:"ROAD_NAME,omitempty"`
		RoadNo    string  `json:"ROAD_NO,omitempty"`
		RoadRank  string  `json:"ROAD_RANK,omitempty"`
		RoadType  string  `json:"ROAD_TYPE,omitempty"`
		RoadUse   string  `json:"ROAD_USE,omitempty"`
		TNode     string  `json:"T_NODE,omitempty"`
	} `json:"properties"`
	Geometry struct {
		Type        string        `json:"type"`
		Coordinates [][][]float64 `json:"coordinates"`
	} `json:"geometry"`
}

var b bytes.Buffer

func main() {
	fmt.Println("start")

	start := time.Now()
	data, _ := os.Open("moct_node.geojson")
	byteValue, _ := ioutil.ReadAll(data)

	var geojsonFormat GeojsonFormat
	json.Unmarshal(byteValue, &geojsonFormat)

	elapsed := time.Since(start)
	fmt.Printf("unmarshal complete. %v \n", elapsed)

	for _, f := range geojsonFormat.Features {
		key := ""
		value, err := json.Marshal(f)
		if err != nil {
			fmt.Println(err)
			return
		}

		switch geojsonFormat.Name {
		case "node":
			key = "node:" + f.Properties.NodeID
		case "link":
			key = "link:" + f.Properties.LinkID
		}

		genRedisProto("set", key, string(value))
	}

	elapsed = time.Since(start)
	fmt.Printf("remarshal complete. %v \n", elapsed)

	// fmt.Println(output)
	// 처리 시간 측정

	fmt.Println("end")
	fmt.Println("write file")

	ioutil.WriteFile("moct-node.db", []byte(b.String()), 0644)
}

func genRedisProto(cmd string, key string, value string) {
	str := "" +
		"*3\r\n" +
		"$" + strconv.Itoa(len(cmd)) + "\r\n" +
		cmd + "\r\n" +
		"$" + strconv.Itoa(len([]byte(key))) + "\r\n" +
		key + "\r\n" +
		"$" + strconv.Itoa(len([]byte(value))) + "\r\n" +
		value + "\r\n"

	b.WriteString(str)
}
