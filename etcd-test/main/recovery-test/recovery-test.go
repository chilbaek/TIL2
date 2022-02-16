package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func main() {

	// Adder section
	adderClient, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer adderClient.Close()

	for i := 0; i < 100; i++ {
		fmt.Printf("%d ", i)
		adderClient.Put(adderClient.Ctx(), "rcvr", strconv.Itoa(i))
		time.Sleep(2 * time.Second)
	}
}
