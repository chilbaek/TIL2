// example from https://www.pixelstech.net/article/1615108646-Service-discovery-with-etcd

package main

import (
	"context"
	"log"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// ServiceRegister create and register service
type ServiceRegister struct {
	cli           *clientv3.Client
	leaseID       clientv3.LeaseID
	keepAliveChan <-chan *clientv3.LeaseKeepAliveResponse
	key           string
	val           string
}

// NewServiceRegister create new service
func NewServiceRegister(endpoints []string, key, val string, lease int64) (*ServiceRegister, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}

	ser := &ServiceRegister{
		cli: cli,
		key: key,
		val: val,
	}

	if err := ser.putKeyWithLease(lease); err != nil {
		return nil, err
	}

	return ser, nil
}

// Set lease
func (s *ServiceRegister) putKeyWithLease(lease int64) error {
	// Grant lease time
	resp, err := s.cli.Grant(context.Background(), lease)
	if err != nil {
		return err
	}

	// Register service and bind lease
	_, err = s.cli.Put(context.Background(), s.key, s.val, clientv3.WithLease(resp.ID))
	if err != nil {
		return err
	}

	// Set keep-alivd logic
	leaseRespChan, err := s.cli.KeepAlive(context.Background(), resp.ID)
	if err != nil {
		return err
	}
	s.leaseID = resp.ID
	s.keepAliveChan = leaseRespChan
	log.Println(s.leaseID)
	log.Printf("Put key:%s   val:%s   success!", s.key, s.val)
	return nil
}

// ListenLeaseRespChan listen and watch
func (s *ServiceRegister) ListenLeaseRespChan() {
	for leaseKeepResp := range s.keepAliveChan {
		log.Println("Success", leaseKeepResp)
	}
	log.Println("ListenLease closed.")
}

// Close
func (s *ServiceRegister) Close() error {
	// Revoke lease
	if _, err := s.cli.Revoke(context.Background(), s.leaseID); err != nil {
		return err
	}
	log.Println("Revoke lease")
	return s.cli.Close()
}

func main() {
	var endpoints = []string{"localhost:2379"}
	ser, err := NewServiceRegister(endpoints, "/web/node2", "value2-2379", 5)
	if err != nil {
		log.Fatal(err)
	}

	// Listen to keep alive chan
	go ser.ListenLeaseRespChan()
	select {
	// case <-time.After(20 * time.Second):
	// 	ser.Close()
	}
}
