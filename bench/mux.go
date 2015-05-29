package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/funkygao/golib/gofmt"
	s "github.com/funkygao/golib/server"
	"github.com/funkygao/yamux"
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

func client() {
	t1 := time.Now()
	var wg sync.WaitGroup
	conn, err := net.Dial("tcp", addr)
	dieIfError(err)

	go stats()

	session, err := yamux.Client(conn, nil)
	dieIfError(err)
	fmt.Printf("session created for %s\n", addr)
	for i := 0; i < opts.c; i++ {
		wg.Add(1)
		go func(seq int) {
			defer wg.Done()

			msg := []byte(strings.Repeat("X", opts.sz))
			stream, err := session.OpenStream()
			dieIfError(err)
			for i := 0; i < opts.n; i++ {
				n, err := stream.Write(msg)
				dieIfError(err)
				atomic.AddInt64(&bytesW, int64(n))
			}
		}(i)
	}

	wg.Wait()
	fmt.Printf("%s\n", time.Since(t1))
}

func server() {
	go s.RunSysStats(time.Now(), time.Second*9)
	go stats()

	l, err := net.Listen("tcp", addr)
	dieIfError(err)
	fmt.Printf("listen on %s\n", addr)

	for {
		conn, err := l.Accept()
		dieIfError(err)

		go handleConn(conn)
	}

}

func handleConn(conn net.Conn) {
	session, err := yamux.Server(conn, nil)
	dieIfError(err)
	fmt.Printf("session created for %s\n", conn.RemoteAddr())

	for {
		stream, err := session.AcceptStream()
		if err == io.EOF {
			break
		}
		dieIfError(err)

		go handleStream(stream)
	}

}

func handleStream(st *yamux.Stream) {
	b := make([]byte, opts.sz)
	for {
		n, err := st.Read(b)
		if err == io.EOF {
			return
		}
		dieIfError(err)
		atomic.AddInt64(&bytesR, int64(n))
	}
}
