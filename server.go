package gopcp_stream

// at server side
// (1) will get client stream request
// (2) can send data piece by piece (if error happened, stop sending data)
// (3) can send end or error (but only once)

// eg: some_remote_stream(streamId) = { sendData(streamId, [...]); sendEnd(streamId); /*sendError(streamId, err)*/ }
