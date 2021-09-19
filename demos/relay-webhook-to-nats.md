## Relaying data from a Webhook to a NATS Broker

In one terminal, run `nats-server` (on macOS you can install using Homebrew `brew install nats-server`):

```
> nats-server
[22080] 2021/09/18 14:34:03.933806 [INF] Starting nats-server version 2.1.9
[22080] 2021/09/18 14:34:03.933927 [INF] Git commit [not set]
[22080] 2021/09/18 14:34:03.934138 [INF] Listening for client connections on 0.0.0.0:4222
[22080] 2021/09/18 14:34:03.934144 [INF] Server id is NDBKLGPSIJOPXX2JAVTKL23W27S6YDWSXU3QSQLYXTZTW7GDCLMR3N5C
[22080] 2021/09/18 14:34:03.934146 [INF] Server is ready
```

Them start a `wasm-pipes` process on another terminal that receives data over HTTP and sends to NATS broker.

```
> SOURCE_URI=http://:9000/events SINK_URI=nats://sink NATS_SERVER_URL=0.0.0.0:4222 wasm-pipes
2021/09/18 14:39:05 Starting source on port 9000
2021/09/18 14:39:05  info wasm-pipes started
```

To test the pipeline, you can run `curl` to send data to the source:

```
> curl -X POST -H "Content-Type: application/json" -d '{"event_type": "test", "data": {"test": "test"}}' http://localhost:9000/events
```

To make sure the pipeline is working, we can create a loop with another isntance of `wasm-pipes`:

```
> SOURCE_URI=nats://sink?pipe=queue NATS_SERVER_URL=0.0.0.0:4222 SINK_URI=http://localhost:9000/events  wasm-pipes
2021/09/18 14:39:05  info wasm-pipes started
```

Now you're gonna be able to see data arriving back and forth between the two pipeline instances, arriving over HTTP, passing through NATS and back over HTTP.
