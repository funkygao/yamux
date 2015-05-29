package main

import (
	"io"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

func client() {
	t1 := time.Now()
	var wg sync.WaitGroup
	for i := 0; i < opts.c; i++ {
		wg.Add(1)

		go func(seq int) {
			defer wg.Done()

			conn, err := net.Dial("tcp", addr)
			dieIfError(err)
			log.Printf("[%3d]connected with %s\n", seq, addr)

			msg := []byte(strings.Repeat("X", opts.sz))
			b := make([]byte, opts.sz)
			for j := 0; j < opts.n; j++ {
				n, err := conn.Write(msg)
				dieIfError(err)
				addByteWritten(n)
				n, err = conn.Read(b)
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

		log.Printf("got conn from %s\n", conn.RemoteAddr())
		go handleConn(conn)
	}

}

func handleConn(conn net.Conn) {
	b := make([]byte, opts.sz)
	response := []byte(strings.Repeat("Y", opts.sz))
	for {
		b = b[:]
		n, err := conn.Read(b)
		if err == io.EOF {
			return
		}
		dieIfError(err)
		addByteRead(n)

		// simulate biz logic overhead
		time.Sleep(time.Millisecond * 50)

		n, err = conn.Write(response)
		dieIfError(err)
		addByteWritten(n)
	}

}
