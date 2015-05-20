package main

import (
	"fmt"
	"io"
	"net"
	"os"

	"github.com/funkygao/yamux"
)

const (
	addr = "localhost:1234"
)

func init() {
	yamux.Debug = true
}

func dieIfError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	switch os.Args[1] {
	case "server":
		server()

	case "client":
		client()
	}

}

func client() {
	conn, err := net.Dial("tcp", addr)
	dieIfError(err)

	session, err := yamux.Client(conn, nil)
	dieIfError(err)

	// 在一个socket上打开一个stream，每个stream长的都像个独立的socket
	// 每个stream可以独立地IO
	stream1, err := session.OpenStream()
	dieIfError(err)
	stream1.Write([]byte("hello"))
	b := make([]byte, 100)
	stream1.Read(b)
	println(string(b))

	// 再打开一个stream
	stream2, err := session.OpenStream()
	dieIfError(err)
	stream2.Write([]byte("hello1"))
	stream2.Read(b)
	println(string(b))

	// 现在问题来了，如何让stream1和stream2并发而非串行地运行?
	// 除了goroutine，还有别的办法吗?
}

func server() {
	l, err := net.Listen("tcp", addr)
	dieIfError(err)

	for {
		conn, err := l.Accept()
		dieIfError(err)

		// 每个conn都是一个socket
		go handleConn(conn)
	}

}

func handleConn(conn net.Conn) {
	session, err := yamux.Server(conn, nil)
	dieIfError(err)

	i := 0
	for {
		stream, err := session.AcceptStream()
		if err == io.EOF {
			stream.Close()
			break
		}
		dieIfError(err)

		// 这样的话，所有accept的stream就都是串行的，而非并行
		b := make([]byte, 100)
		stream.Read(b)
		println(string(b))

		i++
		stream.Write([]byte(fmt.Sprintf("world%d", i)))
	}

	/*
		// 为了让每个stream可以独立并行，应该是这样，而不是上面那样
		for {
			stream, err := session.AcceptStream()
			if err == io.EOF {
				stream.Close()
				break
			}

			go handleStream(stream)
		}
	*/

}

func handleStream(stream *yamux.Stream) {
}
