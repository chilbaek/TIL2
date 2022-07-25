package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	auth := `{
    "username" : "admin",
    "password" : "[password입력]",
    "hostname" : "[host정보입력]",
    "port" : "27017"
}`

	type Info struct {
		Username string
		Password string
		Hostname string
		Port     string
	}

	var info Info
	json.Unmarshal([]byte(auth), &info)
	fmt.Println(info.Port)

}
