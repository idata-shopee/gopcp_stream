package gopcp_stream

import (
	"errors"
	"github.com/lock-free/gopcp"
	"time"
)

// at server side
// (1) will get client stream request
// (2) can Send data piece by piece (if error happened, stop Sending data)
// (3) can Send end or error (but only once)

// eg: some_remote_stream(streamId) = { SendData(streamId, [...]); SendEnd(streamId); /*SendError(streamId, err)*/ }

// call(clientAcceptName, streamId, t, d)

type StreamServer struct {
	clientAcceptName string
	callFun          CallFunc
	pcpClient        gopcp.PcpClient
}

// call function type
type CallFunc = func(command string, timeout time.Duration) (interface{}, error)

func (ss *StreamServer) SendData(streamId string, data interface{}, timeout time.Duration) (interface{}, error) {

	if cmd, err := ss.pcpClient.ToJSON(ss.pcpClient.Call(ss.clientAcceptName, streamId, STREAM_DATA, data)); err != nil {
		return nil, err
	} else {
		return ss.callFun(cmd, timeout)
	}
}

func (ss *StreamServer) SendEnd(streamId string, timeout time.Duration) (interface{}, error) {
	if cmd, err := ss.pcpClient.ToJSON(ss.pcpClient.Call(ss.clientAcceptName, streamId, STREAM_END)); err != nil {
		return nil, err
	} else {
		return ss.callFun(cmd, timeout)
	}
}

func (ss *StreamServer) SendError(streamId string, errMsg string, timeout time.Duration) (interface{}, error) {
	if cmd, err := ss.pcpClient.ToJSON(ss.pcpClient.Call(ss.clientAcceptName, streamId, STREAM_ERROR, errMsg)); err != nil {
		return nil, err
	} else {
		return ss.callFun(cmd, timeout)
	}
}

type StreamProducer struct {
	streamId string
	ss       *StreamServer
}

func (ps *StreamProducer) SendData(d interface{}, timeout time.Duration) {
	ps.ss.SendData(ps.streamId, d, timeout)
}

func (ps *StreamProducer) SendEnd(timeout time.Duration) {
	ps.ss.SendEnd(ps.streamId, timeout)
}

func (ps *StreamProducer) SendError(errMsg string, timeout time.Duration) {
	ps.ss.SendError(ps.streamId, errMsg, timeout)
}

func (ss *StreamServer) StreamApi(handle func(StreamProducer, []interface{}, interface{}, *gopcp.PcpServer) (interface{}, error)) *gopcp.BoxFunc {
	// (...args, streamId)
	return gopcp.ToSandboxFun(func(args []interface{}, attachment interface{}, pcpServer *gopcp.PcpServer) (interface{}, error) {
		if len(args) < 1 {
			return nil, errors.New("missing stream id at the stream request")
		} else if streamId, ok := args[len(args)-1].(string); !ok {
			return nil, errors.New("missing stream id (string) at the stream request")
		} else {
			streamProducer := StreamProducer{streamId, ss}
			return handle(streamProducer, args[:len(args)-1], attachment, pcpServer)
		}
	})
}

// TODO support stream lazy api

func GetStreamServer(clientAcceptName string, callFun CallFunc) *StreamServer {
	pcpClient := gopcp.PcpClient{}

	return &StreamServer{clientAcceptName, callFun, pcpClient}
}
