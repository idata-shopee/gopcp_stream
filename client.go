package gopcp_stream

import (
	"errors"
	"fmt"
	"github.com/idata-shopee/gopcp"
	"github.com/satori/go.uuid"
	"math"
	"sync"
)

// at client side
// (1) will request remote stream
// (2) can accept stream data or end or error from remote
// In this design, we seperate those two procedures, and use uuid to connect them.
// In calling interface view, it should be one procedure, but at bottom it could be seperated.
// eg: call("some_remote_stream", StreamCallback((type, data) => {}))

// if t is 0, d is the data of this chunk
// if t is 1, d is the error object
// if t is 2, d is nil
type StreamCallbackFunc = func(t int, d interface{})

// streamMap {streamId, callbackFunction}
type StreamClient struct {
	callbackMap sync.Map
	pcpClient   gopcp.PcpClient
}

// clean stream client to avoid memory leak
// eg: when connction is broken, clean it
func (sc *StreamClient) Clean() {
	sc.callbackMap = sync.Map{}
}

// register callback
func (sc *StreamClient) StreamCallback(callbackFunc StreamCallbackFunc) string {
	id := uuid.NewV4().String()
	sc.callbackMap.Store(id, callbackFunc)
	return id
}

// simple convension for stream calling
// (streamFunName, ...params, streamCallback)
func (sc *StreamClient) StreamCall(streamFunName string, params ...interface{}) (*gopcp.CallResult, error) {
	if len(params) < 1 {
		return nil, errors.New("missing stream callback function for stream call.")
	} else if callbackFun, ok := params[len(params)-1].(StreamCallbackFunc); !ok {
		return nil, errors.New("missing stream callback function for stream call.")
	} else {
		callbackParam := sc.StreamCallback(callbackFun)
		args := append(params[:len(params)-1], callbackParam)
		callExp := sc.pcpClient.Call(streamFunName, args...)
		return &callExp, nil
	}
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
		// when finished, remove callback from map
		if t == STREAM_END || t == STREAM_ERROR {
			sc.callbackMap.Delete(sid)
		}
		// call callback
		fun(t, d)
		return nil
	}
}

// define the stream accept sandbox function at client
func GetPcpStreamAcceptBoxFun(sc *StreamClient) *gopcp.BoxFunc {
	// args = [streamId: string, t: int, d: interface{}]
	return gopcp.ToSandboxFun(func(args []interface{}, attachment interface{}, pcpServer *gopcp.PcpServer) (interface{}, error) {
		if len(args) < 3 {
			return nil, streamFormatError(args)
		} else if streamId, ok := args[0].(string); !ok {
			return nil, streamFormatError(args)
		} else if t, ok := args[1].(float64); !ok {
			return nil, streamFormatError(args)
		} else {
			return nil, sc.Accept(streamId, int(math.Trunc(t)), args[2])
		}
	})
}

func streamFormatError(args []interface{}) error {
	return errors.New(fmt.Sprintf("stream chunk format: [streamId: string, t: int, d: interface{}]. args=%v", args))
}

func GetStreamClient() *StreamClient {
	var cm sync.Map
	sc := StreamClient{cm, gopcp.PcpClient{}}
	return &sc
}
