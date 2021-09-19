package engine

import (
	"context"

	"github.com/eko/gocache/cache"
	"github.com/wasmerio/wasmer-go/wasmer"
	"gocloud.dev/blob"
	"gocloud.dev/pubsub"
)

type Engine interface {
	Start()
	Shutdown(ctx context.Context) error
}

type EngineDeps struct {
	Filename string
	Cache    *cache.Cache
	Input    *pubsub.Subscription
	Output   *pubsub.Topic
	Bucket   *blob.Bucket
}

type ExecVars struct {
	Message  *pubsub.Message
	Instance *wasmer.Instance
}

type Plugin interface {
	Init(deps *EngineDeps, wasmStore *wasmer.Store, vars *ExecVars) error
	ImportedObjects() map[string]wasmer.IntoExtern // Register import objects
}
