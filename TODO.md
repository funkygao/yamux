TODO
====

### Terms

- socket
- session
- stream

### Internals

            socket
               |
               | wrap  +---- sendLoop
               |       |
            session ---|---- [keepaliveLoop]
               |       |
       map[id]stream   +---- recvLoop

### Why

- not reset stream id

### Fixed Sized Header

     0   1   2   3   4   5   6   7   8   9   A   B  
    +---+---+---+---+---+---+---+---+---+---+---+---+
    |ver|typ|flags  |streamId       |length         |
    +---+---+---+---+---+---+---+---+---+---+---+---+


### FrameType/MsgType

- data
- window

- ping
- bye
