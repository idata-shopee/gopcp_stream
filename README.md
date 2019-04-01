# gopcp_stream

Stream protocol supporting for golang pcp

## Stream behaviour

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
