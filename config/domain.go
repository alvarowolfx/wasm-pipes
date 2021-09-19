package config

type Config struct {
	SourceURI     string `env:"SOURCE_URI" default:"http://:9000/events"`
	WasmEnabled   bool   `env:"WASM_ENABLED" default:"false"`
	WasmBucketURI string `env:"WASM_BLOB_URI" default:"file://./wasm/assemblyscript/examples"`
	WasmFilename  string `env:"WASM_FILENAME" default:"main.wasm"`
	SinkURI       string `env:"SINK_URI" default:"nats://sink-topic"`
}

func (c Config) String() string {
	out := "Config{SourceURI: " + c.SourceURI + ", "
	out += "SinkURI: " + c.SinkURI + ", "
	if c.WasmEnabled {
		out += "WasmEnabled: true, "
		out += "WasmBucketURI: " + c.WasmBucketURI + ", "
		out += "WasmFilename: " + c.WasmFilename + "}"
	} else {
		out += "WasmEnabled: false }"
	}
	return out
}
