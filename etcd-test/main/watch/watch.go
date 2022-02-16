package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// go run watch.go -key=foo
func main() {
	key := flag.String("key", "", "Key to watch")
	flag.Parse()

	go func() {
		watcherClient, err := clientv3.New(clientv3.Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
		})
		if err != nil {
			log.Fatal(err)
		}
		defer watcherClient.Close()

		rch := watcherClient.Watch(context.Background(), *key)
		for wresp := range rch {
			for _, ev := range wresp.Events {
				fmt.Printf("Watcher - %s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			}
		}

		// rch2 := watcherClient.Watch(context.Background(), "fooForLease2")
		// for wresp := range rch2 {
		// 	for _, ev := range wresp.Events {
		// 		fmt.Printf("Watcher - %s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
		// 	}
		// }
	}()

	var ch chan bool
	<-ch // Block forever

}
