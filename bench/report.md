        free    virt    res     cpu     net     cs      interrupt   time
sock                                    30                          1m14.762873999s
mux                                                                 1m50.469296776s

在云机房的2台linux之间进行测试，一台服务器，另外一台客户端
比较mux(multiplex)和sock(多socket)在发送相同量数据的情况下的各自表现

总计发送500万个长度为100的消息
消息总计：476MB
mux用时2分半，sock用时5秒，mux网络上的数据包有100万次，而sock只有15万次

[root@imserver1 bench]# ./mux -l=false -m c
2015/05/28 19:54:30 connected with 10.77.144.193:10123
2015/05/28 19:57:08 2m37.117577196s
    net peak:       82Mbps
    cs peak:        80784
    tcp packets:    1048575     100万

[root@imserver1 bench]# ./sock -l=false -m c
2015/05/28 19:57:34 r:   0.00B w:231.35MB
2015/05/28 19:57:36 r:   0.00B w:449.98MB
2015/05/28 19:57:36 4.288915481s
    net peak:       907Mbps
    cs peak:        9541
    tcp packets:    147155      15万

