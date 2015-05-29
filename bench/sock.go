package main

import (
	"fmt"
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
			fmt.Printf("connected with %s\n", addr)

			msg := []byte(strings.Repeat("X", opts.sz))
			for j := 0; j < opts.n; j++ {
				n, err := conn.Write(msg)
				dieIfError(err)
				addByteWritten(n)
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

		fmt.Printf("got conn from %s\n", conn.RemoteAddr())
		go handleConn(conn)
	}

}

func handleConn(conn net.Conn) {
	b := make([]byte, opts.sz)
	for {
		b = b[:]
		n, err := conn.Read(b)
		dieIfError(err)
		addByteRead(n)
	}

}
