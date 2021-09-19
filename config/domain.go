package config

type Config struct {
	SourceURI     string `env:"SOURCE_URI" default:"http://:9000/events"`
	WasmBucketURI string `env:"WASM_BLOB_URI" default:"file://./wasm/assemblyscript/examples"`
	WasmFilename  string `env:"WASM_FILENAME" default:"main.wasm"`
	SinkURI       string `env:"SINK_URI" default:"nats://sink-topic"`
}

func (c Config) String() string {
	return "Config{SourceURI: " + c.SourceURI + ", WasmBucketURI: " + c.WasmBucketURI + ", WasmFilename: " + c.WasmFilename + ", SinkURI: " + c.SinkURI + "}"
}
