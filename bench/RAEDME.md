benchmark comparison between multiplex and multiple socket
==========================================================

### sock.go

client fork n threads, each create a new socket, and write to the socket 
for many loops.

server accpeted socket each will fork a thread to read the socket and discard.

### mux.go

client dial server only 1 socket, and fork n threads, each open a new 
stream on this socket, and each thread write to the stream for many loops.

server accepted stream each will fork a thread to read from the stream and discard.

### Usage

go run sock.go -h

go run mux.go -h

### What to compare

- mem usage
- cpu usage(sys vs usr)
- network throughput
- interrupts per second
- context switch per second
- which is fast to finish 1M messages


