package gopcp_stream

// at client side
// (1) will request remote stream
// (2) can accept stream data or end or error from remote
// In this design, we seperate those two procedures, and use uuid to connect them.
// In calling interface view, it should be one procedure, but at bottom it could be seperated.
// eg: call("some_remote_stream", streamCallback((type, data) => {}))
