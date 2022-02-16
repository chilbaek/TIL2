# About etcd

## etcd 

"/etc" + "distributed" : 거대 분산 시스템의 configuration 정보를 저장

etcd는 일관된 분산 KV저장소. 주로 분산 시스템에서 별도의 조정 서비스로 사용된다. 그리고 메모리에 완전히 들어갈 수 있는 적은 양의 데이터를 보유하도록 설계되었다.

* 클라이언트는 굳이 리더를 콕 집어 요청을 보내지 않아도 된다. 클러스터 멤버 누구에게라도 요청하면 그 요청은 리더로 전달된다.
* etcd는 디스크에 데이터를 쓴다. [하드웨어 스펙 권장 사항](https://etcd.io/docs/v3.5/op-guide/hardware/)
* 클러스터 구성시 홀수개의 노드, 7개 이하로 구성할 것을 권장
* `snapshot`명령을 통해 백업 가능 ([more](https://etcd.io/docs/v3.5/op-guide/recovery/#snapshotting-the-keyspace))

## ZooKeeper Vs. Etcd

|                            | ZooKeeper                                                    | Etcd                                                         |
| -------------------------- | ------------------------------------------------------------ | ------------------------------------------------------------ |
| Foundation                 | Apache Hadoop                                                | Google Kubernetes                                            |
| Based                      | Java                                                         | Go                                                           |
| Extension                  | Apache Curator                                               | Etcd API (http curl 명령으로도 가능)                         |
| Leader Election Algorithm  | [ZAB(Zookeeper Atomic Broadcast)](https://cwiki.apache.org/confluence/display/ZOOKEEPER/Zab1.0) <br />가장 최신 기록을 가진 리더를 선출 | RAFT                                                         |
| Storage                    | Znode (Tree)                                                 | Key-Value                                                    |
| Data Storage               | Memory                                                       | Disk                                                         |
| In memory                  | Data                                                         | in-memory index<br />페이지 캐시                             |
| Linearizable Reads         | X                                                            | [O](https://etcd.io/docs/v3.3/learning/api_guarantees/#linearizability) |
| Data version               | znode 데이터의 버전만 기록                                   | [버전에 따른 값도 기록](https://etcd.io/docs/v3.3/learning/data_model/) |
| HTTP/JSON API              | X                                                            | O                                                            |
| Membership Reconfiguration | O (> 3.5)                                                    | O                                                            |
| Maximum reliable db size   | 수백 MB                                                      | 수 GB                                                        |

Zookeeper를 개조하여 만든 etcd의 장점

* Dynamic cluster membership reconfiguration
* Stable read/write under high load
* A multi-version concurrency control data model
* Reliable key monitoring which never which never silently drop events
* Lease primitives decoupling connections from sessions
  세션에서 연결을 분리하는 lease 처리
* APIs for safe distributed shared locks

![image-20220109234254039](md-images/image-20220109234254039.png)

* https://etcd.io/docs/v3.3/learning/why/
* https://dzone.com/articles/apache-zookeeper-vs-etcd3
* https://www.slideshare.net/CoreOS_Slides/etcd-mission-critical-keyvalue-store



### Ephemeral znode

# 데이터 모델

* etcd는 드물게 업데이트 되는 데이터 저장에 적합하다.
* key-value 값의 모든 버전을 기록한다. (히스토리) 
* key를 생성하면 해당 키의 버전은 1부터 시작하며, 키를 수정할 때 마다 버전이 증가한다. 삭제될 경우엔 0으로 초기화 된다.
* 공간 절약을 위해 오래된 버전은 압축되고 압축개정 이전의 버전은 삭제된다.
* b+tree 구조를 사용하여 물리적 데이터를 저장한다.
* btree 인덱스를 보조메모리에 사용하여 쿼리 속도를 높인다.

# 클라이언트 디자인

## 용어

* clientv3: etcd v3 API용 etcd 공식 Go client
* clientv3-grpc1.0: etcd v3.1에서 사용되는 공식 클라이언트
* clientv3-grpc1.7: etcd v3.2 
* clientv3-grpc1.23: etcd v3.4
* Balancer: 재시도 및 장애조치 메커니즘을 구현하는 etcd 로드 밸런서. 클라이언트는 엔드포인트 간 로드 균형을 자동으로 조정한다.
* Endpoints: 클라이언트가 연결할 수 있는 etcd 서버 목록. 일반적으로 etcd 클라이언트는 3개 또는 5개.

# 클라이언트 개요

## clientv3-grpc1.0

<img src="https://etcd.io/docs/v3.4/learning/img/client-balancer-figure-01.png" alt="클라이언트 밸런서 그림 01.png" style="zoom:50%;" />

* A, B, C 모든 TCP 연결을 유지하며, 그 중 하나만 선택하여 요청을 보낸다.
* 장애 발생시 빠른 조치가 가능하나 더 많은 리소스를 필요로 한다. 

## clientv3-grpc1.7

<img src="https://etcd.io/docs/v3.5/learning/img/client-balancer-figure-03.png" alt="클라이언트 밸런서 그림-03.png" style="zoom: 50%;" />

* grpc1.0과 달리 하나의 서버와 TCP연결을 유지한다.
* 에러 핸들러가 따로 있다.
* 로드 밸런서는 Endpoint의 상태 체크를 위해 주기적으로 Ping을 보낸다.

<img src="https://etcd.io/docs/v3.5/learning/img/client-balancer-figure-07.png" alt="클라이언트 밸런서 그림 07.png" style="zoom:50%;" />

* 밸런서는 위 그림과 같이 헬스 체크하던 A가 나머지 B, C와 연결이 끊어진 경우도 감지할 수 있다.
* keepAlive Ping time-out, Context time-out 시간차를 비교해서 판단하는 듯

<img src="https://etcd.io/docs/v3.5/learning/img/client-balancer-figure-08.png" alt="클라이언트 밸런서-그림-08.png" style="zoom:50%;" />

## clientv3-grpc1.23

<img src="https://etcd.io/docs/v3.5/learning/img/client-balancer-figure-08.png" alt="client-balancer-figure-08.png" style="zoom:50%;" />

# [etcd3 API](https://etcd.io/docs/v3.5/learning/api/)

## gRPC 서비스

모든 API 요청은 gRPC 원격 프로시저 호출이다. etcd의 키 공간을 다루는 데 중요한 서비스는 다음과 같다.

* KV - Key-Value 페어를 Create, Update, Fetch, Delete
* Watch - Key의 변경 사항을 모니터링 
* Lease - 클라이언트의 keep-alive 메시지 consume 목적 (keepAlive(TTL) 속성 부여)

클러스터 자체를 관리하는 서비스에는 다음이 포함된다.

* Auth - 사용자 인증 목적의 역할 기반 인증 메커니즘
* Cluster - 멤버십 정보 및 구성 기능 제공
* Maintance - 복구 스냅샷 생성, 저장소 최적화, 구성원별 상태 정보 반환

## Key-Value API

### Key ranges



## Put

키는 `Put` 을 수행하는 호출을 실행하여 KV 저장소에 저장된다.

```protobuf
message PutRequest {
	bytes key = 1;
  bytes value = 2;
  int64 lease = 3;
  bool prev_kv = 4;
  bool ignore_value = 5;
  bool ignore_lease = 6;
}
```

* lease : KV 저장소의 키와 연결하기 위한 lease ID. 값이 0이면 lease가 없다는 뜻
* prev_kv : 이 값이 설정되면, 해당 Put 요청을 업데이트 하기 전의 K-V 데이터로 응답한다.
* ignore_value : 이 값이 설정되면, 현재 값을 변경하지 않고 키를 업데이트 한다. 키가 없으면 오류 반환.
* ignore_lease : 이 값이 설정되면, lease를 변경하지 않고 키를 업데이트 한다. 키가 없으면 오류 반환.

```protobuf
message PutResponse {
  ResponseHeader header = 1;
  mvccpb.KeyValue prev_kv = 2;
}
```

* PutRequest의 `prev_kv` 가 설정되었다면 KV쌍은 `Put`에 의해 덮어 씌여졌다.

## Delete Range

## Transaction

## Watch API

Watch API는 키의 변화를 비동기로 모니터링하기 위한 이벤트 기반 인터페이스를 제공한다. Etcd3 watch 는 현재 또는 과거의 주어진 수정본으로 부터의 키 변경을 계속 주시하고, 키 업데이트 내용을 클라이언트에 전달한다.

### Events

모든 키에 대한 모든 변화는 `Event` 메시지로 표현된다. 하나의 이벤트 메시지는 업데이트의 데이터와 유형을 포함한다.

```protobuf
message Event {
  enum EventType {
    PUT = 0;
    DELETE = 1;
  }
  EventType type = 1;
  KeyValue kv = 2;
  KeyValue prev_kv = 3;
}
```

* Type - 이벤트 유형.
  * PUT 타입은 키에 새로운 데이터가 저장된다. 
  * DELETE 타입은 키가 삭제된다.
* KV - 키밸류는 이벤트와 연관되어 있다. 
  * PUT 이벤트는 현재의 kv 쌍을 포함한다. kv.Version=1 인 PUT 이벤트는 키의 생성을 나타낸다. 
  * DELETE 이벤트는 삭제로 설정된 삭제된 키를 포함한다. (?)
* Prev_KV - 이벤트 직전 수정본의 키에 대한 KV 쌍. 대역폭을 저장하기 위해 watch가 명시적으로 활성화한 경우에만 값을 정의한다.

### Watch streams

워치들은 장기 실행 요청이며 스트림 이벤트 데이터에 gRPC 스트림을 사용한다. 워치 스트림은 양방향이다. 클라이언트는 워치를 설정하기 위해 스트림에 쓰고, 워치 이벤트를 수신하기 위해 읽는다. 싱글 워치 스트림은 워치당(per-watch) 식별자로 이벤트에 태그를 지정하여 별개의 워치를 다중화할 수 있다. 이 다중화(multplexing)은 코어 etcd 클러스터에서 메모리 공간과 연결 오버헤드를 줄이는 데에 도움된다.

워치는 이벤트에 대해 세 가지를 보장한다.

* Ordered - 이벤트는 리비전revision 순으로 정렬된다. 이벤트가 이미 게시된 시간에 이벤트보다 선행된 경우 워치에 이벤트가 표시되지 않는다?
* Reliable - 일련의 이벤트가 이벤트의 하위 시퀀스를 삭제하지 않는다. a < b < c 순으로 이벤트가 정렬된 경우 워치가 a와 c 이벤트를 받았다면, b를 받는 것이 보장된다.
* Atomic - 이벤트 목록은 완전한 리비전을 포함하도록 보장된다. 다중 키에 걸친 똑같은 리비전은 여러 이벤트 목록으로 분할되지 않는다.

클라이언트 `WatchCreateRequest`는 `Watch` 다음에서 반환된 스트림을 통해 워치를 생성한다.

```protobuf
message WatchCreateRequest {
  bytes key = 1;
  bytes range_end = 2;
  int64 start_revision = 3;
  bool progress_notify = 4;
  
  enum FilterType {
    NOPUT = 0;
    NODELETE = 1;
  }
  repeated FilterType filters = 5;
  bool prev_kv = 6;
}
```

* Key, Range_End - 워치를 수행할 키 범위
* Start_Revision - 포괄적으로 워치를 시작할 위치에 대한 옵셔널 리비전. 지정하지 않으면 워치 생성 응답헤더 리비전에 따라 이벤트를 스트리밍한다. 사용가능한 전체 이벤트 기록은 마지막 압축 리비전부터 볼 수 있다.
* Progress_Notify - 설정되면 최근 이벤트가 없는 경우 워치 이벤트 없이 WatchResponse 를 주기적으로 수신한다. 클라이언트가 최근 알려진 리비전에서 시작하여 연결이 끊긴 워쳐를 복구하는 경우 유용하다. Etcd 서버는 현재 서버 로드에 기반하여 얼마나 자주 노티를 보낼지 결정한다.
* Filters - 서버 측에서 필터링할 이벤트 유형 목록
* Prev_Kv - 설정되면 워치는 이벤트가 발생하기 전의 KV 데이터를 수신한다. 덮어쓴 데이터를 알고자 할 때 유용하다.



`WatchCreateRequest`에 대한 응답 또는 생성된 워치를 위한 새로운 이벤트가 있다면 클라이언트는 `WatchResponse`를 받는다.



## Lease API

Lease란 클라이언트 상태를 감지하기 위한 메커니즘이다. 클러스터는 TTL이 있는 리스를 부여한다. etcd 클러스터가 주어진 TTL 안에 keepAlive를 수신하지 못하면 리스는 만료된다.

리스를 kv저장소에 연결하기 위해, 각 키는 하나의 리스와 붙을 수 있다. 리스가 만료 또는 취소되면 리스에 붙은 모든 키들은 삭제된다. 각 만료된 키는 이벤트 히스토리에 삭제 이벤트를 생성한다.

### Optaining leases

리스는`LeaseGrantRequest`를 갖는 `LeaseGrant` API호출을 통해 얻을 수 있다.

```protobuf
message LeaseGrantRequest {
	int64 TTL = 1;
	int64 ID = 2;
}
```

* TTL - 초 단위
* ID - 리스를 위한 ID. 0으로 설정하면 etcd는 ID를 선택한다. (임의부여한다는 말인가?)



```protobuf
message LeaseGrantResponse {
	ResponseHeader header = 1;
	int64 ID = 2;
	int64 TTL = 3;
}
```

* ID - 부여된 리스의 ID

* TTL - 리스에 대해 서버가 선택한 TTL (초 단위)

  

```protobuf
message LeaseRevokeRequest {
	int64 ID = 1;
}
```

* ID - 취소할 리스ID. 리스가 취소되면 여기 붙은 모든 키들은 삭제된다.



### Keep alives

리스들은 `LeaseKeepAlive` API 호출과 함께 생성된 양방향 스트림을 사용하면서 리프레쉬 된다. 클라이언트가 리스를 리프레쉬하고 싶을때 클라이언트는 스트림을 통해 `LeaseKeepAliveRequest` 를 보낸다.

```protobuf
message LeaseKeepAliveRequest {
	int64 ID = 1;
}
```

```protobuf
message LeaseKeepAliveResponse {
	ResponseHeaser header = 1;
	int64 ID = 2;
	int64 TTL = 3;
}
```







# Developer Guide

## [API refence: concurrency](https://etcd.io/docs/v3.5/dev-guide/api_concurrency_reference_v3/#service-lock-serveretcdserverapiv3lockv3lockpbv3lockproto)

API 레퍼런스는 `.proto` 파일에 의해 자동 생성된다.

> 분산 락은 데이터베이스 등 공통된 저장소를 이용하여 자원이 사용 중인지를 체크한다. 그래서 전체 서버에서 동기화된 처리가 가능하다.

> 여러 독립된  프로세스에서 하나의 자원을 공유해야 할 때, 데이터에 결함이 발생하지 않도록 하기 위해 분산 락을 활용할 수 있다. 분산 락을 구현하려면 데이터베이스 등 여러 프로세스가 공통으로 사용하는 저장소를 활용해야 한다. 

### service Lock (server/etcdserver/api/v3lock/v3lockpb/v3lock.proto)

락 서비스는 클라이언트 쪽 락 기능을 gRPC 인터페이스로 제공한다.

| Method | Request Type  | Response Type  | Description                                                  |
| ------ | ------------- | -------------- | ------------------------------------------------------------ |
| Lock   | LockRequest   | LockResponse   | 락은 지정된 락에서 분산된 공유 락을 얻는다. 성공하면 호출자가 락을 유지하는 한 존재하는 고유 키를 반환한다. 이 키는 락 소유권을 보유하는 동안에만 etcd 업데이트가 발생하도록 안전하게 트랜잭션과 함께 사용할 수 있다. 락은 키에 대한 잠금해제가 호출되거나 소유자와의 lease 연결이 만료될 때 까지 유지된다. |
| Unlock | UnlockRequest | UnlockResponse | 언락은 락에 의해 반환받은 키를 갖는다. 그러면 락을 기다리는 다음 락 호출자가 깨워지고 락에 대한 소유권을 갖는다. |

* message LockRequest
* message LockResponse
* message UnlockRequest
* message UnlockResponse

### service Election (server/etcdserver/api/v3election/v3electionpb/v3election.proto)

일렉션 서비스는 클라이언트 쪽 일렉션 기능을 gRPC 인터페이스로 제공한다.

| Method   | Request Type    | Response Type    | Description                                                  |
| -------- | --------------- | ---------------- | ------------------------------------------------------------ |
| Campaign | CampaignRequest | CampaignResponse | 캠페인은 election에서 리더십 획득을 기다리고, 성공하면 리더십을 나타내는 리더키LeaderKey를 반환한다. 그 다음 리더키는 election에서 새로운 값을 발행하며, 아직까지 리더십에 들어오는 API 요청을 트랜잭션 방식으로 보호하고 election을 그만둔다. |
| Proclaim | ProclaimRequest | ProclaimResponse | Proclaim공표는 리더의 게시된 값을 새로운 값으로 업데이트한다. |
| Leader   | LeaderRequest   | LeaderResponse   | 리더는 현재 election 선언문을 반환한다. (있는 경우)          |
| Observe  | LeaderRequest   | LeaderResponse   | Observe는 선출된 리더에 의해 만들어진 election 선언문을 stream한다. |
| Resign   | ResignRequest   | ResignResponse   | Resign은 일렉션 리더십을 해제하여 다른 캠페이너들이 election에서 리더십을 획득할 수 있도록 한다. |



# How to conduct leader election in etcd cluster





# etcdctl cmd

`-w table` 은 테이블로 보여주는 옵션

## 리더/멤버 상태 확인

```bash 
etcdctl endpoint status --cluster -w table
```

## 헬스 체크

```bash
etcdctl --endpoints=127.0.0.1:2379,127.0.0.1:22379,127.0.0.1:32379 endpoint health -w table
```

## [모든 키 확인](https://github.com/etcd-io/etcd/issues/6904)

```bash
etcdctl --endpoints=localhost:2379 get --prefix --keys-only ''
```

