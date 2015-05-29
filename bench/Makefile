all:sock mux

clean:
	rm -f mux sock

sock:sock.go common.go
	go build -o sock sock.go common.go

mux:mux.go common.go
	go build -o mux mux.go common.go

