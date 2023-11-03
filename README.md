# Trying Reliable ZeroMQ 1-to-1 PubSub With Socket Monitor

http://api.zeromq.org/4-1:zmq-socket-monitor

When sub disconnects, pub will store message to a queue, and will send to sub once sub reconnects.

## How to run

1. make sure your go version is >= 1.21
```
$ go mod tidy
```

2. run the publisher
```
$ cd pub/
$ go run pub.go
```

3. run the subscriber on another terminal
```
$ cd sub/
$ go run sub.go
```

4. try to `ctrl+c` the subscriber, notice how publisher is now sending to a `queue` instead.

5. try to run the subscriber again, notice how it is getting data from the queue first before getting new datas.

## Flowchart

### Publisher
![Alt text](assets/publisher-flowchart.png?raw=true "Publisher Flowchart")

### Subscriber
![Alt text](assets/subscriber-flowchart.png?raw=true "Subscriber Flowchart")