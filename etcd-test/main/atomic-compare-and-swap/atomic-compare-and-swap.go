package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strconv"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// Atomic Comapre-And-Swap(CAS)을 활용하여 동시성 보장 확인
func main() {
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

	_, err = adderClient.Put(context.Background(), "foo", strconv.Itoa(50))

	requestTimeoout := 10 * time.Second
	ctx1, _ := context.WithTimeout(context.Background(), requestTimeoout)

	go func() {
		for i := 0; i < 50; i++ {

			// Try compare-and-swap until succeeded
			for {
				resp1, _ := adderClient.Get(context.Background(), "foo")
				value1 := string(resp1.Kvs[0].Value)
				num1, _ := strconv.Atoi(value1)
				num1++
				fmt.Printf("[Adder] Value1(%v), num1(%v)\n", value1, num1)

				// Compare-and-Swap (CAS)
				txResp1, err1 := adderClient.Txn(ctx1).
					If(clientv3.Compare(clientv3.Value("foo"), "=", value1)).
					Then(clientv3.OpPut("foo", strconv.Itoa(num1))).
					Else(clientv3.OpGet("foo")).
					Commit()

				if err1 != nil {
					fmt.Printf("[Adder] i(%v) - Err1: %v\n", i, err1)
					// cancel1()
				}
				fmt.Printf("[Adder] i(%v) - txResp1: %v\n", i, txResp1)

				if txResp1.Succeeded {
					break
				}
				// time.Sleep(1 * time.Millisecond)
			}
			// time.Sleep(1 * time.Millisecond)
		}
	}()

	// Substarcor Section
	substractorClient, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer substractorClient.Close()

	ctx2, _ := context.WithTimeout(context.Background(), requestTimeoout)

	go func() {
		for j := 0; j < 50; j++ {

			// Try compare-and-swap until succeeded
			for {
				resp2, _ := substractorClient.Get(context.Background(), "foo")
				value2 := string(resp2.Kvs[0].Value)
				num2, _ := strconv.Atoi(value2)
				num2--
				fmt.Printf("[Substractor] value2(%v), num2(%v)\n", value2, num2)

				// Compare-and-Swap (CAS)
				txResp2, err2 := substractorClient.Txn(ctx2).
					If(clientv3.Compare(clientv3.Value("foo"), "=", value2)).
					Then(clientv3.OpPut("foo", strconv.Itoa(num2))).
					Else(clientv3.OpGet("foo")).
					Commit()

				if err2 != nil {
					fmt.Printf("[Substractor] j(%v) - Err2: %v\n", j, err2)
				}

				fmt.Printf("[Substractor] j(%v) - txResp2: %v\n", j, txResp2)

				if txResp2.Succeeded {
					break
				}
				// time.Sleep(1 * time.Millisecond)
			}
			// time.Sleep(1 * time.Millisecond)
		}
	}()

	var ch chan bool
	<-ch // Block forever

}
