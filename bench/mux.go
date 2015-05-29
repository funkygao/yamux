package main

import (
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	s "github.com/funkygao/golib/server"
	"github.com/funkygao/yamux"
)

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
