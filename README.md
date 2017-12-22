# Mobile robot fleet

Distributed control of a fleet of robotino using etcd.

![demo](https://user-images.githubusercontent.com/10092554/34301113-68036970-e72b-11e7-81e4-64c656bc24be.gif)


## Usage

Download this repo - e.g. with Go:
```bash
$ go get github.com/kindermoumoute/mobilerobotfleet
```

### Simulated with Docker

```bash
$ ./run.sh
```

### On the actual system

Build binaries with Go:
```bash
$ go build
```

Then run them on every node at the same time, specifying the pool ips. For example:
```bash
$ mobilerobotfleet.exe -pool 192.168.0.2,192.168.0.3
```

## Try it

Add jobs to the pool with the following curl:
```bash
$ curl http://127.0.0.1 -d A
```
