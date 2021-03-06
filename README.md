# gopcp_stream

Stream protocol supporting for golang pcp

## Stream Behaviour

```
                                      request stream

+-----------------+ +-------------------------------------------> +-------------------+
|                 |                                               |                   |
|                 |                                               |                   |
|                 |               response chunk by chunk         |                   |
|      Client     | <-------------------------------------------+ |      Server       |
|                 |                                               |                   |
|                 |                                               |                   |
|                 |              response end or error            |                   |
+-----------------+ <-------------------------------------------+ +-------------------+
```

## Quick Example

```go
pcpClient := gopcp.PcpClient{}
streamClient := GetStreamClient()

// client side
clientSide := gopcp.NewPcpServer(gopcp.GetSandbox(map[string]*gopcp.BoxFunc{
    "stream_accept": GetPcpStreamAcceptBoxFun(streamClient),
}))

// server side
streamServer := GetStreamServer("stream_accept", func(command string, timeout time.Duration) (interface{}, error) {
    return clientSide.Execute(command, nil)
})
serverSide := gopcp.NewPcpServer(gopcp.GetSandbox(map[string]*gopcp.BoxFunc{
	"streamApi": streamServer.StreamApi(func(streamProducer StreamProducer, args []interface{}, attachment interface{}, pcpServer *gopcp.PcpServer) (interface{}, error) {
		streamProducer.SendData("1", 10*time.Second)
		streamProducer.SendData("2", 10*time.Second)
		streamProducer.SendData("3", 10*time.Second)
		streamProducer.SendEnd(10 * time.Second)
		return nil, nil
	}),
}))
// client call server stream api
sum := ""
callExp, _ := streamClient.StreamCall("streamApi", func(t int, d interface{}) {
	if t == STREAM_DATA {
		sum += d.(string)
	}
})
cmd, _ := pcpClient.ToJSON(callExp)
serverSide.Execute(cmd, nil)
```
