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
