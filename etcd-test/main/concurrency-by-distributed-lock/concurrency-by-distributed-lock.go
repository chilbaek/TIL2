package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strconv"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

// Adder Vs. Substractor - Race Condition
// 분산락 Distributed Lock 활용하여 동시성 보장 확인
// 더하기, 빼기 똑같이 50번 반복: 분산락 걸때 <-> 안 걸때 결과값 다름
func main() {
	var name = flag.String("name", "foo", "Give a name")
	flag.Parse()

	// Watcher Section
	go func() {
		watcherClient, err := clientv3.New(clientv3.Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
		})
		if err != nil {
			log.Fatal(err)
		}
		defer watcherClient.Close()

		rch := watcherClient.Watch(context.Background(), "foo")
		for wresp := range rch {
			for _, ev := range wresp.Events {
				fmt.Printf("Watch - %s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			}
		}

	}()

	// Adder Section
	adderClient, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer adderClient.Close()

	// Create a sessions to acquire a lock
	s1, _ := concurrency.NewSession(adderClient)
	defer s1.Close()

	l1 := concurrency.NewMutex(s1, "/distributed-lock/")
	ctx1 := context.Background()

	adderClient.Put(context.Background(), "foo", strconv.Itoa(50))

	go func() {
		for i := 0; i < 50; i++ {
			// Acquire lock (or wait to have it)
			if err := l1.Lock(ctx1); err != nil {
				log.Fatal(err)
			}
			fmt.Println("Acquire lock for ", *name)

			// Do value+1 and put the value
			resp1, _ := adderClient.Get(context.Background(), "foo")
			num1, _ := strconv.Atoi(string(resp1.Kvs[0].Value))
			num1++
			fmt.Printf("Adder: %v\n", num1)
			adderClient.Put(context.Background(), "foo", strconv.Itoa(num1))
			// time.Sleep(10 * time.Millisecond)

			// Release lock
			if err := l1.Unlock(ctx1); err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Released lock for %s \n\n", *name)
		}
	}()

	// Substractor Section
	substractorClient, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer substractorClient.Close()

	// Create a sesssions to acquire a lock
	s2, _ := concurrency.NewSession(substractorClient)
	defer s2.Close()

	l2 := concurrency.NewMutex(s2, "/distributed-lock/")
	ctx2 := context.Background()

	go func() {
		for j := 0; j < 50; j++ {
			// Acquire lock (or wait to have it)
			if err := l2.Lock(ctx2); err != nil {
				log.Fatal(err)
			}
			fmt.Println("Acquired lock for ", *name)

			// Do value-1 and put the value
			resp2, _ := substractorClient.Get(context.Background(), "foo")
			num2, _ := strconv.Atoi(string(resp2.Kvs[0].Value))
			num2--
			fmt.Printf("Substractor: %v\n", num2)
			substractorClient.Put(context.Background(), "foo", strconv.Itoa(num2))
			// time.Sleep(10 * time.Millisecond)

			// Release Rock
			if err := l2.Unlock(ctx2); err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Releases lock for %s \n\n", *name)
		}
	}()

	var ch chan bool
	<-ch // Block forever

}
