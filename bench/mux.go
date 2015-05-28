package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/funkygao/yamux"
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
	conn, err := net.Dial("tcp", addr)
	dieIfError(err)

	session, err := yamux.Client(conn, nil)
	dieIfError(err)
	for i := 0; i < opts.c; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			msg := []byte(strings.Repeat("X", opts.sz))
			stream, err := session.OpenStream()
			dieIfError(err)
			for i := 0; i < opts.n; i++ {
				_, err := stream.Write(msg)
				dieIfError(err)

				if i%(opts.n/100) == 1 {
					println(i)
				}
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
	session, err := yamux.Server(conn, nil)
	dieIfError(err)

	for {
		stream, err := session.AcceptStream()
		if err == io.EOF {
			stream.Close()
			break
		}
		dieIfError(err)

		go handleStream(stream)
	}

}

func handleStream(st *yamux.Stream) {
	b := make([]byte, opts.sz)
	for {
		_, err := st.Read(b)
		if err != nil {
			fmt.Println(err.Error())
			break
		}
	}
}
