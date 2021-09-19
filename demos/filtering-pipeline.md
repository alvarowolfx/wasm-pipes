## Sample filtering/routing pipeline - Cold/Hot Storage

[Add drawing of pipeline]

### HTTP Entry point pipeline (1)

```
> SOURCE_URI=http://:9000/events SINK_URI=nats://check-temp NATS_SERVER_URL=0.0.0.0:4222 WASM_ENABLED=false wasm-pipes
```

### Cold storage Pipeline (2)

```
> SINK_URI=https://xyz.m.pipedream.net?storage=cold SOURCE_URI=nats://check-temp?queue=cold NATS_SERVER_URL=0.0.0.0:4222 WASM_ENABLED=true WASM_BLOB_URI=file://./wasm/assemblyscript/examples/cold wasm-pipes
```

### Hot storage Pipeline (2)

```
> SINK_URI=https://xyz.m.pipedream.net?storage=hot SOURCE_URI=nats://check-temp?queue=hot NATS_SERVER_URL=0.0.0.0:4222 WASM_ENABLED=true WASM_BLOB_URI=file://./wasm/assemblyscript/examples/hot wasm-pipes
```

### Sending events

```
# Generate hot storage event
> curl -X POST http://localhost:9000/events -d '{ "temp" : 45.0 }' -H "Content-Type: application/json
# No event triggered
> curl -X POST http://localhost:9000/events -d '{ "temp" : 30.0 }' -H "Content-Type: application/json
# Generate cold storage event
> curl -X POST http://localhost:9000/events -d '{ "temp" : 15.0 }' -H "Content-Type: application/json
```
