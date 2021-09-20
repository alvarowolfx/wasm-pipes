## Dedup pipeline - Filter duplicated messages arriving in sequence

Uses cache service.

[Add drawing of pipeline]

### HTTP Entry point pipeline

```
> SOURCE_URI=http://:9000/events SINK_URI=https://xyz.m.pipedream.net WASM_ENABLED=true WASM_BLOB_URI=file://./wasm/assemblyscript/examples/dedup-name wasm-pipes
```

### Sending events

```
# Triggers event
> curl -X POST http://localhost:9000/events -d '{ "name" : "My message" }' -H "Content-Type: application/json"
# Filtered
> curl -X POST http://localhost:9000/events -d '{ "name" : "My message" }' -H "Content-Type: application/json"
# Triggers event
> curl -X POST http://localhost:9000/events -d '{ "name" : "New message" }' -H "Content-Type: application/json"
# Triggers event
> curl -X POST http://localhost:9000/events -d '{ "name" : "My message" }' -H "Content-Type: application/json"
# Triggers event
> curl -X POST http://localhost:9000/events -d '{ "name" : "New message" }' -H "Content-Type: application/json"
# Filtered
> curl -X POST http://localhost:9000/events -d '{ "name" : "New message" }' -H "Content-Type: application/json"
```

### Output

```
2021/09/19 00:19:39  info got message on engine: "{ \"name\" : \"My message\" }"
2021/09/19 00:19:39  info received event: { "name" : "My message" } module=wasm-log
2021/09/19 00:19:39  info name: My message          module=wasm-log
2021/09/19 00:19:39  info name value changed, setting new value on cache module=wasm-log
2021/09/19 00:19:39  info sending name forward      module=wasm-log
2021/09/19 00:19:39  info event processed successfully
2021/09/19 00:19:39  info got message to send to sink: "{ \"name\" : \"My message\" }"
2021/09/19 00:19:39  info response from sink: 200 OK
2021/09/19 00:19:47  info got message on engine: "{ \"name\" : \"My message\" }"
2021/09/19 00:19:47  info received event: { "name" : "My message" } module=wasm-log
2021/09/19 00:19:47  info name: My message          module=wasm-log
2021/09/19 00:19:47  info repeated name, not sending module=wasm-log
2021/09/19 00:19:47  info event processed successfully
2021/09/19 00:20:09  info got message on engine: "{ \"name\" : \"New message\" }"
2021/09/19 00:20:09  info received event: { "name" : "New message" } module=wasm-log
2021/09/19 00:20:09  info name: New message         module=wasm-log
2021/09/19 00:20:09  info name value changed, setting new value on cache module=wasm-log
2021/09/19 00:20:09  info sending name forward      module=wasm-log
2021/09/19 00:20:09  info event processed successfully
2021/09/19 00:20:09  info got message to send to sink: "{ \"name\" : \"New message\" }"
2021/09/19 00:20:09  info response from sink: 200 OK
```
