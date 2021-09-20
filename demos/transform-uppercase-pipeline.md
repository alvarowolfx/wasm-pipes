## Transform pipeline - Uppercase

[Add drawing of pipeline]

### HTTP Entry point pipeline

```
> SOURCE_URI=http://:9000/events SINK_URI=https://xyz.m.pipedream.net WASM_ENABLED=true WASM_BLOB_URI=file://./wasm/assemblyscript/examples/uppercase wasm-pipes
```

### Sending events

```
# Converted to { "NAME" : "MY MESSAGE" }
> curl -X POST http://localhost:9000/events -d '{ "name" : "My message" }' -H "Content-Type: application/json"
# Converted to { "NAME" : "ALL LOWERCASE" }
> curl -X POST http://localhost:9000/events -d 'all lowercase'
```

### Output

```
2021/09/19 00:14:04  info got message on engine: "{ \"name\" : \"My message\" }"
2021/09/19 00:14:04  info received event: { "name" : "My message" } module=wasm-log
2021/09/19 00:14:04  info uppercased event: { "NAME" : "MY MESSAGE" } module=wasm-log
2021/09/19 00:14:04  info event processed successfully
2021/09/19 00:14:05  info got message to send to sink: "{ \"NAME\" : \"MY MESSAGE\" }"
2021/09/19 00:16:35  info got message on engine: "all lowercase"
2021/09/19 00:16:35  info received event: all lowercase module=wasm-log
2021/09/19 00:16:35  info uppercased event: ALL LOWERCASE module=wasm-log
2021/09/19 00:16:35  info event processed successfully
2021/09/19 00:16:35  info got message to send to sink: "ALL LOWERCASE"
```
