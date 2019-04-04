package gopcp_stream

import (
	"github.com/idata-shopee/gopcp"
	"time"
)

// at server side
// (1) will get client stream request
// (2) can send data piece by piece (if error happened, stop sending data)
// (3) can send end or error (but only once)

// eg: some_remote_stream(streamId) = { sendData(streamId, [...]); sendEnd(streamId); /*sendError(streamId, err)*/ }

// call(clientAcceptName, streamId, t, d)

type StreamServer struct {
	clientAcceptName string
	callFun          CallFunc
	pcpClient        gopcp.PcpClient
}

// call function type
type CallFunc = func(command string, timeout time.Duration) (interface{}, error)

func (ss *StreamServer) sendData(streamId string, data interface{}, timeout time.Duration) (interface{}, error) {

	if cmd, err := ss.pcpClient.ToJSON(ss.pcpClient.Call(ss.clientAcceptName, streamId, STREAM_DATA, data)); err != nil {
		return nil, err
	} else {
		return ss.callFun(cmd, timeout)
	}
}

func (ss *StreamServer) sendEnd(streamId string, timeout time.Duration) (interface{}, error) {
	if cmd, err := ss.pcpClient.ToJSON(ss.pcpClient.Call(ss.clientAcceptName, streamId, STREAM_END)); err != nil {
		return nil, err
	} else {
		return ss.callFun(cmd, timeout)
	}
}

func (ss *StreamServer) sendError(streamId string, errMsg string, timeout time.Duration) (interface{}, error) {
	if cmd, err := ss.pcpClient.ToJSON(ss.pcpClient.Call(ss.clientAcceptName, streamId, STREAM_ERROR, errMsg)); err != nil {
		return nil, err
	} else {
		return ss.callFun(cmd, timeout)
	}
}

func GetStreamServer(clientAcceptName string, callFun CallFunc) StreamServer {
	pcpClient := gopcp.PcpClient{}

	return StreamServer{clientAcceptName, callFun, pcpClient}
}
