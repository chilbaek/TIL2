package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

func main() {
	var name = flag.String("name", "", "give a name")
	flag.Parse()

	// Create a etcd client
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})

	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	// Create a sesssion to elect a leader
	s, err := concurrency.NewSession(cli)
	if err != nil {
		log.Fatal(err)
	}
	defer s.Close()

	e := concurrency.NewElection(s, "/leader-election")
	ctx := context.Background()
	fmt.Println(e.Key())

	// Elect a leader (or wait that the leader resign)
	if err := e.Campaign(ctx, "e"); err != nil {
		log.Fatal(err)
	}
	fmt.Println("leader election for ", *name)

	fmt.Println("Do some work in", *name)
	time.Sleep(5 * time.Second)

	if err := e.Resign(ctx); err != nil {
		log.Fatal(err)
	}
	fmt.Println("resign ", *name)
	fmt.Println(e.Key())
	fmt.Println(e)

	fmt.Println(e.Leader(ctx))

}
