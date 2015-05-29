package main

import (
	"flag"
	"fmt"
	"os"
	"sync/atomic"
	"time"

	"github.com/funkygao/golib/gofmt"
)

const (
	addr = "localhost:10123"
	//addr = "10.77.144.193:10123"
)

var (
	opts struct {
		c    int
		sz   int
		n    int
		mode string
	}

	bytesR int64
	bytesW int64
)

func init() {
	flag.IntVar(&opts.c, "c", 100, "concurrency")
	flag.IntVar(&opts.sz, "s", 100, "size of each msg")
	flag.IntVar(&opts.n, "n", 50000, "loops count")
	flag.StringVar(&opts.mode, "m", "s", "<c|s> mode, c for client and s for server mode")

	flag.Parse()
}

func dieIfError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	switch opts.mode {
	case "s":
		server()

	case "c":
		client()

	default:
		flag.Usage()
		os.Exit(1)
	}

}

func stats() {
	tick := time.NewTicker(time.Second * 2)
	defer tick.Stop()

	for _ = range tick.C {
		r := atomic.LoadInt64(&bytesR)
		w := atomic.LoadInt64(&bytesW)
		fmt.Printf("r:%10s w:%10s\n", gofmt.ByteSize(r), gofmt.ByteSize(w))
	}
}
