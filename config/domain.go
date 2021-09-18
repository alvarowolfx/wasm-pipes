package config

type Config struct {
	SourceURI   string `env:"SOURCE_URI" default:"http://:9000/events"`
	WasmBlobURI string `env:"WASM_BLOB_URI" default:"file://function.wasm"`
	SinkURI     string `env:"SINK_URI" default:"nats://sink-topic"`
	//SinkURI string `env:"SINK_URI" default:"http://:9000/events"`
}

func (c Config) String() string {
	return "Config{SourceURI: " + c.SourceURI + ", WasmBlobURI: " + c.WasmBlobURI + ", SinkURI: " + c.SinkURI + "}"
}
