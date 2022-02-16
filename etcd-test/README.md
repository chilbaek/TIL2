# ETCD client examples

## ETCD 클러스터 구성

1. etcd 설치 ([공식문서 가이드](https://etcd.io/docs/v3.5/install/))

   ```bash
   git clone https://github.com/etcd-io/etcd.git
   cd etcd
   ./build
   ```

   

2. 클러스터 구성 (goreman 활용)

   1. [goreman 설치](https://github.com/mattn/goreman)

      ```bash
      go install github.com/mattn/goreman@latest
      ```

      

   2. Procfile 설정

      > * etcd1: localhost:2379
      > * etcd2: localhost:22379
      > * etcd3: localhost:32379

      위 내용 기준으로 [예제 Procfile](https://github.com/etcd-io/etcd/blob/main/Procfile) 과 동일하게 설정

      

   3. 클러스터 기동
      (Procfile이 있는 위치에서 실행)

      ```bash
      goreman start
      ```

      

   4. etcdctl command로 정상구동 확인

   5. * 리더/멤버 상태 확인

        ```bash
        etcdctl endpoint status --cluster -w table
        ```

      * 헬스 체크

        ```bash
        etcdctl --endpoints=127.0.0.1:2379,127.0.0.1:22379,127.0.0.1:32379 endpoint health -w table
        ```

      * 모든 키 확인

        ```bash
        etcdctl --endpoints=localhost:2379 get --prefix --keys-only ''
        ```

        



## Go client v3 examples

```bash
go get go.etcd.io/etcd/client/v3
```



### 1. Basic (Get & Put)

클라이언트가 ETCD 클러스터 KV저장소에 읽기/쓰기(Get/Put)를 수행

#### basic/basic.go

1. ETCD 클러스터 정보를 정의하여 클라이언트 선언

   ```go
   cli, err := clientv3.New(clientv3.Config{
   	Endpoints:   []string{"http://localhost:2379"},
   	DialTimeout: 5 * time.Second,
   })
   ```

2. Put

   ```go
   _, err = cli.Put(cli.Ctx(), "sample_keys", "sample_value")

3. Get

   ```go
   resp, err2 := cli2.Get(cli2.Ctx(), "sample_keys")
   fmt.Println(resp.Header)
   fmt.Println(resp.Kvs)
   ```



### 2. Watch API

Key의 변경을 주시하는 API

#### watch/watch.go

```go
// watcher 클라이언트 선언
watcherClient, err := clientv3.New(clientv3.Config{
	Endpoints:   []string{"localhost:2379"},
	DialTimeout: 5 * time.Second,
})
defer watcherClient.Close()

// key값 변경 이벤트 발생시 콘솔 출력
rch := watcherClient.Watch(context.Background(), *key)
for wresp := range rch {
	for _, ev := range wresp.Events {
		fmt.Printf("Watcher - %s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
	}
}
```

[스크린샷]

키의 값이 변경될 때 마다 콘솔에 출력된다.

### 3. Lease API

* Lease : 클라이언트 상태를 감지하기 위한 메커니즘 
* 클러스터는 TTL(Time-To-Live)이 있는 Lease를 부여
* 클러스터가 주어진 TTL 안에 keepAlive를 수신하지 못하면 리스는 만료되어 해당 클라이언트가 생성한 KV를 삭제한다.

#### lease/lease.go

```go
// Lease Section
cli, err := clientv3.New(clientv3.Config{
	Endpoints:   []string{"localhost:2379"},
	DialTimeout: 5 * time.Second,
})
defer cli.Close()

// minimum lease TTL is 5-second
resp, err := cli.Grant(context.TODO(), 5)

// after 5 seconds, the key 'fooForLease' will be removed
_, err = cli.Put(context.TODO(), "fooForLease", "bar", clientv3.WithLease(

println(resp.ID)
// renew the lease only once
ka, kaerr := cli.KeepAliveOnce(context.TODO(), resp.ID)
fmt.Println("ttl: ", ka.TTL)
```



### 4. Service Register

* 서비스들은 서로를 인지할 수 있도록 공통된 공간에 자신의 정보를 등록해야 한다.
* 서비스가 시작되면 서비스는 서버 주소를 etcd에 기록하고 lease를 부여하며, health check를 통해 lease 생명주기를 유지해야 한다.

#### sr/service-register.go

서비스가 시작될 때 etcd에 자신의 정보(K-V)를 저장하고 lease(생명주기) 속성을 부여받는다. 서비스가 죽고 keep-alive Time-out이 발생하면 etcd에 저장된 해당 서비스 정보는 삭제된다.

```go
type ServiceRegister struct {
	cli           *clientv3.Client	// etcd client
	leaseID       clientv3.LeaseID	// lease id
	keepAliveChan <-chan *clientv3.LeaseKeepAliveResponse
	key           string	// service key
	val           string	// service value
}

// Register new service
func NewServiceRegister(endpoints []string, key, val string, lease int64) (*ServiceRegister, error) {
	cli, err := clientv3.New(clientv3.Config{	// Init client
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})

	ser := &ServiceRegister{ 
		cli: cli,
		key: key,
		val: val,
	}

  // Put service information
	if err := ser.putKeyWithLease(lease); err != nil { 
		return nil, err
	}
	return ser, nil
}

// Set lease
func (s *ServiceRegister) putKeyWithLease(lease int64) error {
	// Grant lease time
	resp, err := s.cli.Grant(context.Background(), lease)

	// Register service and bind lease
	_, err = s.cli.Put(context.Background(), s.key, s.val, clientv3.WithLease(resp.ID))

	// Set keep-alive logic
	leaseRespChan, err := s.cli.KeepAlive(context.Background(), resp.ID)
	s.leaseID = resp.ID
	s.keepAliveChan = leaseRespChan
	log.Println(s.leaseID)
	log.Printf("Put key:%s   val:%s   success!", s.key, s.val)
	return nil
}

// '/web' 하위에 서비스(node1) 등록
func main() {
	var endpoints = []string{"localhost:2379"}
	ser, err := NewServiceRegister(endpoints, "/web/node1", "value-2379", 5)
	...
}
```



### 5. Service Discovery

* 서비스는 특정 서비스 이름으로 클라이언트를 초기화하여 서비스 검색 프로세스를 시작한다.
* Watch API를 활용하여 서버 추가, 삭제 등을 수행한다.

#### sd/service-discovery.go

```go
type ServiceDiscovery struct {
	cli        *clientv3.Client  // etcd client
	serverList map[string]string // server list
	lock       sync.Mutex
}

// Service discovery client
func NewServiceDiscovery(endpoints []string) *ServiceDiscovery {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})

	return &ServiceDiscovery{
		cli:        cli,
		serverList: make(map[string]string),
	}
}

func (s *ServiceDiscovery) WatchService(prefix string) error {
	resp, err := s.cli.Get(context.Background(), prefix, clientv3.WithPrefix())
	for _, ev := range resp.Kvs {
		s.SetServiceList(string(ev.Key), string(ev.Value))
	}

	go s.watcher(prefix)
	return nil
}

// Watch API를 활용하여 '/web1' 하위에 생성/삭제되는 노드를 주시한다
func main() {
	var endpoints = []string{"localhost:2379"}
	ser := NewServiceDiscovery(endpoints)
	defer ser.Close()
	ser.WatchService("/web/")
	for {
		select {
		case <-time.Tick(10 * time.Second):
			log.Println(ser.GetServices())
		}
	}
}

```



### 6. Distributed Lock

#### concurrency-by-distributed-lock/concurrency-by-distributed-lock.go

* Adder Vs. Substractor - Race Condition
* 분산락 Distributed Lock 활용하여 동시성 보장 확인
* 더하기, 빼기 똑같이 50번 반복: 분산락 걸때 <-> 안 걸때 결과값 다름을 확인

