package main

import (
	"io"
	"log"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/funkygao/yamux"
)

func client() {
	t1 := time.Now()
	var wg sync.WaitGroup

	conn, err := net.Dial("tcp", addr)
	dieIfError(err)
	log.Printf("connected with %s\n", addr)

	session, err := yamux.Client(conn, nil)
	dieIfError(err)
	log.Printf("session created for %s\n", addr)
	for i := 0; i < opts.c; i++ {
		wg.Add(1)

		go func(seq int) {
			defer wg.Done()

			msg := []byte(strings.Repeat("X", opts.sz))
			b := make([]byte, opts.sz)
			stream, err := session.OpenStream()
			dieIfError(err)
			for i := 0; i < opts.n; i++ {
				n, err := stream.Write(msg)
				dieIfError(err)
				addByteWritten(n)
				n, err = stream.Read(b)
				dieIfError(err)
				addByteRead(n)
			}
		}(i)
	}

	wg.Wait()
	log.Printf("%s\n", time.Since(t1))
}

func server() {
	l, err := net.Listen("tcp", addr)
	dieIfError(err)
	log.Printf("listen on %s\n", addr)

	for {
		conn, err := l.Accept()
		dieIfError(err)

		go handleConn(conn)
	}

}

func handleConn(conn net.Conn) {
	session, err := yamux.Server(conn, nil)
	dieIfError(err)
	log.Printf("session created for %s\n", conn.RemoteAddr())

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
	response := []byte(strings.Repeat("Y", opts.sz))
	for {
		n, err := st.Read(b)
		if err == io.EOF {
			return
		}
		dieIfError(err)
		addByteRead(n)

		// simulate the biz logic overhead
		time.Sleep(time.Millisecond * 50) // 50ms

		n, err = st.Write(response)
		dieIfError(err)
		addByteWritten(n)
	}
}
