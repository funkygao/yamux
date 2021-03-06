package main

import (
	"flag"
	"log"
	"os"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/funkygao/golib/gofmt"
	s "github.com/funkygao/golib/server"
)

const (
	local  = "localhost:10123"
	remote = "10.77.144.193:10123"
)

var (
	addr string
	opts struct {
		c     int
		sz    int
		n     int
		mode  string
		local bool
	}

	bytesR int64
	bytesW int64
)

func init() {
	flag.IntVar(&opts.c, "c", 100, "concurrency")
	flag.IntVar(&opts.sz, "s", 100, "size of each msg")
	flag.IntVar(&opts.n, "n", 50000, "loops count")
	flag.StringVar(&opts.mode, "m", "s", "<c|s> mode, c for client and s for server mode")
	flag.BoolVar(&opts.local, "l", true, "local mode")

	flag.Parse()
	if opts.local {
		addr = local
	} else {
		addr = remote
	}

	runtime.GOMAXPROCS(runtime.NumCPU())
}

func dieIfError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	switch opts.mode {
	case "s":
		go s.RunSysStats(time.Now(), time.Second*9)
		go stats()
		server()

	case "c":
		go stats()
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
		log.Printf("r:%8s w:%8s\n", gofmt.ByteSize(r), gofmt.ByteSize(w))
	}
}

func addByteRead(n int) {
	atomic.AddInt64(&bytesR, int64(n))
}

func addByteWritten(n int) {
	atomic.AddInt64(&bytesW, int64(n))
}
