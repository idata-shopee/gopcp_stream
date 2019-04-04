package gopcp_stream

import (
	"fmt"
	"github.com/idata-shopee/gopcp"
	"testing"
	"time"
)

func assertEqual(t *testing.T, expect interface{}, actual interface{}, message string) {
	if expect == actual {
		return
	}
	if len(message) == 0 {
		message = fmt.Sprintf("expect %v !=  actual %v", expect, actual)
	}
	t.Fatal(message)
}

func TestBase(t *testing.T) {
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
			seed := args[0].(string)
			streamProducer.SendData(seed+"1", 10*time.Second)
			streamProducer.SendData(seed+"2", 10*time.Second)
			streamProducer.SendData(seed+"3", 10*time.Second)
			streamProducer.SendEnd(10 * time.Second)
			return nil, nil
		}),
	}))

	// client call server stream api
	sum := ""
	callExp, _ := streamClient.StreamCall("streamApi", "(", func(t int, d interface{}) {
		if t == STREAM_DATA {
			sum += d.(string)
		}
	})

	cmd, _ := pcpClient.ToJSON(*callExp)
	serverSide.Execute(cmd, nil)
	assertEqual(t, "(1(2(3", sum, "")
	streamClient.Clean()
}
