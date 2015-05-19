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

	stream, err := session.OpenStream()
	dieIfError(err)
	stream.Write([]byte("hello"))
	b := make([]byte, 100)
	stream.Read(b)
	println(string(b))

	stream1, err := session.OpenStream()
	dieIfError(err)
	stream1.Write([]byte("hello1"))
	stream1.Read(b)
	println(string(b))
}

func server() {
	l, err := net.Listen("tcp", addr)
	dieIfError(err)

	for {
		conn, err := l.Accept()
		dieIfError(err)

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
			break
		}
		dieIfError(err)

		b := make([]byte, 100)
		stream.Read(b)
		println(string(b))

		i++
		stream.Write([]byte(fmt.Sprintf("world%d", i)))
	}

}
