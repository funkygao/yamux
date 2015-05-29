package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"
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
	flag.IntVar(&opts.n, "n", 100000, "loops count")
	flag.StringVar(&opts.mode, "m", "server", "client or server mode")

	flag.Parse()
}

func dieIfError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	switch opts.mode {
	case "server":
		server()

	case "client":
		client()

	default:
		flag.Usage()
		os.Exit(0)
	}

}

func client() {
	t1 := time.Now()
	var wg sync.WaitGroup
	for i := 0; i < opts.c; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			println("connecting", addr)
			conn, err := net.Dial("tcp", addr)
			dieIfError(err)

			msg := []byte(strings.Repeat("X", opts.sz))
			for j := 0; j < opts.n; j++ {
				_, err := conn.Write(msg)
				//println(i, j)
				dieIfError(err)
			}
		}(i)
	}

	wg.Wait()
	fmt.Printf("%s\n", time.Since(t1))
}

func server() {
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
	b := make([]byte, opts.sz)

	for {
		b = b[:]
		_, err := conn.Read(b)
		//println(string(b[:n]))
		dieIfError(err)
	}

}
