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
    "streamApi": gopcp.ToSandboxFun(func(args []interface{}, attachment interface{}, pcpServer *gopcp.PcpServer) (interface{}, error) {
		streamId := args[0].(string)

		streamServer.SendData(streamId, "1", 10*time.Second)
		streamServer.SendData(streamId, "2", 10*time.Second)
		streamServer.SendData(streamId, "3", 10*time.Second)
		streamServer.SendEnd(streamId, 10*time.Second)

		return nil, nil
	}),
}))
// client call server stream api
sum := ""
cmd, _ := pcpClient.ToJSON(
    pcpClient.Call("streamApi", streamClient.StreamCallback(func(t int, d interface{}) {
	    if t == STREAM_DATA {
	        sum += d.(string)
	    }
	})),
)
serverSide.Execute(cmd, nil)
```
