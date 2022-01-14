package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func main() {

	// Watcher
	go func() {
		watcherClient1, err := clientv3.New(clientv3.Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
		})
		watcherClient2, err := clientv3.New(clientv3.Config{
			Endpoints:   []string{"localhost:22379"},
			DialTimeout: 5 * time.Second,
		})
		watcherClient3, err := clientv3.New(clientv3.Config{
			Endpoints:   []string{"localhost:32379"},
			DialTimeout: 5 * time.Second,
		})
		if err != nil {
			log.Fatal(err)
		}
		defer watcherClient1.Close()
		defer watcherClient2.Close()
		defer watcherClient3.Close()

		rch1 := watcherClient1.Watch(context.Background(), "Akey")
		for wresp := range rch1 {
			for _, ev := range wresp.Events {
				fmt.Printf("Watch 1 - %s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			}
		}
		rch2 := watcherClient2.Watch(context.Background(), "Akey")
		for wresp := range rch2 {
			for _, ev := range wresp.Events {
				fmt.Printf("Watch 2 - %s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			}
		}
		rch3 := watcherClient3.Watch(context.Background(), "Akey")
		for wresp := range rch3 {
			for _, ev := range wresp.Events {
				fmt.Printf("Watch 3 - %s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			}
		}
	}()

	// Adder client
	adderClient, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	adderClient.Close()

	_, err = adderClient.Put(context.Background(), "Akey", strconv.Itoa(0))

	go func() {
		for i := 0; i < 50; i++ {
			getResp, _ := adderClient.Get(context.Background(), "Akey")
			respVal, _ := strconv.Atoi(string(getResp.Kvs[0].Value))
			respVal++
			fmt.Printf("Adder: %v\n", respVal)
			adderClient.Put(context.Background(), "Akey", strconv.Itoa(respVal))
			time.Sleep(1 * time.Second)
		}
	}()

}
