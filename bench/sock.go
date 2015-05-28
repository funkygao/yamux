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
	addr = "localhost:1234"
)

var opts struct {
	c    int
	sz   int
	n    int
	mode string
}

func init() {
	flag.IntVar(&opts.c, "c", 100, "concurrency")
	flag.IntVar(&opts.sz, "s", 100, "size of each msg")
	flag.IntVar(&opts.n, "n", 1000000, "loops count")
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
		flag.PrintDefaults()
		os.Exit(0)
	}

}

func client() {
	t1 := time.Now()
	var wg sync.WaitGroup
	for i := 0; i < opts.c; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			conn, err := net.Dial("tcp", addr)
			dieIfError(err)

			msg := []byte(strings.Repeat("X", opts.sz))
			for i := 0; i < opts.n; i++ {
				_, err := conn.Write(msg)
				dieIfError(err)
			}
		}()
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
		dieIfError(err)
	}

}
