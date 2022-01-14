package main

import (
	"context"
	"fmt"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func main() {
	// Expect dial timeout on ipv4 blackhole
	_, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"http://localhost:2373"},
		DialTimeout: 2 * time.Second,
	})

	// etcd clientv3 >= v3.2.10, grpc/grpc-go >= v1.7.3
	if err == context.DeadlineExceeded {
		// Handle errors
	}

	cli2, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"http://localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		// Handle errors
	}

	_, err = cli2.Put(cli2.Ctx(), "sample_keys", "sample_value")
	if err != nil {
		// Handle errors
	}

	resp, err2 := cli2.Get(cli2.Ctx(), "sample_keys")
	if err2 != nil {
		// Handle errors
	}

	fmt.Println(resp.Header)
	fmt.Println(resp.Kvs)

}
