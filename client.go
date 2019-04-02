package gopcp_stream

import (
	"errors"
	"fmt"
	"github.com/idata-shopee/gopcp"
	"github.com/satori/go.uuid"
	"sync"
)

// at client side
// (1) will request remote stream
// (2) can accept stream data or end or error from remote
// In this design, we seperate those two procedures, and use uuid to connect them.
// In calling interface view, it should be one procedure, but at bottom it could be seperated.
// eg: call("some_remote_stream", StreamCallback((type, data) => {}))

// callback type: data, end, error

const STREAM_DATA = 0
const STREAM_END = 1
const STREAM_ERROR = 2

// if t is 0, d is the data of this chunk
// if t is 1, d is the error object
// if t is 2, d is nil
type StreamCallbackFunc = func(t int, d interface{})

// streamMap {streamId, callbackFunction}
type StreamClient struct {
	callbackMap sync.Map
}

// register callback
func (sc *StreamClient) StreamCallback(callbackFunc StreamCallbackFunc) string {
	id := uuid.NewV4().String()
	sc.callbackMap.Store(id, callbackFunc)
	return id
}

// accept stream response from server
func (sc *StreamClient) Accept(sid string, t int, d interface{}) error {
	if t != STREAM_DATA && t != STREAM_END && t != STREAM_ERROR {
		return errors.New("unepxpected stream chunk type.")
	} else if callbackFun, ok := sc.callbackMap.Load(sid); !ok {
		return errors.New(fmt.Sprintf("missing stream callback function for id: %s", sid))
	} else if fun, ok := callbackFun.(StreamCallbackFunc); !ok {
		return errors.New(fmt.Sprintf("stream callback function type error for id: %s. callbackFun=%v", sid, callbackFun))
	} else {
		if t == STREAM_END || t == STREAM_ERROR {
			sc.callbackMap.Delete(sid)
		}
		// call callback
		fun(t, d)
		return nil
	}
}

func GetPcpStreamAcceptBoxFun(sc *StreamClient) {
	// args = [streamId: string, t: int, d: interface{}]
	return gopcp.ToSandboxFun(func(args []interface{}, pcpServer *gopcpc.PcpServer) (interface{}, error) {
		if len(args) < 3 {
			return nil, errors.New("stream chunk format: [streamId: string, t: int, d: interface{}]")
		} else if streamId, ok := args[0].(string); !ok {
			return nil, errors.New("stream chunk format: [streamId: string, t: int, d: interface{}]")
		} else if t, ok := args[1].(int); !ok {
			return nil, errors.New("stream chunk format: [streamId: string, t: int, d: interface{}]")
		} else {
			d, ok := args[2]
			return nil, sc.Accept(streamId, t, d)
		}
	})
}

func GetStreamClient() *StreamClient {
	var cm sync.Map
	sc := StreamClient{cm}
	return &sc
}
